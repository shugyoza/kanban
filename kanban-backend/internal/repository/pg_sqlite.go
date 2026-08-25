package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type SQLBoardRepository struct {
	db *sql.DB
}

// NewSQLBoardRepository creates a new instance of SQLBoardRepository with the provided database connection.
func NewSQLBoardRepository(db *sql.DB) *SQLBoardRepository {
	return &SQLBoardRepository{db: db}
}

func (r *SQLBoardRepository) GetBoardTree(ctx context.Context, boardID string) (*BoardModel, []ColumnModel, []TaskModel, error) {
	query := `
		SELECT 
			b.id, b.title,
			c.id, c.title, c.position,
			t.id, t.column_id, t.title, t.description, t.position
		FROM boards b
		LEFT JOIN columns c ON b.id = c.board_id
		LEFT JOIN tasks t ON c.id = t.column_id
		WHERE b.id = $1
		ORDER BY c.position ASC, t.position ASC;
	`

	rows, err := r.db.QueryContext(ctx, query, boardID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to query board tree: %w", err)
	}
	defer rows.Close()

	var board BoardModel
	columnMap := make(map[string]ColumnModel)
	taskMap := make(map[string]TaskModel)

	boardPopulated := false

	for rows.Next() {
		var bID, bTitle sql.NullString
		var cID, cTitle sql.NullString
		var cPosition sql.NullInt32
		var tID, tColumnID, tTitle, tDescription sql.NullString
		var tPosition sql.NullInt32

		err := rows.Scan(&bID, &bTitle, &cID, &cTitle, &cPosition, &tID, &tColumnID, &tTitle, &tDescription, &tPosition)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Populate the parent board details once
		if !boardPopulated && bID.Valid {
			board.ID = bID.String
			board.Title = bTitle.String
			boardPopulated = true
		}

		// 2. Collect unique columns found in the join
		if cID.Valid {
			columnMap[cID.String] = ColumnModel{
				ID: cID.String,
				BoardID: board.ID,
				Title: cTitle.String,
				Position: int(cPosition.Int32),
			}
		}

		// 3. Collect unique tasks found in the join
		if tID.Valid {
			taskMap[tID.String] = TaskModel{
				ID: tID.String,
				ColumnID: tColumnID.String,
				Title: tTitle.String,
				Description: tDescription.String,
				Position: int(tPosition.Int32),
			}
		}
	}

	// Flatten maps into slices to pass back to the business domain layer.
	var columns []ColumnModel
	for _, col := range columnMap {
		columns = append(columns, col)
	}

	var tasks []TaskModel
	for _, task := range taskMap {
		tasks = append(tasks, task)
	}

	return &board, columns, tasks, nil
}
