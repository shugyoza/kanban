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
	"context"
	"fmt"
	"kanban-backend/internal/domain"
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
	boardAggregate, err := uc.repo.GetBoardTree(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("usecase failed to get board details: %w", err)
	}

	// 2. Handle empty state (if the board ID was not found)
	if boardAggregate == nil || boardAggregate.ID == "" {
		return nil, fmt.Errorf("board not found")
	}

	return boardAggregate, nil
}

func (uc *KanbanInteractor) MoveTask(ctx context.Context, taskID string, targetColumnID string, targetPosition int) error {
	// 1. Enforce strict data input validations
	if taskID == "" {
		return fmt.Errorf("business rule violation: task ID cannot be empty")
	}

	if targetColumnID == "" {
		return fmt.Errorf("business rule violation: target column ID cannot be empty")
	}

	if targetPosition < 0 {
		return fmt.Errorf("business rule violation: position index cannot be negative")
	}

	// 2. Delegate the data execution safely to your SQL transaction repository port
	err := uc.repo.UpdateTaskPositions(ctx, taskID, targetColumnID, targetPosition)
	if err != nil {
		return fmt.Errorf("usecase failed to execute task move: %w", err)
	}

	return nil
}

func (uc *KanbanInteractor) CreateTask(ctx context.Context, columnID string, title string, description string) (*domain.Task, error) {
	if columnID == "" {
		return nil, fmt.Errorf("business rule violation: parent column ID is mandatory for task creation")
	}
	if (title == "") {
		return nil, fmt.Errorf("business rule violation: task title cannot be left blank")
	}

	createdTask, err := uc.repo.InsertTask(ctx, columnID, title, description)
	if err != nil {
		return nil, fmt.Errorf("usecase failed to create and insert task: %w", err)
	}

	return createdTask, nil
}

func (uc *KanbanInteractor) DeleteTask(ctx context.Context, columnID string, deletedTaskID string, deletedTaskPosition int) error {
	if columnID == "" {
		return fmt.Errorf("business rule violation: parent column ID is mandatory for task deletion")
	}

	if deletedTaskID == "" {
		return fmt.Errorf("business rule violation: task ID to delete is mandatory for task deletion")
	}

	if deletedTaskPosition < 0 {
		return fmt.Errorf("business rule violation: task position to delete is mandatory for task deletion")
	}

	err := uc.repo.DeleteTask(ctx, columnID, deletedTaskID, deletedTaskPosition)
	if err != nil {
		return fmt.Errorf("usecase failed to execute task deletion: %w", err)
	}

	return nil
}