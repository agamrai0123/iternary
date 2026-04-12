package auth

import (
	"time"
)

// User represents a user in the system
type AuthUser struct {
	ID        string    `json:"id"`
	Username  string    `json:"username" binding:"required,min=3,max=50"`
	Email     string    `json:"email" binding:"required,email"`
	Password  string    `json:"password,omitempty" binding:"required,min=6"`
	FullName  string    `json:"full_name"`
	Bio       string    `json:"bio"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Session represents a user session
type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token     string    `json:"token"`
	User      *AuthUser `json:"user"`
	ExpiresAt time.Time `json:"expires_at"`
}

// SignupRequest represents a signup request
type SignupRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name"`
}

// ProfileUpdateRequest represents a profile update request
type ProfileUpdateRequest struct {
	FullName string `json:"full_name"`
	Bio      string `json:"bio"`
	Avatar   string `json:"avatar"`
}
