// *_test.go is skipped during normal compilation and only runs when running `go test`.

package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"backend/internal/repository"
)

// 1. Define the Mock Repository struct
type mockBoardRepository struct {
	mockGetBoardTree func(ctx context.Context, boardID string) (*repository.BoardModel, []repository.ColumnModel, []repository.TaskModel, error)

	mockUpdateTaskPositions func(ctx context.Context, taskID string, targetColumnID string, targetPosition int)
}

// 2. Implement the interface method so it satisfies domain.BoardRepository
func (m *mockBoardRepository) GetBoardTree(ctx context.Context, boardID string) (*repository.BoardModel, []repository.ColumnModel, []repository.TaskModel, error) {
	return m.mockGetBoardTree(ctx, boardID)
}

func (m *mockBoardRepository) UpdateTaskPositions(ctx context.Context, taskID string, targetColumnID string, targetPosition int) error {
	return nil
}

// 3. The unit test function
func TestGetBoardDetails_Success(t *testing.T) {
	// Arrange: set up dummy data that mimics raw SQL rows returned from the database.
	now := time.Now()
	targetBoardID := "board-1"

	dummyBoard := &repository.BoardModel{ID: targetBoardID, Title: "Test Board", CreatedAt: now}

	dummyColumns := []repository.ColumnModel{
		{ID: "col-todo", BoardID: targetBoardID, Title: "To Do", Position: 0, CreatedAt: now},
		{ID: "col-done", BoardID: targetBoardID, Title: "Done", Position: 1, CreatedAt: now},
	}

	dummyTasks := []repository.TaskModel{
		{ID: "task-1", ColumnID: "col-todo", Title: "Task 1", Description: "Setup Go Project", Position: 0, CreatedAt: now, UpdatedAt: now},
		{ID: "task-2", ColumnID: "col-todo", Title: "Task 2", Description: "Configure Git", Position: 1, CreatedAt: now, UpdatedAt: now},
		{ID: "task-3", ColumnID: "col-done", Title: "Task 3", Description: "Design Schema", Position: 0, CreatedAt: now, UpdatedAt: now},
	}

	// Instantiate the mock repo and define its specific behavior for this test case.
	mockRepo := &mockBoardRepository{
		mockGetBoardTree: func(ctx context.Context, boardID string) (*repository.BoardModel, []repository.ColumnModel, []repository.TaskModel, error) {
			if boardID == targetBoardID {
				return dummyBoard, dummyColumns, dummyTasks, nil
			}

			return nil, nil, nil, errors.New("unexpected board ID requested, and not found")
		},
	}

	// Inject the mock repo into real UseCase implementation.
	interactor := NewKanbanInteractor(mockRepo)

	// Act: Call the UseCase method under test.
	result, err := interactor.GetBoardDetails(context.Background(), targetBoardID)

	// Assert: Validate the output match architectural expectations.
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if result.ID != targetBoardID {
		t.Errorf("Expected board ID %s, but got %s", targetBoardID, result.ID)
	}

	if len(result.Columns) != 2 {
		t.Errorf("Expected 2 columns, but got %d", len(result.Columns))
	}

	// Validate that tasks are correctly grouped under their respective columns.
	todoColumn := result.Columns[0] // Will be "To Do" column because of the order in dummyColumns
	if len(todoColumn.Tasks) != 2 {
		t.Errorf("Expected 2 tasks in 'To Do' column, but got %d", len(todoColumn.Tasks))
	}

	doneColumn := result.Columns[1] // Will be "Done" column
	if len(doneColumn.Tasks) != 1 {
		t.Errorf("Expected 1 task in 'Done' column, but got %d", len(doneColumn.Tasks))
	}
}