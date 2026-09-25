package domain

import (
	"context"
	"time"
)

// User represents a core account profile entity inside our system layout
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Never expose the raw password hash strings over json payload wires
	Email        string    `json:"email"`
	Phone        *string   `json:"phone,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type Invitation struct {
	Token     string    `json:"token"`
	CreatedBy string    `json:"createdBy"`
	Email     string    `json:"email"`
	ExpiresAt time.Time `json:"expiresAt"`
	UsedAt    time.Time `json:"usedAt"`
}

// UserRepository defines the database contract required to resolve user identities
type UserRepository interface {
	CreateUser(ctx context.Context, username string, passwordHash string, email string, inviteToken string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByID(ctx context.Context, userID string) (*User, error)
	CreateSession(ctx context.Context, userID string) (*Session, error)
	GetSessionByID(ctx context.Context, sessionID string) (*Session, error)
	DeleteSessionByID(ctx context.Context, sessionID string) error
	IsEmailRegistered(ctx context.Context, email string) (bool, error)
	CreateAccountRegistrationInvitation(ctx context.Context, userID string, email string, expirationTime time.Duration) (string, error)
	GetAccountRegistrationInvitationByEmail(ctx context.Context, email string) (*Invitation, error)
	GetAccountRegistrationInvitationByToken(ctx context.Context, token string) (*Invitation, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

// AuthUserCase orchestrates credentials checks and login state verifications
type AuthUseCase interface {
	Register(ctx context.Context, username string, password string, email string, inviteToken string) (*User, error)
	Login(ctx context.Context, username string, password string) (string, error) // Returns a secure session ID string token
	AuthenticateSession(ctx context.Context, sessionID string) (*User, error)
	Logout(ctx context.Context, sessionID string) error
	CreateInvitationToken(ctx context.Context, userID string, email string, registerURL string) (string, error)
	ValidateInvitationToken(ctx context.Context, token string) (*Invitation, error)
	ValidateEmail(ctx context.Context, email string) (string, string, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	IsEmailRegistered(ctx context.Context, email string) (bool, error)
}
