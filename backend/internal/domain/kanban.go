// business entities, aggregators & interface ports

package domain

import (
	"context"
)

// KanbanUseCase defines the business rules available for the Kanban board.
type KanbanUseCase interface {
	GetBoardDetails(ctx context.Context, boardID string) (*BoardAggregate, error) // use * to pass a pointer to the BoardAggregate struct (instead of copied value), allowing for efficient memory usage and the ability to modify the original struct if needed.
}