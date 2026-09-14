package domain

import (
	"context"
	"time"
)

// User represents a core account profile entity inside our system layout
type User struct {
	ID string `json:"id"`
	Username string `json:"username"`
	PasswordHash string `json:"-"` // Never expose the raw password hash strings over json payload wires
	CreatedAt time.Time `json:"createdAt"`
}

type Session struct {
	ID string `json:"id"`
	UserID string `json:"userId"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// UserRepository defines the database contract required to resolve user identities
type UserRepository interface {
	CreateUser(ctx context.Context, username string, passwordHash string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByID(ctx context.Context, userID string) (*User, error)
	CreateSession(ctx context.Context, userID string) (*Session, error)
	GetSessionByID(ctx context.Context, sessionID string) (*Session, error)
}

// AuthUserCase orchestrates credentials checks and login state verifications
type AuthUseCase interface {
	Register(ctx context.Context, username string, password string) (*User, error)
	Login(ctx context.Context, username string, password string) (string, error) // Returns a secure session ID string token
	AuthenticateSession(ctx context.Context, sessionID string) (*User, error)
}