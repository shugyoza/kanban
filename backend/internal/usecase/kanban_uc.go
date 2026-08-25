/*
Calls the database repository method , loops through the flat rows, and packages them into the clean nested BoardAggregate tree that the frontend expects.

In Hexagonal (Ports and Adapters) and Clean Architecture, this is the "UseCase Interactor" layer that sits between the delivery layer (HTTP Handler) and the data access layer (Repository).
It orchestrates the business logic and data transformation, and represent a single use case of the application (in this case, "GetBoardDetails").

Core Responsibilities of an Interactor:
* Orchestrate the flow of data between the delivery layer and the repository layer.
* Execute One Use Case: It represents a specific user action (e.g., "GetBoardDetails").
* Enforce business rules and validation (e.g., check if board exists).
* Transform raw data models into a format suitable for the delivery layer (e.g., nested aggregates).
* Handle errors and propagate them back to the delivery layer in a meaningful way.
* Remains technology-agnostic: It doesn't know about HTTP, JSON, or any specific database technology. It only knows about the domain models and interfaces.
*/

package usecase

import (
	"backend/internal/domain"
	"context"
	"fmt"
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
	//  1.  Fetch the raw models from the database repository layer
	boardModel, columnModels, taskModels, err := uc.repo.GetBoardTree(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("usecase failed to get board details: %w", err)
	}

	// 2. Handle empty state (if the board ID was not found)
	if boardModel == nil || boardModel.ID == "" {
		return nil, fmt.Errorf("board not found")
	}

	// 3. Map tasks into a fast-lookup map grouped by column ID
	taskMap := make(map[string][]domain.Task)
	for _, task := range taskModels {
		taskMap[task.ColumnID] = append(taskMap[task.ColumnID], domain.Task{
			ID:       task.ID,
			Title:    task.Title,
			Description: task.Description,
			Position: task.Position,
			ColumnID: task.ColumnID,
		})
	}

	// 4. Assemble the nested ColumnAggregates
	var columnAggregates []domain.ColumnAggregate
	for _, column := range columnModels {
		// Find any task that belongs to this specific column (default to empty slice if none)
		tasks := taskMap[column.ID]
		if (tasks == nil) {
			tasks = []domain.Task{} // if no tasks exist for this column, ensure we return an empty JSON array instead of null
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