// Package auth implements the authentication service (login, logout, session validation).
package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"awo.so/framework/internal/iam/domain"
	"awo.so/framework/internal/iam/session"
)

// ErrInvalidCredentials is returned when email/password do not match.
var ErrInvalidCredentials = errors.New("invalid credentials")

// ErrAccountLocked is returned when the account is temporarily locked.
var ErrAccountLocked = errors.New("account locked")

// ErrMFARequired signals that MFA verification is required before a full session is issued.
var ErrMFARequired = errors.New("mfa required")

// ErrInvalidSession is returned when a session token is invalid or expired.
var ErrInvalidSession = errors.New("invalid or expired session")

// Service handles credential verification, session creation, and logout.
type Service struct {
	db       *pgxpool.Pool
	sessions *session.Service
	log      *slog.Logger
}

// New creates an auth Service.
func New(db *pgxpool.Pool, sessions *session.Service, log *slog.Logger) *Service {
	return &Service{db: db, sessions: sessions, log: log}
}

// Login verifies credentials and returns a ResolvedSession + raw token.
// If MFA is enabled, returns ErrMFARequired and a pending token instead.
func (s *Service) Login(ctx context.Context, p domain.LoginParams) (*domain.ResolvedSession, string, error) {
	// Step 1: load user by email
	user, err := s.loadUserByEmail(ctx, p.Email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// Step 2: account lock check
	if user.IsLocked() {
		return nil, "", ErrAccountLocked
	}

	// Step 3: verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(p.Password)); err != nil {
		if updateErr := s.incrementFailedAttempts(ctx, user); updateErr != nil {
			s.log.Error("failed to increment failed attempts", slog.String("user_id", user.ID.String()), slog.Any("err", updateErr))
		}
		return nil, "", ErrInvalidCredentials
	}

	// Step 4: reset failed attempts
	if err := s.resetFailedAttempts(ctx, user.ID); err != nil {
		s.log.Warn("reset failed attempts error", slog.String("user_id", user.ID.String()), slog.Any("err", err))
	}

	// Step 5: MFA check — return pending token if MFA enabled
	if user.MFAEnabled {
		pendingToken, err := s.issueMFAPendingToken(ctx, user.ID)
		if err != nil {
			return nil, "", fmt.Errorf("mfa pending token: %w", err)
		}
		return nil, pendingToken, ErrMFARequired
	}

	// Step 6: build full session
	resolved, rawToken, err := s.sessions.Create(ctx, session.CreateParams{
		User:      user,
		IP:        p.IP,
		UserAgent: p.UserAgent,
	})
	if err != nil {
		return nil, "", fmt.Errorf("create session: %w", err)
	}
	return resolved, rawToken, nil
}

// ValidateSession looks up a session by raw token.
func (s *Service) ValidateSession(ctx context.Context, p domain.ValidateSessionParams) (*domain.ResolvedSession, error) {
	resolved, err := s.sessions.Get(ctx, p.RawToken)
	if err != nil || resolved == nil {
		return nil, ErrInvalidSession
	}
	return resolved, nil
}

// Logout invalidates the session identified by the raw token.
func (s *Service) Logout(ctx context.Context, p domain.LogoutParams) error {
	return s.sessions.Invalidate(ctx, p.RawToken)
}

// --- internal helpers ---

func (s *Service) loadUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	row := s.db.QueryRow(ctx, `
		SELECT id, tenant_id, email, full_name, password_hash,
		       user_type, mfa_enabled, mfa_secret,
		       is_active, is_suspended, is_super,
		       failed_login_attempts, locked_until, last_login_at,
		       created_at, updated_at
		FROM users
		WHERE email = $1 AND is_active = TRUE`, email)

	var u domain.User
	if err := row.Scan(
		&u.ID, &u.TenantID, &u.Email, &u.FullName, &u.PasswordHash,
		&u.UserType, &u.MFAEnabled, &u.MFASecret,
		&u.IsActive, &u.IsSuspended, &u.IsSuper,
		&u.FailedAttempts, &u.LockedUntil, &u.LastLoginAt,
		&u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("load user: %w", err)
	}
	return &u, nil
}

func (s *Service) incrementFailedAttempts(ctx context.Context, u *domain.User) error {
	const maxAttempts = 5
	const lockDuration = 15 * time.Minute

	newCount := u.FailedAttempts + 1
	var lockedUntil *time.Time
	if newCount >= maxAttempts {
		t := time.Now().Add(lockDuration)
		lockedUntil = &t
	}
	_, err := s.db.Exec(ctx, `
		UPDATE users
		SET failed_login_attempts = $1,
		    locked_until          = $2,
		    updated_at            = NOW()
		WHERE id = $3`,
		newCount, lockedUntil, u.ID)
	return err
}

func (s *Service) resetFailedAttempts(ctx context.Context, userID uuid.UUID) error {
	_, err := s.db.Exec(ctx, `
		UPDATE users
		SET failed_login_attempts = 0,
		    locked_until          = NULL,
		    last_login_at         = NOW(),
		    updated_at            = NOW()
		WHERE id = $1`, userID)
	return err
}

func (s *Service) issueMFAPendingToken(ctx context.Context, userID uuid.UUID) (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	token := hex.EncodeToString(raw)
	// Store in Redis with 5-min TTL via session service
	if err := s.sessions.StoreMFAPending(ctx, token, userID); err != nil {
		return "", err
	}
	return token, nil
}

// HashPassword bcrypt-hashes a plaintext password.
func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// SHA256Hex returns the lower-hex SHA-256 hash of s.
func SHA256Hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}
