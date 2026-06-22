// Package domain contains IAM domain types shared across IAM sub-packages.
package domain

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// UserType constants.
const (
	UserTypeInternal  = "INTERNAL"
	UserTypeSysAdmin  = "SYSADMIN"
	UserTypeCustomer  = "CUSTOMER"
	UserTypePortal    = "PORTAL"
	UserTypeAPI       = "API"
)

// EntityScopeType constrains which entity records a session can access.
type EntityScopeType string

const (
	EntityScopeAll     EntityScopeType = "all"
	EntityScopeSubtree EntityScopeType = "subtree"
	EntityScopeEntity  EntityScopeType = "entity"
)

// EntityScope is stored as JSONB in user_sessions.entity_scope.
type EntityScope struct {
	Type       EntityScopeType `json:"type"`
	EntityID   string          `json:"entity_id,omitempty"`
	PathPrefix string          `json:"path_prefix,omitempty"` // ltree path for subtree
}

// Configuration holds pre-computed flags, settings and user preferences
// embedded in each session. Stored as JSONB in user_sessions.configuration.
type Configuration struct {
	Flags    map[string]bool   `json:"flags"`
	Settings map[string]string `json:"settings"`
	Prefs    map[string]string `json:"prefs"`
}

// ResolvedSession is the fully-hydrated session object. It is computed once
// at login and cached in Redis for the session lifetime. It is the authoritative
// identity context for every request — NOT a Principal.
type ResolvedSession struct {
	UserID        uuid.UUID
	UserType      string
	TenantID      uuid.UUID
	PrincipalID   *uuid.UUID // non-nil for portal users only
	DisplayName   string
	EntityScope   EntityScope
	Configuration Configuration
	ExpiresAt     time.Time
	SessionToken  string // SHA-256 hex — for session invalidation only
}

// FeatureEnabled returns true if flagKey is enabled and all parent module flags
// in the dot-separated hierarchy are also enabled.
// finance.transactions → checks flags["finance"] AND flags["finance.transactions"].
func (s *ResolvedSession) FeatureEnabled(flagKey string) bool {
	if !s.Configuration.Flags[flagKey] {
		return false
	}
	if idx := strings.Index(flagKey, "."); idx > 0 {
		moduleKey := flagKey[:idx]
		if !s.Configuration.Flags[moduleKey] {
			return false
		}
	}
	return true
}

// SettingString returns the effective setting string value or def.
func (s *ResolvedSession) SettingString(key, def string) string {
	if v, ok := s.Configuration.Settings[key]; ok {
		return v
	}
	return def
}

// SettingBool parses the effective setting as a boolean or returns def.
func (s *ResolvedSession) SettingBool(key string, def bool) bool {
	if v, ok := s.Configuration.Settings[key]; ok {
		b, err := strconv.ParseBool(v)
		if err == nil {
			return b
		}
	}
	return def
}

// SettingInt parses the effective setting as an int or returns def.
func (s *ResolvedSession) SettingInt(key string, def int) int {
	if v, ok := s.Configuration.Settings[key]; ok {
		i, err := strconv.Atoi(v)
		if err == nil {
			return i
		}
	}
	return def
}

// IsPlatform returns true for platform staff (SYSADMIN user type).
func (s *ResolvedSession) IsPlatform() bool {
	return s.UserType == UserTypeSysAdmin
}

// IsPortal returns true for portal (customer-facing) users.
func (s *ResolvedSession) IsPortal() bool {
	return s.UserType == UserTypePortal
}

// IsAPI returns true for API key sessions.
func (s *ResolvedSession) IsAPI() bool {
	return s.UserType == UserTypeAPI
}
