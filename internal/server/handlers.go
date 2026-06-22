package server

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"awo.so/framework/internal/authz"
	"awo.so/framework/internal/core"
	internalHooks "awo.so/framework/internal/hooks"
	"awo.so/framework/internal/middleware"
	internalPerms "awo.so/framework/internal/permissions"
	"awo.so/framework/internal/sdui"
	"awo.so/framework/internal/store"
	"awo.so/framework/internal/tenant"
	"awo.so/framework/internal/workflow"
	"awo.so/framework/pkg/filter"
	"awo.so/framework/pkg/hooks"
	"awo.so/framework/pkg/permissions"
	wftypes "awo.so/framework/pkg/workflow"
)

// Deps holds the dependencies injected into every handler.
type Deps struct {
	SystemRegistry *core.EntityRegistry
	TenantRegistry *tenant.Registry
	HookExecutor   *internalHooks.Executor
	Evaluator      *internalPerms.Evaluator
	Authz          *authz.Service     // may be nil when DB is not configured; used for Casbin RBAC
	Dispatcher     *workflow.Dispatcher // may be nil when Temporal is not configured
	PageBuilder    *sdui.CachedBuilder  // may be nil when Redis is not configured
	RepoFor        func(t *tenant.Tenant) store.EntityRepository
	Log            *slog.Logger
}

// checkAccess enforces RBAC for (principal, entity, action).
// When Authz (Casbin) is available it is the authoritative check.
// Superusers bypass Casbin but not privacy policies.
// Falls back to Evaluator (entity-declared permissions) when Authz is nil.
func checkAccess(
	ctx context.Context,
	d *Deps,
	def *core.EntityDefinition,
	principal *permissions.Principal,
	action permissions.Action,
	record *core.EntityRecord,
) error {
	if principal.IsSuper {
		return nil
	}

	if d.Authz != nil {
		subject := authz.TenantSubject(principal.UserID)
		domain := authz.TenantDomain(principal.TenantID)
		ok, err := d.Authz.Enforce(ctx, authz.Request{
			Subject: subject,
			Domain:  domain,
			Object:  def.Name,
			Action:  string(action),
		})
		if err != nil {
			return internalPerms.ErrForbidden
		}
		if !ok {
			return &internalPerms.ForbiddenError{
				EntityName: def.Name,
				Action:     action,
				Reason:     "casbin: access denied",
			}
		}
		return nil
	}

	// Casbin not available — fall back to entity-declared permissions.
	return d.Evaluator.Check(ctx, def, principal, action, record)
}

// resolve looks up the EntityDefinition for the entity named in the URL parameter.
func resolve(c *fiber.Ctx, d *Deps) (*core.EntityDefinition, *tenant.Tenant, *permissions.Principal, error) {
	t := tenant.FromCtx(c.UserContext())
	if t == nil {
		return nil, nil, nil, errNotFound(c, "tenant not found")
	}
	principal := middleware.PrincipalFromCtx(c.UserContext())
	if principal == nil {
		return nil, nil, nil, errForbidden(c, "unauthorized")
	}

	entityName := c.Params("entity")
	entry, err := d.TenantRegistry.GetOrLoad(t.ID.String(), func() (*tenant.Entry, error) {
		return &tenant.Entry{Tenant: t, Registry: core.NewEntityRegistry()}, nil
	})
	if err != nil {
		return nil, nil, nil, errInternal(c)
	}
	resolver := core.NewEntityResolver(d.SystemRegistry, entry.Registry)
	def, err := resolver.Resolve(entityName)
	if err != nil {
		return nil, nil, nil, errNotFound(c, "entity not found")
	}
	return def, t, principal, nil
}

// HandleList handles GET /api/:entity
func HandleList(d *Deps) fiber.Handler {
	return func(c *fiber.Ctx) error {
		def, t, principal, handlerErr := resolve(c, d)
		if handlerErr != nil {
			return handlerErr
		}

		if err := checkAccess(c.UserContext(), d, def, principal, permissions.ActionRead, nil); err != nil {
			return errForbidden(c, err.Error())
		}

		f := filter.New().Limit(pageSize(c)).Build()
		if err := d.Evaluator.ApplyPrivacy(c.UserContext(), def, principal, &f); err != nil {
			return errInternal(c)
		}

		repo := d.RepoFor(t)
		records, total, err := repo.FindMany(c.UserContext(), def.Name, f)
		if err != nil {
			d.Log.Error("list query failed", slog.String("entity", def.Name), slog.Any("err", err))
			return errInternal(c)
		}

		return okList(c, records, PageMeta{Total: total, Limit: f.Limit})
	}
}

// HandleGet handles GET /api/:entity/:id
func HandleGet(d *Deps) fiber.Handler {
	return func(c *fiber.Ctx) error {
		def, t, principal, handlerErr := resolve(c, d)
		if handlerErr != nil {
			return handlerErr
		}

		repo := d.RepoFor(t)
		record, err := repo.FindByID(c.UserContext(), def.Name, c.Params("id"))
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				return errNotFound(c, "record not found")
			}
			return errInternal(c)
		}

		if err := checkAccess(c.UserContext(), d, def, principal, permissions.ActionRead, record); err != nil {
			return errForbidden(c, err.Error())
		}

		return ok(c, record)
	}
}

// HandleCreate handles POST /api/:entity
func HandleCreate(d *Deps) fiber.Handler {
	return func(c *fiber.Ctx) error {
		def, t, principal, handlerErr := resolve(c, d)
		if handlerErr != nil {
			return handlerErr
		}

		if err := checkAccess(c.UserContext(), d, def, principal, permissions.ActionCreate, nil); err != nil {
			return errForbidden(c, err.Error())
		}

		var data map[string]any
		if err := c.BodyParser(&data); err != nil {
			return errBadRequest(c, "invalid JSON body")
		}

		repo := d.RepoFor(t)

		// before_validate hook
		if err := d.HookExecutor.Run(c.UserContext(), def, hooks.BeforeValidate, nil, nil, data, principal, repo); err != nil {
			return errUnprocessable(c, []FieldError{{Code: "hook_error", Message: err.Error()}})
		}

		// before_save hook
		if err := d.HookExecutor.Run(c.UserContext(), def, hooks.BeforeSave, nil, nil, data, principal, repo); err != nil {
			return errUnprocessable(c, []FieldError{{Code: "hook_error", Message: err.Error()}})
		}

		record, err := repo.Create(c.UserContext(), def.Name, data)
		if err != nil {
			if errors.Is(err, store.ErrConflict) {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "record already exists"})
			}
			d.Log.Error("create failed", slog.String("entity", def.Name), slog.Any("err", err))
			return errInternal(c)
		}

		// after_save hook
		if err := d.HookExecutor.Run(c.UserContext(), def, hooks.AfterSave, record, nil, data, principal, repo); err != nil {
			d.Log.Error("after_save hook failed", slog.String("entity", def.Name), slog.Any("err", err))
		}

		// Workflow dispatch (non-blocking)
		if d.Dispatcher != nil {
			d.Dispatcher.DispatchAll(c.UserContext(), def, wftypes.TriggerAfterSave, record)
		}

		return okCreated(c, record)
	}
}

// HandleUpdate handles PUT /api/:entity/:id
func HandleUpdate(d *Deps) fiber.Handler {
	return func(c *fiber.Ctx) error {
		def, t, principal, handlerErr := resolve(c, d)
		if handlerErr != nil {
			return handlerErr
		}

		repo := d.RepoFor(t)
		existing, err := repo.FindByID(c.UserContext(), def.Name, c.Params("id"))
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				return errNotFound(c, "record not found")
			}
			return errInternal(c)
		}

		if err := checkAccess(c.UserContext(), d, def, principal, permissions.ActionUpdate, existing); err != nil {
			return errForbidden(c, err.Error())
		}

		var data map[string]any
		if err := c.BodyParser(&data); err != nil {
			return errBadRequest(c, "invalid JSON body")
		}

		if err := d.HookExecutor.Run(c.UserContext(), def, hooks.BeforeValidate, existing, existing.Clone(), data, principal, repo); err != nil {
			return errUnprocessable(c, []FieldError{{Code: "hook_error", Message: err.Error()}})
		}
		if err := d.HookExecutor.Run(c.UserContext(), def, hooks.BeforeSave, existing, existing.Clone(), data, principal, repo); err != nil {
			return errUnprocessable(c, []FieldError{{Code: "hook_error", Message: err.Error()}})
		}

		record, err := repo.Update(c.UserContext(), def.Name, c.Params("id"), data)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				return errNotFound(c, "record not found")
			}
			return errInternal(c)
		}

		if err := d.HookExecutor.Run(c.UserContext(), def, hooks.AfterSave, record, existing, data, principal, repo); err != nil {
			d.Log.Error("after_save hook failed", slog.String("entity", def.Name), slog.Any("err", err))
		}

		if d.Dispatcher != nil {
			d.Dispatcher.DispatchAll(c.UserContext(), def, wftypes.TriggerAfterSave, record)
		}

		return ok(c, record)
	}
}

// HandleDelete handles DELETE /api/:entity/:id
func HandleDelete(d *Deps) fiber.Handler {
	return func(c *fiber.Ctx) error {
		def, t, principal, handlerErr := resolve(c, d)
		if handlerErr != nil {
			return handlerErr
		}

		repo := d.RepoFor(t)
		existing, err := repo.FindByID(c.UserContext(), def.Name, c.Params("id"))
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				return errNotFound(c, "record not found")
			}
			return errInternal(c)
		}

		if err := checkAccess(c.UserContext(), d, def, principal, permissions.ActionDelete, existing); err != nil {
			return errForbidden(c, err.Error())
		}

		if err := d.HookExecutor.Run(c.UserContext(), def, hooks.BeforeDelete, existing, nil, nil, principal, repo); err != nil {
			return errUnprocessable(c, []FieldError{{Code: "hook_error", Message: err.Error()}})
		}

		if err := repo.Delete(c.UserContext(), def.Name, c.Params("id")); err != nil {
			return errInternal(c)
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

// HandleSubmit handles POST /api/:entity/:id/submit
func HandleSubmit(d *Deps) fiber.Handler {
	return func(c *fiber.Ctx) error {
		def, t, principal, handlerErr := resolve(c, d)
		if handlerErr != nil {
			return handlerErr
		}
		if !def.IsSubmittable {
			return errNotFound(c, "entity does not support submit")
		}

		repo := d.RepoFor(t)
		existing, err := repo.FindByID(c.UserContext(), def.Name, c.Params("id"))
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				return errNotFound(c, "record not found")
			}
			return errInternal(c)
		}

		if err := checkAccess(c.UserContext(), d, def, principal, permissions.ActionSubmit, existing); err != nil {
			return errForbidden(c, err.Error())
		}

		if err := d.HookExecutor.Run(c.UserContext(), def, hooks.OnSubmit, existing, nil, nil, principal, repo); err != nil {
			return errUnprocessable(c, []FieldError{{Code: "hook_error", Message: err.Error()}})
		}

		record, err := repo.Update(c.UserContext(), def.Name, existing.ID, map[string]any{"doc_status": string(core.DocStatusSubmitted)})
		if err != nil {
			return errInternal(c)
		}

		if d.Dispatcher != nil {
			d.Dispatcher.DispatchAll(c.UserContext(), def, wftypes.TriggerOnSubmit, record)
		}

		return ok(c, record)
	}
}

// HandleCancel handles POST /api/:entity/:id/cancel
func HandleCancel(d *Deps) fiber.Handler {
	return func(c *fiber.Ctx) error {
		def, t, principal, handlerErr := resolve(c, d)
		if handlerErr != nil {
			return handlerErr
		}
		if !def.IsSubmittable {
			return errNotFound(c, "entity does not support cancel")
		}

		repo := d.RepoFor(t)
		existing, err := repo.FindByID(c.UserContext(), def.Name, c.Params("id"))
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				return errNotFound(c, "record not found")
			}
			return errInternal(c)
		}

		if err := checkAccess(c.UserContext(), d, def, principal, permissions.ActionCancel, existing); err != nil {
			return errForbidden(c, err.Error())
		}

		if err := d.HookExecutor.Run(c.UserContext(), def, hooks.OnCancel, existing, nil, nil, principal, repo); err != nil {
			return errUnprocessable(c, []FieldError{{Code: "hook_error", Message: err.Error()}})
		}

		record, err := repo.Update(c.UserContext(), def.Name, existing.ID, map[string]any{"doc_status": string(core.DocStatusCancelled)})
		if err != nil {
			return errInternal(c)
		}
		return ok(c, record)
	}
}

// HandlePage handles GET /api/:entity/__page — returns amis page JSON.
func HandlePage(d *Deps) fiber.Handler {
	return func(c *fiber.Ctx) error {
		def, t, principal, handlerErr := resolve(c, d)
		if handlerErr != nil {
			return handlerErr
		}
		if err := checkAccess(c.UserContext(), d, def, principal, permissions.ActionRead, nil); err != nil {
			return errForbidden(c, err.Error())
		}

		role := "default"
		if len(principal.Roles) > 0 {
			role = principal.Roles[0]
		}

		if d.PageBuilder != nil {
			page, err := d.PageBuilder.GetOrBuild(c.UserContext(), t.ID.String(), def, role)
			if err != nil {
				return errInternal(c)
			}
			c.Set("Content-Type", "application/json")
			return c.Send(page)
		}

		// No Redis — build uncached.
		page, err := sdui.BuildPage(def)
		if err != nil {
			return errInternal(c)
		}
		c.Set("Content-Type", "application/json")
		return c.Send(page)
	}
}

// pageSize extracts the limit from query params with a safe default and ceiling.
func pageSize(c *fiber.Ctx) int {
	limit := c.QueryInt("limit", 20)
	if limit <= 0 {
		limit = 20
	}
	if limit > 200 {
		limit = 200
	}
	return limit
}
