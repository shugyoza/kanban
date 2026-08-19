/* Calls the database repository method , loops through the flat rows, and packages -
them into the clean nested BoardAggregate tree that the frontend expects. */

package usecase

import (
	"context"
	"fmt"
	"kanban/internal/domain"
)

type KanbanInteractor struct {
	repo domain.BoardRepository
}

// NewKanbanInteractor initializes business logic layer with its data
func NewKanbanInteractor(repo domain.BoardRepository) *KanbanInteractor {
	return &KanbanInteractor{repo: repo}
}

// GetBoardDetails orchestrates the data retrieval and formats the nested tree
func (uc *KanbanInteractor) GetBoardDetails(ctx context.Context, boardID string) (*domain.BoardAggregate, error) {
	//  1.  Fetch the raw modles from the database repository layer
	boardModel, columnModels, taskModels, err := uc.repo.GetBoardDetails(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("usecase failed to get board details: %w", err)
	}

	// 2. Handle empty state (if the board ID was not found)
	if boardModel == nil || boardModel.ID == "" {
		return nil, fmt.Errorf("board not found")
	}

	// 3. Map tasks into a fast-lookup map grouped by column ID
	taskMap := make(map[string][]domain.TaskModel)
	for _, task := range taskModels {
		taskMap[task.ColumnID] = append(taskMap[task.ColumnID], task)
	}

	// 4. Assemble the nested ColumnAggregates
	var columnAggregates []domain.ColumnAggregate
	for _, column := range columnModels {
		// Find any task that belongs to this specific column (default to empty slice if none)
		tasks := taskMap[column.ID]
		if (tasks == nil) {
			tasks = []domain.TaskModel{} // if no tasks exist for this column, ensure we return an empty JSON array instead of null
		}

		columnAggregates = append(columnAggregates, domain.ColumnAggregate{
			ID:       column.ID,
			Title:    column.Title,
			Position: column.Position,
			Tasks:    tasks, // Attach tasks to the corresponding column
		})
	}

	// 5. Construct and return the final unified aggregate tree
	return &domain.BoardAggregate{
		ID:      boardModel.ID,
		Title:   boardModel.Title,
		Columns: columnAggregates,
	}, nil
}