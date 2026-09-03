package repository

import (
	"context"
	"database/sql"
	"fmt"
	"kanban-backend/internal/domain"
	"slices"
	"time"
)

type SQLBoardRepository struct {
	db *sql.DB
}

// NewSQLBoardRepository creates a new instance of SQLBoardRepository with the provided database connection.
func NewSQLBoardRepository(db *sql.DB) *SQLBoardRepository {
	return &SQLBoardRepository{db: db}
}

func (r *SQLBoardRepository) GetBoardTree(ctx context.Context, boardID string) (*domain.BoardAggregate, error) {
	query := `
		SELECT 
			b.id, b.title,
			c.id, c.title, c.position,
			t.id, t.column_id, t.title, t.description, t.position
		FROM boards b
		LEFT JOIN columns c ON b.id = c.board_id
		LEFT JOIN tasks t ON c.id = t.column_id AND t.is_archived = 0
		WHERE b.id = $1
		ORDER BY c.position ASC, t.position ASC;
	`

	rows, err := r.db.QueryContext(ctx, query, boardID)
	if err != nil {
		return nil, fmt.Errorf("failed to query board tree: %w", err)
	}
	defer rows.Close()

	var board domain.BoardAggregate // Uses domain model directly
	columnMap := make(map[string]domain.ColumnAggregate)
	taskMap := make(map[string][]domain.Task)

	boardPopulated := false

	for rows.Next() {
		var bID, bTitle sql.NullString
		var cID, cTitle sql.NullString
		var cPosition sql.NullInt32
		var tID, tColumnID, tTitle, tDescription sql.NullString
		var tPosition sql.NullInt32

		err := rows.Scan(&bID, &bTitle, &cID, &cTitle, &cPosition, &tID, &tColumnID, &tTitle, &tDescription, &tPosition)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Populate the parent board details once
		if !boardPopulated && bID.Valid {
			board.ID = bID.String
			board.Title = bTitle.String
			boardPopulated = true
		}

		// 2. Collect unique columns found in the join
		if cID.Valid {
			if _, exists := columnMap[cID.String]; !exists {
			columnMap[cID.String] = domain.ColumnAggregate{
				ID: cID.String,
				Title: cTitle.String,
				Position: int(cPosition.Int32),
				Tasks: []domain.Task{},
			}

			}
		}

		// 3. Collect unique tasks found in the join
		if tID.Valid && tColumnID.Valid {
			taskMap[tColumnID.String] = append(taskMap[tColumnID.String], domain.Task{
				ID: tID.String,
				ColumnID: tColumnID.String,
				Title: tTitle.String,
				Description: tDescription.String,
				Position: int(tPosition.Int32),
			})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate board tree rows: %w", err)
	}

	// Flatten maps into slices to pass back to the business domain layer.
	var columns []domain.ColumnAggregate
	for colID, col := range columnMap {
		tasks := taskMap[colID]
		if tasks == nil {
			tasks = []domain.Task{} // Guarantees empty JSON array [] instead of null
		}

		// Explicitly sort the task cards vertically inside this column to override Go's map randomization, and ensure the sort order is enforced
		slices.SortFunc(tasks, func(a, b domain.Task) int {
			return a.Position - b.Position
		})

		col.Tasks = tasks
		columns = append(columns, col)
	}

	// Explicitly sort the columns array slice by position to override Go's map extraction randomization
	slices.SortFunc(columns, func(a, b domain.ColumnAggregate) int {
		return a.Position - b.Position
	})

	board.Columns = columns
	return &board, nil
}

func (r *SQLBoardRepository) UpdateTaskPositions(ctx context.Context, taskID string, targetColumnID string, targetPosition int) error {
	// 1. Initialize a strict ACID db transaction block
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start database index transaction: %w", err)
	}

	// Defer a rollback safety mechanism. If tx.Commit() is executed successfully, this becomes a safe no-op.
	defer tx.Rollback()

	// 2. Fetch the task's current column tracking metadata before the move
	var currentColumnID string
	var currentPosition int
	err = tx.QueryRowContext(ctx, "SELECT column_id, position FROM tasks WHERE id = $1", taskID).Scan(&currentColumnID, &currentPosition)
	if err != nil {
		return fmt.Errorf("failed to trace targeted task: %w", err)
	}

	// 3. Re-order items based on whether it moved columns or stayed within the same column
	if currentColumnID == targetColumnID {
		// Moving within the same column lane
		if currentPosition < targetPosition {
			// shifting down: push intermediate cards up
			_, err = tx.ExecContext(
				ctx, 
				"UPDATE tasks SET position = position - 1 WHERE column_id = $1 AND position > $2 AND position <= $3 AND is_archived = 0", 
				targetColumnID, currentPosition, targetPosition,
			)
		} else if currentPosition > targetPosition {
			// shifting up: push intermediate cards down
			_, err = tx.ExecContext(
				ctx,
				"UPDATE tasks SET position = position + 1 WHERE column_id = $1 AND position >= $2 AND position < $3 AND is_archived = 0",
				targetColumnID, targetPosition, currentPosition,
			)
		}
	} else {
		// Moving across columns: clear the gap in the old column lane
		_, err = tx.ExecContext(
			ctx,
			"UPDATE tasks SET position = position - 1 WHERE column_id = $1 AND position > $2 AND is_archived = 0",
			currentColumnID, currentPosition,
		)
		if err != nil {
			return fmt.Errorf("failed to re-index source column lane: %w", err)
		}

		// open up a placeholder index slot in the new target column lane
		_, err = tx.ExecContext(
			ctx,
			"UPDATE tasks SET position = position + 1 WHERE column_id = $1 AND position >= $2 AND is_archived = 0",
			targetColumnID, targetPosition,
		)
	}

	if err != nil {
		return fmt.Errorf("failed to adjust intermediate card positions: %w", err)
	}

	// 4. Update target card to point to its new column and destination position index
	_, err = tx.ExecContext(
		ctx,
		"UPDATE tasks SET column_id = $1, position = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3",
		targetColumnID, targetPosition, taskID,
	)
	if err != nil {
		return fmt.Errorf("failed to write updated target task metrics: %w", err)
	}

	// 5. Explicitly commit the transaction permanently to disk file storage
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to finalize task move transaction: %w", err)
	}

	return nil
}

func (r *SQLBoardRepository) InsertTask(ctx context.Context, columnID string, title string, description string) (*domain.Task, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start insert transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. OPTIMISTIC POSITIONING: Push all existing tasks in this column lane down by 1
	// to make room for a brand new task at position 0
	_, err = tx.ExecContext(
		ctx,
		"UPDATE tasks SET position = position + 1 WHERE column_id = $1 AND is_archived = 0",
		columnID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to re-index column for new task insertion: %w", err)
	}

	// 2. GENERATE TRACKING ID: Create a randomized high-entropy tracking string
	newID := fmt.Sprintf("task-%d", time.Now().UnixNano())
	defaultPosition := 0

	// 3. EXECUTE THE WRITE: Save the clean parameters permanently to the table rows
	query := `
	INSERT INTO tasks (id, column_id, title, description, position, is_archived, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	_, err = tx.ExecContext(ctx, query, newID, columnID, title, description, defaultPosition)
	if err != nil {
		return nil, fmt.Errorf("failed to insert raw task row: %w", err)
	}

	// 4. FINALIZE: Commit the transaction permanently to disk
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to finalize task creation transaction: %w", err)
	}

	// 5. RETURN VALUE: pass a pointer to the clean domain entity back up the stack
	return &domain.Task{
		ID: newID,
		ColumnID: columnID,
		Title: title,
		Description: description,
		Position: defaultPosition,
	}, nil
}

func (r *SQLBoardRepository) DeleteTask(ctx context.Context, columnID string, deletedTaskID string, deletedTaskPosition int) error {
	// 1. Initialize a strict ACID db transaction block
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start delete transaction: %w", err)
	}

	// Defer a rollback safety mechanism. If tx.Commit() is executed successfully, this becomes a safe no-op.
	defer tx.Rollback()

	// 2. Query to delete the row
	_, err = tx.ExecContext(
		ctx,
		"DELETE FROM tasks WHERE column_id = $1 AND id = $2",
		columnID,
		deletedTaskID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete raw task row: %w", err)
	}

	// Re-index all the rows within the column
	_, err = tx.ExecContext(
		ctx, 
		"UPDATE tasks SET position = position - 1 WHERE column_id = $1 AND position > $2 AND is_archived = 0", 
		columnID, deletedTaskPosition)
	if err != nil {
		return fmt.Errorf("failed to close position gap following deletion: %w", err)
	}


	// Explicitly commit the transaction permanently to disk file storage
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to finalize task delete transaction: %w", err)
	}

	return nil
}

func (r *SQLBoardRepository) UpdateTaskDetails(ctx context.Context, taskID string, title string, description string) error {
	query := `UPDATE tasks SET title = $1, description = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query, title, description, taskID)
	if err != nil {
		return fmt.Errorf("failed to update task details: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to verify update rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no task found with ID: %s", taskID)
	}

	return nil
}

func (r *SQLBoardRepository) ArchiveTask(ctx context.Context, columnID string,taskID string, taskPosition int) error {
	// 1. Initialize a strict ACID db transaction block
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start archiving transaction: %w", err)
	}

	defer tx.Rollback()

	// 2. Mark the task as archived in the database
	_, err = tx.ExecContext(
		ctx, 
		"UPDATE tasks SET is_archived = 1, updated_at = CURRENT_TIMESTAMP WHERE id = $1", 
		taskID)
	if err != nil {
		return fmt.Errorf("failed to archive task: %w", err)
	}

	// 3. Re-index all the rows within the column to close any gaps left by the archived task
	_, err = tx.ExecContext(
		ctx,
		"UPDATE tasks SET position = position - 1 WHERE column_id = $1 AND position > $2 AND is_archived = 0",
		columnID, taskPosition,
	)
	if err != nil {
		return fmt.Errorf("failed to re-index tasks after archiving: %w", err)
	}

	// 4. Explicitly commit the transaction permanently to disk file storage
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to finalize task archiving transaction: %w", err)
	}

	return nil
}

func (r *SQLBoardRepository) UnarchiveTask(ctx context.Context, columnID string, taskID string, taskPosition int) error {
	// 1. Initialize a strict ACID db transaction block
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start un-archiving transaction: %w", err)
	}
	defer tx.Rollback()

	// 2. Re-index all the rows within the column to make room for the un-archived task
	_, err = tx.ExecContext(
		ctx,
		"UPDATE tasks SET position = position + 1 WHERE column_id = $1 AND position >= $2 AND is_archived = 0",
		columnID, taskPosition,
	)
	if err != nil {
		return fmt.Errorf("failed to re-index tasks before un-archiving: %w", err)
	}

	// 3. Mark the task as un-archived in the database
	_, err = tx.ExecContext(
		ctx,
		"UPDATE tasks SET is_archived = 0, updated_at = CURRENT_TIMESTAMP WHERE id = $1",
		taskID,
	)
	if err != nil {
		return fmt.Errorf("failed to un-archive task: %w", err)
	}

	// 4. Explicitly commit the transaction permanently to disk file storage
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to finalize task un-archiving transaction: %w", err)
	}

	return nil
}