// *_test.go is skipped during normal compilation and only runs when running `go test`.

package usecase

import (
	"context"
	"strings"
	"testing"

	"kanban-backend/internal/domain"
)

// 1. Define the Mock Repository struct
type mockBoardRepository struct {
	getBoardTree func(ctx context.Context, boardID string) (*domain.BoardAggregate, error)
	updateTaskPositions func(ctx context.Context, taskID string, targetColumnID string, targetPosition int) error
	insertTask func(ctx context.Context, columnID string, title string, description string) (*domain.Task, error)
	deleteTask func(ctx context.Context, columnID string, deletedTaskID string, deletedTaskPosition int) error
	updateTaskDetails func(ctx context.Context, taskID string, title string, description string) error
	archiveTask func(ctx context.Context, columnID string, taskID string, taskPosition int) error
}

// 2. Implement the interface method so it satisfies domain.BoardRepository
func (m *mockBoardRepository) GetBoardTree(ctx context.Context, boardID string) (*domain.BoardAggregate, error) {
	return m.getBoardTree(ctx, boardID)
}

func (m *mockBoardRepository) UpdateTaskPositions(ctx context.Context, taskID string, targetColumnID string, targetPosition int) error {
	return m.updateTaskPositions(ctx, taskID, targetColumnID, targetPosition)
}

func (m *mockBoardRepository) InsertTask(ctx context.Context, columnID string, title string, description string) (*domain.Task, error) {
	return m.insertTask(ctx, columnID, title, description)
}

func (m *mockBoardRepository) DeleteTask(ctx context.Context, columnID string, deletedTaskID string, deletedTaskPosition int) error {
	return m.deleteTask(ctx, columnID, deletedTaskID, deletedTaskPosition)
}

func (m *mockBoardRepository) UpdateTaskDetails(ctx context.Context, taskID string, title string, description string) error {
	return m.updateTaskDetails(ctx, taskID, title, description)
}

func (m *mockBoardRepository) ArchiveTask(ctx context.Context, columnID string, taskID string, taskPosition int) error {
	return m.archiveTask(ctx, columnID, taskID, taskPosition)
}

// 3. The unit test function
func TestGetBoardDetails_Success(t *testing.T) {
	targetBoardID := "board-1"

	// Instantiate the mock repo and define its specific behavior for this test case.
	mockRepo := &mockBoardRepository{
		getBoardTree: func(ctx context.Context, boardID string) (*domain.BoardAggregate, error) {	
			return &domain.BoardAggregate{
				ID: targetBoardID,
				Title: "Engineering Board",
				Columns: []domain.ColumnAggregate{
					{
						ID: "col-todo",
						Title: "To Do",
						Position: 0,
						Tasks: []domain.Task{
							{ID: "task-1", ColumnID: "col-todo", Title: "Task 1", Description: "Setup Go Project", Position: 0},
							{ID: "task-2", ColumnID: "col-todo", Title: "Task 2", Description: "Configure Git", Position: 1},
						},
					},
					{
						ID:       "col-done",
						Title:    "Done",
						Position: 1,
						Tasks: []domain.Task{
							{ID: "task-3", ColumnID: "col-done", Title: "Task 3", Description: "Design Schema", Position: 0},
						},
					},
				}, 
			}, nil
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

func TestGetBoardDetails_NotFound(t *testing.T) {
	mockRepo := &mockBoardRepository{
		getBoardTree: func(ctx context.Context, boardID string) (*domain.BoardAggregate, error) {
			return nil, nil // simulate empty dataset recprd
		},
	}

	interactor := NewKanbanInteractor(mockRepo)
	_, err := interactor.GetBoardDetails(context.Background(), "missing-id")

	if err == nil || err.Error() != "board not found" {
		t.Errorf("Expected strict 'board not found' error exception, received: %v", err)
	}
}

func TestMoveTask_BusinessRuleViolations(t *testing.T) {
	interactor := NewKanbanInteractor(&mockBoardRepository{})

	// Test 1: Missing Task ID
	err := interactor.MoveTask(context.Background(), "", "col-1", 0)
	if err == nil || err.Error() != "business rule violation: task ID cannot be empty" {
		t.Errorf("Expected blank task validation break, got: %v", err)
	}

	// Test 2: Negative Position Index
	err = interactor.MoveTask(context.Background(), "task-1", "col-1", -5)
	if err == nil || err.Error() != "business rule violation: position index cannot be negative" {
		t.Errorf("Expected negative position index protection error, got: %v", err)
	}
}

func TestCreateTask_Success(t *testing.T) {
	mockRepo := &mockBoardRepository{
		insertTask: func(ctx context.Context, columnID string, title string, description string) (*domain.Task, error) {
			return &domain.Task{
				ID: "task-123",
				ColumnID: columnID,
				Title: title,
				Description: description,
				Position: 0,
			}, nil
		},
	}

	interactor := NewKanbanInteractor(mockRepo)
	task, err := interactor.CreateTask(context.Background(), "col-todo", "Write Unit Tests", "Cover backend core components")

	if err != nil {
		t.Fatalf("Expected clean execution pass, received unexpected error: %v", err)
	}

	if task.ID != "task-123" || task.Position != 0 {
		t.Errorf("Task properties corrupted during entity generation pipeline")
	}
}

func TestCreateTask_ValidationFails(t *testing.T) {
	interactor := NewKanbanInteractor(&mockBoardRepository{})

	// Test: attempt to insert task with empty title header
	_, err := interactor.CreateTask(context.Background(), "col-todo", "", "Missing Title Details")
	if err == nil || err.Error() != "business rule violation: task title cannot be left blank" {
		t.Errorf("Expected strict blank title check violation rule break, got: %v", err)
	}
}

// DELETE TASK VALIDATION & EXECUTION TEST CASES
func TestDeleteTask_Success(t *testing.T) {
	repositoryCalled := false

	mockRepo := &mockBoardRepository{
		deleteTask: func(ctx context.Context, columnID string, deletedTaskID string, deletedTaskPosition int) error {
			repositoryCalled = true

			// Verify that the UseCase safely forwarded the correct arguments down the wire
			if columnID != "col-todo" || deletedTaskID != "task-1" || deletedTaskPosition != 2 {
				t.Errorf("UseCase passed corrupted arguments to repository layer")
			}
			
			return nil
		},
	}

	interactor := NewKanbanInteractor(mockRepo)
	err := interactor.DeleteTask(context.Background(), "col-todo", "task-1", 2)

	if err != nil {
		t.Fatalf("Expected zero errors during clean task deletion, go: %v", err)
	}

	if !repositoryCalled {
		t.Errorf("UseCase failed to delegate data removal execution to the repository layer")
	}
}

func TestEditTask_Success(t *testing.T) {
	repositoryCalled := false

	mockRepo := &mockBoardRepository{
		updateTaskDetails: func(ctx context.Context, taskID string, title string, description string) error {
			repositoryCalled = true
			if taskID != "task-1" || title != "Updated Title" {
				t.Errorf("UseCase passed corrupted editing arguments to repository layer")
			}

			return nil
		},
	}

	interactor := NewKanbanInteractor(mockRepo)
	err := interactor.EditTask(context.Background(), "task-1", "Updated Title", "Updated Description")

	if err != nil {
		t.Fatalf("Expected zero errors during clean task edit, got: %v", err)
	}
	if !repositoryCalled {
		t.Errorf("UseCase failed to delegate task edit execution to the repository layer")
	}
}

func TestArchiveTask_Success(t *testing.T) {
	repositoryCalled := false

	mockRepo := &mockBoardRepository{
		archiveTask: func(ctx context.Context, columnID string, taskID string, taskPosition int) error {
			repositoryCalled = true
			if columnID != "col-1" || taskID != "task-1" || taskPosition != 0 {

				t.Errorf("UseCase passed corrupted arguments to repository layer")
			}

			return nil
		},
	}

	interactor := NewKanbanInteractor(mockRepo)
	err := interactor.ArchiveTask(context.Background(), "col-1", "task-1", 0)
	if err != nil {
		t.Fatalf("Expected zero errors during clean task archiving, got: %v", err)
	}
	if !repositoryCalled {
		t.Errorf("UseCase failed to delegate task archiving execution to the repository layer")
	}
}

func TestArchiveTask_BusinessRuleViolations(t *testing.T) {
	interactor := NewKanbanInteractor(&mockBoardRepository{})

	// Test case A: Blank parent track container identification exception
	err := interactor.ArchiveTask(context.Background(), "", "task-1", 0)
	if err == nil || !strings.Contains(err.Error(), "parent column ID is mandatory") {
		t.Errorf("Expected blank column verification fault, received: %v", err)
	}

	// Test case B: Negative position array grid check validation exception
	err = interactor.ArchiveTask(context.Background(), "col-done", "task-1", -1)
	if err == nil || !strings.Contains(err.Error(), "task position to archive is mandatory") {
		t.Errorf("Expected negative position blocker check trigger, received: %v", err)
	}
}