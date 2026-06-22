// Package authz provides Casbin-based RBAC enforcement for the Awo framework.
//
// Subject format:   "tenant:{userID}" | "platform:{userID}" | "api:{clientID}"
// Domain format:    "{tenantUUID}" | "_platform_" | "{tenantUUID}:api"
// Object format:    entity name or resource path (supports keyMatch2 wildcards)
// Action format:    "create" | "read" | "update" | "delete" | "submit" | "cancel" | "*"
package authz

import (
	"time"

	"github.com/google/uuid"
)

// Subject builder functions — enforce consistent naming conventions.
func PlatformSubject(userID string) string { return "platform:" + userID }
func TenantSubject(userID string) string   { return "tenant:" + userID }
func PortalSubject(userID string) string   { return "portal:" + userID }
func APISubject(clientID string) string    { return "api:" + clientID }

// Domain builder functions.
const DomainPlatform = "_platform_"

func TenantDomain(tenantID string) string { return tenantID }
func PortalDomain(tenantID string) string { return tenantID + ":portal" }
func APIDomain(tenantID string) string    { return tenantID + ":api" }

// Effect constants.
const (
	EffectAllow = "allow"
	EffectDeny  = "deny"
)

// Request is a single Casbin enforcement request.
type Request struct {
	Subject string
	Domain  string
	Object  string
	Action  string
}

// Policy is a Casbin p-rule.
type Policy struct {
	Subject string
	Domain  string
	Object  string
	Action  string
	Effect  string // "allow" | "deny"
}

// RoleAssignment mirrors the role_assignments table.
type RoleAssignment struct {
	ID          uuid.UUID
	TenantID    *uuid.UUID
	Subject     string
	RoleName    string
	Domain      string
	AssignedBy  string
	DelegatedBy string
	ExpiresAt   *time.Time
	IsActive    bool
	CreatedAt   time.Time
}

// AssignOptions holds optional parameters for AssignRole.
type AssignOptions struct {
	ExpiresAt   *time.Time
	AssignedBy  string
	DelegatedBy string
}

// AssignOpt is a functional option for AssignRole.
type AssignOpt func(*AssignOptions)

func WithExpiry(t time.Time) AssignOpt {
	return func(o *AssignOptions) { o.ExpiresAt = &t }
}

func WithAssignedBy(subject string) AssignOpt {
	return func(o *AssignOptions) { o.AssignedBy = subject }
}

func WithDelegatedBy(subject string) AssignOpt {
	return func(o *AssignOptions) { o.DelegatedBy = subject }
}

func applyOpts(opts []AssignOpt) AssignOptions {
	o := AssignOptions{}
	for _, fn := range opts {
		fn(&o)
	}
	return o
}
