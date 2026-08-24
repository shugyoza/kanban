package handler

import (
	"backend/internal/domain"
	"encoding/json"
	"net/http"
)

type KanbanHandler struct {
	useCase domain.KanbanUseCase
}

// NewKanbanHandler initializes the delivery layer with its required business logic dependency.
func NewKanbanHandler(uc domain.KanbanUseCase) *KanbanHandler {
	return &KanbanHandler{
		useCase: uc,
	}
}

// GetBoard handles requests matching: GET /api/boards/boards?id=xxx
func (handler *KanbanHandler) GetBoard(response http.ResponseWriter, request *http.Request) {
	// 1. Enforce strict HTTP method checking
	if request.Method != http.MethodGet {
		http.Error(response, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Extract the target board ID from the query string params
	boardID := request.URL.Query().Get("id")
	if boardID == "" {
		http.Error(response, "Missing required board 'id' query parameter", http.StatusBadRequest)
		return
	}

	// 3/ Trigger our core Hexagonal Business Interactor
	boardTree, err := handler.useCase.GetBoardDetails(request.Context(), boardID)
	if err != nil {
		if err.Error() == "board not found" {
			http.Error(response, "Board not found", http.StatusNotFound)
		} else {
			http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		}
	}

	// 4. Set standard web API response headers
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	// 5. Serialize the boardTree pointer struct to camelCase JSON and write it to the response
	if err := json.NewEncoder(response).Encode(boardTree); err != nil {
		http.Error(response, "Failed to encode response payload", http.StatusInternalServerError)
	}
}