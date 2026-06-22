package domain

import (
	"time"

	"github.com/google/uuid"
)

// User is the domain representation of a platform/tenant user.
type User struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	Email          string
	FullName       string
	PasswordHash   string
	UserType       string
	MFAEnabled     bool
	MFASecret      string // AES-256-GCM encrypted TOTP secret; empty if MFA disabled
	IsActive       bool
	IsSuspended    bool
	IsSuper        bool
	FailedAttempts int
	LockedUntil    *time.Time
	LastLoginAt    *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// IsLocked returns true if the account is temporarily locked due to failed attempts.
func (u *User) IsLocked() bool {
	return u.LockedUntil != nil && u.LockedUntil.After(time.Now())
}

// LoginParams are passed to AuthService.Login.
type LoginParams struct {
	Email     string
	Password  string
	IP        string
	UserAgent string
}

// ValidateMFAParams are passed to AuthService.ValidateMFA.
type ValidateMFAParams struct {
	PendingToken string
	Code         string
	IP           string
	UserAgent    string
}

// ValidateSessionParams are passed to AuthService.ValidateSession.
type ValidateSessionParams struct {
	RawToken string
}

// LogoutParams are passed to AuthService.Logout.
type LogoutParams struct {
	RawToken string
	UserID   uuid.UUID
}

// PasswordResetRequestParams initiates a password reset.
type PasswordResetRequestParams struct {
	Email string
}

// PasswordResetParams completes the password reset with a token.
type PasswordResetParams struct {
	Token    string
	Password string
}

// APIKey is the domain representation of an API key.
type APIKey struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	Name      string
	KeyHash   string // SHA-256 hex; raw key never stored
	Scopes    []string
	CreatedBy *uuid.UUID
	ExpiresAt *time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

// IsValid returns false if the key is revoked or expired.
func (k *APIKey) IsValid() bool {
	if k.RevokedAt != nil {
		return false
	}
	if k.ExpiresAt != nil && k.ExpiresAt.Before(time.Now()) {
		return false
	}
	return true
}

// CreateAPIKeyParams creates a new API key.
type CreateAPIKeyParams struct {
	TenantID  uuid.UUID
	Name      string
	Scopes    []string
	ExpiresAt *time.Time
	CreatedBy uuid.UUID
}
