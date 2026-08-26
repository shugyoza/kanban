// business entities, aggregators & interface ports

package domain

import (
	"backend/internal/repository"
	"context"
)

// Task represents a single task in the Kanban board.
type Task struct {
	ID string `json:"id"`
	ColumnID string `json:"columnId"`
	Title string `json:"title"`
	Description string `json:"description"`
	Position int `json:"position"`
}

type ColumnAggregate struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Position int `json:"position"`
	Tasks []Task `json:"tasks"`
}

// Representing clean tree structure in domain layer that matches the exact tree structure the UI needs.
type BoardAggregate struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Columns []ColumnAggregate `json:"columns"`
}

type BoardRepository interface {
	GetBoardTree(ctx context.Context, boardID string) (*repository.BoardModel, []repository.ColumnModel, []repository.TaskModel, error)

	UpdateTaskPositions(ctx context.Context, taskID string, targetColumnID string, targetPosition int) error
}

// KanbanUseCase defines the business rules available for the Kanban board.
type KanbanUseCase interface {
	GetBoardDetails(ctx context.Context, boardID string) (*BoardAggregate, error) // use * to pass a pointer to the BoardAggregate struct (instead of copied value), allowing for efficient memory usage and the ability to modify the original struct if needed.

	MoveTask(ctx context.Context, taskID string, targetColumnID string, targetPosition int) error
}