package workflow

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"text/template"
	"time"

	"go.temporal.io/sdk/client"

	"awo.so/framework/internal/core"
	wftypes "awo.so/framework/pkg/workflow"
)

const defaultIDTemplate = "{{.TenantID}}.{{.EntityName}}.{{.RecordID}}.{{.WorkflowName}}"

// idTemplateVars is passed to the workflow ID template.
type idTemplateVars struct {
	TenantID     string
	EntityName   string
	RecordID     string
	WorkflowName string
	Event        string
}

// Dispatcher fires Temporal workflows from entity lifecycle events.
// It is called from the after_save / on_submit handler after the DB commit.
type Dispatcher struct {
	client *Client
	log    *slog.Logger
}

// NewDispatcher creates a Dispatcher.
func NewDispatcher(c *Client, log *slog.Logger) *Dispatcher {
	return &Dispatcher{client: c, log: log}
}

// Dispatch starts a Temporal workflow for a binding if its condition is satisfied.
// It is non-blocking — the workflow is started in a goroutine with a short deadline
// so it never blocks the HTTP response.
func (d *Dispatcher) Dispatch(
	ctx context.Context,
	binding wftypes.WorkflowBinding,
	def *core.EntityDefinition,
	record *core.EntityRecord,
) {
	if binding.Condition != nil && !binding.Condition(record.Data) {
		return
	}

	go func() {
		dispatchCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		wfID, err := buildWorkflowID(binding, record, string(binding.Trigger))
		if err != nil {
			d.log.Error(
				"workflow: build ID failed",
				slog.String("workflow", binding.WorkflowName),
				slog.Any("err", err),
			)
			return
		}

		taskQueue := binding.TaskQueue
		if taskQueue == "" {
			taskQueue = fmt.Sprintf("%s.%s", def.Name, binding.WorkflowName)
		}

		opts := client.StartWorkflowOptions{
			ID:        wfID,
			TaskQueue: taskQueue,
			// Default reuse policy: allow duplicate failed only (Temporal default).
		}

		input := map[string]any{
			"tenant_id":   record.TenantID,
			"entity_name": record.EntityName,
			"record_id":   record.ID,
			"trigger":     string(binding.Trigger),
		}

		run, err := d.client.ExecuteWorkflow(dispatchCtx, opts, binding.WorkflowName, input)
		if err != nil {
			d.log.Error(
				"workflow: start failed",
				slog.String("workflow", binding.WorkflowName),
				slog.String("workflow_id", wfID),
				slog.Any("err", err),
			)
			return
		}

		d.log.Info(
			"workflow: started",
			slog.String("workflow", binding.WorkflowName),
			slog.String("workflow_id", wfID),
			slog.String("run_id", run.GetRunID()),
		)
	}()
}

// DispatchAll fires all matching workflow bindings for a lifecycle event.
func (d *Dispatcher) DispatchAll(
	ctx context.Context,
	def *core.EntityDefinition,
	trigger wftypes.TriggerEvent,
	record *core.EntityRecord,
) {
	for _, binding := range def.Workflows {
		if binding.Trigger == trigger {
			d.Dispatch(ctx, binding, def, record)
		}
	}
}

func buildWorkflowID(b wftypes.WorkflowBinding, r *core.EntityRecord, event string) (string, error) {
	tmplStr := b.IDTemplate
	if tmplStr == "" {
		tmplStr = defaultIDTemplate
	}
	tmpl, err := template.New("wfid").Parse(tmplStr)
	if err != nil {
		return "", err
	}
	vars := idTemplateVars{
		TenantID:     r.TenantID,
		EntityName:   r.EntityName,
		RecordID:     r.ID,
		WorkflowName: b.WorkflowName,
		Event:        event,
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vars); err != nil {
		return "", err
	}
	return buf.String(), nil
}
