package usecase

import (
	"context"
	"fmt"
	"kanban-backend/internal/domain"
	"kanban-backend/pkg/util"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthInteractor struct {
	userRepo domain.UserRepository
}

// NewAuthInteractor initializes our security core with its required repository port dependency
func NewAuthInteractor(repo domain.UserRepository) *AuthInteractor {
	return &AuthInteractor{userRepo: repo}
}

// Login verifies incoming string passwords against secure, database-store salted bcrypt hashes
func (uc *AuthInteractor) Login(ctx context.Context, username string, password string) (string, error) {
	if username == "" || password == "" {
		return "", fmt.Errorf("business rule violation: credentials cannot be left blank")
	}

	// 1. Extract the target profile context block from the database row rows
	user, err := uc.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("invalid login: %w", err)
	}

	// 2. Cryptographic comparison, comparing incoming text bytes with salted hashes safely
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid login")
	}

	// 3. Generate session token: For modern token-less architecture layout baseline, we will output a high entropy tracking string to identify this session
	session, err := uc.userRepo.CreateSession(ctx, user.ID)
	if err != nil {
		return "", fmt.Errorf("usecase failed to persist a new session: %w", err)
	}
	return session.ID, nil
}

// AuthenticateSession reads session identifiers to confirm active account authorization properties
func (uc *AuthInteractor) AuthenticateSession(ctx context.Context, sessionID string) (*domain.User, error) {
	if sessionID == "" {
		return nil, fmt.Errorf("business rule violation: active session context is required")
	}

	// Grab registered session by sessionID
	session, err := uc.userRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("unauthorized session bounds: invalid or missing session token: %w", err)
	}

	// validate whether the session has expired or not
	now := time.Now().UTC()
	if now.After(session.ExpiresAt) {
		return nil, fmt.Errorf("unauthorized session bounds: login session has expired")
	}

	user, err := uc.userRepo.GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("unauthorized session bounds: %w", err)
	}

	return user, nil
}

func (uc *AuthInteractor) Register(ctx context.Context, username string, password string) (*domain.User, error) {
	if username == "" || password == "" {
		return nil, fmt.Errorf("business rule violation: credentials cannot be left blank")
	}

	if len(password) < 8 {
		return nil, fmt.Errorf("business rule violation: password must be at least 8 characters long")
	}

	// encrypt password
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("cryptographic failure: failed to securely salt user password: %w", err)
	}

	// pass the hashed byte representation of the password to the repo method to create a user
	user, err := uc.userRepo.CreateUser(ctx, username, string(bytes))
	if err != nil {
		return nil, fmt.Errorf("usecase failed to persist a new account: %w", err)
	}

	return user, nil
}

func (uc *AuthInteractor) Logout(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return fmt.Errorf("business rule violation: active session context is required")
	}

	err := uc.userRepo.DeleteSessionByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("usecase failed to delete a session: %w", err)
	}

	return nil
}

func (uc *AuthInteractor) CreateInvitationToken(ctx context.Context, userID string, email string) (string, error) {
	if userID == "" {

		return "", fmt.Errorf("business rule violation: user id context is required")
	}

	if email == "" {

		return "", fmt.Errorf("business rule violation: email context is required")
	}

	expirationTime := 7 * 24 * time.Hour
	invitationToken, err := uc.userRepo.CreateAccountRegistrationInvitation(ctx, userID, email, expirationTime)
	if err != nil {

		return "", fmt.Errorf("usecase failed to persist a new invitation: %w", err)
	}

	return invitationToken, nil
}

func (uc *AuthInteractor) ValidateInvitationToken(ctx context.Context, token string) (*domain.Invitation, error) {
	if token == "" {

		return nil, fmt.Errorf("business rule violation: invitation token context is required")
	}

	invitation, err := uc.userRepo.GetAccountRegistrationInvitationByToken(ctx, token)
	if err != nil {
		
		return nil, fmt.Errorf("usecase failed to extract invitation for the token: %s: %w", token, err)
	}

	now := time.Now().UTC() // Enforce UTC because the sql table CURRENT_TIMESTAMP default to UTC timezone
	if invitation.ExpiresAt.Before(now) {

		return nil, fmt.Errorf("business rule violation: invitation token expired")
	}

	// enforce single usage for a token rule.
	if !invitation.UsedAt.IsZero() {

		return nil, fmt.Errorf("business rule violation: invitation token has been used")
	}

	return invitation, nil
}

func (uc *AuthInteractor) ValidateEmail(ctx context.Context, email string) (string, string, error) {
	address, domain, err := util.ValidateEmailInput(email)
	if err != nil {
		return "", "", fmt.Errorf("usecase failed to validate email: %w", err)
	}

	return address, domain, nil

}