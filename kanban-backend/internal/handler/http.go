package handler

import (
	"encoding/json"
	"kanban-backend/internal/domain"
	"log"
	"net/http"
	"strings"
)

type KanbanHandler struct {
	useCase domain.KanbanUseCase
}

// MoveTaskPayload defines the strict JSON contract from the frontend
type MoveTaskPayload struct {
	TaskID string `json:"taskId"`
	TargetColumnID string `json:"targetColumnId"`
	TargetPosition int `json:"targetPosition"`
}

// NewKanbanHandler initializes the delivery layer with its required business logic dependency.
func NewKanbanHandler(uc domain.KanbanUseCase) *KanbanHandler {
	return &KanbanHandler{
		useCase: uc,
	}
}

// GetBoard handles requests matching: GET /api/boards/boards?id=xxx
func (handler *KanbanHandler) GetBoard(w http.ResponseWriter, r *http.Request) {
	// 1. Enforce strict HTTP method checking
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Extract the target board ID from the query string params
	boardID := r.URL.Query().Get("id")
	if boardID == "" {
		http.Error(w, "Missing required board 'id' query parameter", http.StatusBadRequest)
		return
	}

	// 3/ Trigger our core Hexagonal Business Interactor
	boardTree, err := handler.useCase.GetBoardDetails(r.Context(), boardID)
	if err != nil {
		if err.Error() == "board not found" {
			http.Error(w, "Board not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}

	// 4. Set standard web API response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 5. Serialize the boardTree pointer struct to camelCase JSON and write it to the response
	if err := json.NewEncoder(w).Encode(boardTree); err != nil {
		http.Error(w, "Failed to encode response payload", http.StatusInternalServerError)
	}
}

// MoveTask handles requests matching : PUT /api/tasks/move
func (h *KanbanHandler) MoveTask(w http.ResponseWriter, r *http.Request) {
	// 1. Enforce strict HTTP method filtering
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Decode and validate the incoming JSON request payload
	var payload MoveTaskPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Malformed JSON request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 3. Trigger core Hexagonal Business Interactor UseCase
	err := h.useCase.MoveTask(r.Context(), payload.TaskID, payload.TargetColumnID, payload.TargetPosition)
	if err != nil {
		log.Printf("Error executing task move workflow: %v", err)

		if strings.Contains(err.Error(), "business rule violation") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, "Internal server update failure", http.StatusInternalServerError)
		return
	}

	// 4. Return a clean HTTP 204 No Content status on a completely successful operation
	w.WriteHeader(http.StatusNoContent)
}