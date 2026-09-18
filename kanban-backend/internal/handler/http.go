package handler

import (
	"encoding/json"
	"fmt"
	"kanban-backend/internal/domain"
	"log"
	"net/http"
	"strings"
	"time"
)

type EmailInputValidationRequest struct {
	Email string `json:"email"`
}

type InvitationRequest struct {
	UserID string `json:"userId"`
	Email string `json:"email"`
}

type InvitationResponse struct {
	Token string `json:"token"`
}

type InvitationTokenValidationRequest = InvitationResponse

// LoginRequest captures credentials incoming from the client form elements
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterRequest alias LoginRequest as they both use identical struct
type RegisterRequest = LoginRequest

// UserResponse filters out sensitive fields when sending account profile to the client
type UserResponse struct {
	ID string `json:"id"`
	Username string `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}

type KanbanHandler struct {
	useCase domain.KanbanUseCase
	authUseCase domain.AuthUseCase
}

// MoveTaskPayload defines the strict JSON contract from the frontend
type MoveTaskPayload struct {
	TaskID string `json:"taskId"`
	TargetColumnID string `json:"targetColumnId"`
	TargetPosition int `json:"targetPosition"`
}

type CreateTaskPayload struct {
	ColumnID string `json:"columnId"`
	Title string `json:"title"`
	Description string `json:"description"`
}

type DeleteTaskPayload struct {
	ColumnID string `json:"columnId"`
	TaskID string `json:"taskId"`
	TaskPosition int `json:"taskPosition"`
}

type UpdateTaskPayload struct {
	TaskID string `json:"taskId"`
	Title string `json:"title"`
	Description string `json:"description"`
}

type ArchiveTaskPayload struct {
	ColumnID string `json:"columnId"`
	TaskID string `json:"taskId"`
	TaskPosition int `json:"taskPosition"`
}

// NewKanbanHandler initializes the delivery layer with its required business logic dependency.
func NewKanbanHandler(uc domain.KanbanUseCase, auc domain.AuthUseCase) *KanbanHandler {
	return &KanbanHandler{
		useCase: uc,
		authUseCase: auc,
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

func (h *KanbanHandler) HandleTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateTask(w, r)
	case http.MethodDelete:
		h.DeleteTask(w, r)
	case http.MethodPatch:
		h.UpdateTask(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *KanbanHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var payload CreateTaskPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Malformed JSON request body", http.StatusBadRequest)

		return
	}
	defer r.Body.Close()

	createdTask, err := h.useCase.CreateTask(r.Context(), payload.ColumnID, payload.Title, payload.Description)
	if err != nil {
		log.Printf("Error executing task creation workflow: %v", err)

		if strings.Contains(err.Error(), "business rule violation") {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		http.Error(w, "Internal server creation failure", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(createdTask); err != nil {
		http.Error(w, "Failed to encode response payload", http.StatusInternalServerError)

		return
	}
}

func (h *KanbanHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	var payload DeleteTaskPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Malformed JSON request body", http.StatusBadRequest)

		return
	}
	defer r.Body.Close()

	err := h.useCase.DeleteTask(r.Context(), payload.ColumnID, payload.TaskID, payload.TaskPosition)
	if err != nil {
		log.Printf("Error executing task deletion workflow: %v", err)

		if strings.Contains(err.Error(), "business rule violation") {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		http.Error(w, "Internal server update failure", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *KanbanHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var payload UpdateTaskPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Malformed JSON request body", http.StatusBadRequest)

		return
	}
	defer r.Body.Close()

	err := h.useCase.EditTask(r.Context(), payload.TaskID, payload.Title, payload.Description)
	if err != nil {
		log.Printf("Error executing task update workflow: %v", err)
	
		if strings.Contains(err.Error(), "business rule violation") {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		http.Error(w, "Internal server update failure", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *KanbanHandler) ArchiveTask(w http.ResponseWriter, r *http.Request) {
	// 1. Enforce strict HTTP method checking
	if r.Method != http.MethodPatch {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)

		return
	}

	var payload ArchiveTaskPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Malformed JSON request body", http.StatusBadRequest)

		return
	}
	defer r.Body.Close()

	err := h.useCase.ArchiveTask(r.Context(), payload.ColumnID, payload.TaskID, payload.TaskPosition)
	if err != nil {
		log.Printf("Error executing task archiving workflow: %v", err)

		if strings.Contains(err.Error(), "business rule violation") {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		http.Error(w, "Internal server update failure", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *KanbanHandler) UnarchiveTask(w http.ResponseWriter, r *http.Request) {
	// 1. Enforce strict HTTP method checking
	if r.Method != http.MethodPatch {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)

		return
	}

	var payload ArchiveTaskPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Malformed JSON request body", http.StatusBadRequest)

		return
	}
	defer r.Body.Close()

	err := h.useCase.UnarchiveTask(r.Context(), payload.ColumnID, payload.TaskID, payload.TaskPosition)
	if err != nil {
		log.Printf("Error executing task un-archiving workflow: %v", err)

		if strings.Contains(err.Error(), "business rule violation") {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		http.Error(w, "Internal server update failure", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *KanbanHandler) GetArchivedTasks(w http.ResponseWriter, r *http.Request) {
	// 1. Enforce strict HTTP method checking
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)

		return
	}

	boardID := r.URL.Query().Get("boardId")
	if boardID == "" {
		http.Error(w, "Missing required board 'boardId' query parameter", http.StatusBadRequest)

		return
	}

	archivedTasks, err := h.useCase.GetArchivedTasks(r.Context(), boardID)
	if err != nil {
		log.Printf("Error retrieving archived tasks: %v", err)

		if strings.Contains(err.Error(), "business rule violation") {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		http.Error(w, "Internal server retrieval failure", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(archivedTasks); err != nil {
		http.Error(w, "Failed to encode response payload", http.StatusInternalServerError)

		return
	}
}

// Login handles requests matching: POST /api/auth/login
func (h *KanbanHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Malformed JSON request body", http.StatusBadRequest)

		return
	}
	defer r.Body.Close()

	// 1. Process secure credential checks via UseCase interactor layer logic
	sessionID, err := h.authUseCase.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		log.Printf("Security validation failed for username %s: %v", req.Username, err)
		http.Error(w, "Invalid credential", http.StatusBadRequest)

		return
	}

	// 2. Set HTTP-Only Cookie wrapper to make it completely invisible to malicious JS (XSS protection)
	http.SetCookie(w, &http.Cookie{
		Name: "kanban_session",
		Value: sessionID,
		Path: "/",
		Expires: time.Now().Add(24 * time.Hour), // extends cookie session tracking window to 24 hours
		HttpOnly: true, // protects token strings from document.cookie queries
		Secure: false, // keep as false strictly for local localhost dev environment
		SameSite: http.SameSiteLaxMode, // guards system state from CSRF cross-origin attack vectors
	})

	w.WriteHeader(http.StatusNoContent)
}

// Logout handles requests matching DELETE /api/auth/logout
func (h *KanbanHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		return
	}

	// 1. Extract active session cookie
	cookie, err := r.Cookie("kanban_session")
	if err != nil  || cookie.Value == "" {
		http.Error(w, "Session cookie not found", http.StatusBadRequest)

		return
	}

	// 2. Remove from the sessions table, the session relevant to the sessionID / cookie value 
	sessionID := cookie.Value
	_ = h.authUseCase.Logout(r.Context(), sessionID)

	// Clear the browser cookie jar by forcing an immediate expiration context
	http.SetCookie(w, &http.Cookie{
		Name: "kanban_session",
		Value: "",
		Path: "/",
		Expires: time.Unix(0, 0), // Sets expiration to Jan 1, 1970, instantly destroying the cookie
		MaxAge: -1, // Ensures immediate eviction in modern browsers
		HttpOnly: true,
		Secure: false, // Keep matching with local dev setup flags
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusNoContent)
}

// Register handles requests matching: POST /api/auth/register
func (h *KanbanHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Malformed JSON request body", http.StatusBadRequest)

		return
	}
	defer r.Body.Close()

	_, err := h.authUseCase.Register(r.Context(), req.Username, req.Password)
	if err != nil {
		log.Printf("Registration failed for username %s: %v", req.Username, err)

		if strings.Contains(err.Error(), "business rule violation") {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		http.Error(w, "Unable to complete registration. Please check your submission constraints.", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusCreated)
}

// Authenticate handles requests matching: GET /api/auth/me
func (h *KanbanHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		return
	}

	// 1. Extract secure cookie
	cookie, err := r.Cookie("kanban_session")
	if err != nil {
		http.Error(w, "Unauthorized: No active session cookie found", http.StatusUnauthorized)

		return
	}

	user, err := h.authUseCase.AuthenticateSession(r.Context(), cookie.Value)
	if err != nil {
		log.Printf("Session validation failed: %v", err)
		http.Error(w, "Unauthorized session bounds", http.StatusUnauthorized)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UserResponse{
		ID: user.ID,
		Username: user.Username,
		CreatedAt: user.CreatedAt,
	})
}

func (h *KanbanHandler) ValidateInvitationToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		return
	}

	var req InvitationTokenValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Malformed JSON request body", http.StatusBadRequest)

		return
	}
	defer r.Body.Close()

	_, err := h.authUseCase.ValidateInvitationToken(r.Context(), req.Token)
	if err != nil {
		log.Printf("Invitation token validation for token: %s, failed: %v", req.Token, err)

		if strings.Contains(err.Error(), "business rule violation") {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		notFoundErrMessage := fmt.Sprintf("invitation token: %s, not found:", req.Token)
		if strings.Contains(err.Error(), notFoundErrMessage) {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		http.Error(w, "Unable to complete invitation token validation.", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *KanbanHandler) CreateInvitationToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		return
	}

	var req InvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Malformed JSON request body", http.StatusBadRequest)

		return
	}
	defer r.Body.Close()

	invitationToken, err := h.authUseCase.CreateInvitationToken(r.Context(), req.UserID, req.Email)
	if err != nil {
		log.Printf("Invitation generating by userId: %s, failed for email: %s: %v", req.UserID, req.Email, err)

		if strings.Contains(err.Error(), "business rule violation") {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		http.Error(w, "Unable to complete invitation process.", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(InvitationResponse{
		Token: invitationToken,
	})
}

func (h *KanbanHandler) ValidateEmailInput(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		return
	}

	var req EmailInputValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Malformed JSON request body", http.StatusBadRequest)

		return
	}
	defer r.Body.Close()

	_, _, err := h.authUseCase.ValidateEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "Invalid email input", http.StatusBadRequest)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}