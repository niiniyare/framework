// Package workflow defines the types for binding Temporal workflows to EntityDefinitions.
package workflow

// TriggerEvent names the lifecycle point at which a workflow is started.
type TriggerEvent string

const (
	TriggerAfterSave TriggerEvent = "after_save"
	TriggerOnSubmit  TriggerEvent = "on_submit"
	TriggerOnCancel  TriggerEvent = "on_cancel"
)

// WorkflowBinding attaches a Temporal workflow to an EntityDefinition lifecycle event.
type WorkflowBinding struct {
	// WorkflowName is the name registered with the Temporal worker.
	WorkflowName string

	// Trigger controls when the workflow is started.
	Trigger TriggerEvent

	// Condition is evaluated against the saved record before dispatching.
	// Nil means the workflow fires unconditionally on the trigger event.
	// The record is passed as map[string]any for decoupling.
	Condition func(record map[string]any) bool

	// TaskQueue overrides the default task queue for this workflow.
	// Empty string uses the framework default.
	TaskQueue string

	// IDTemplate is a Go text/template string for the Temporal workflow ID.
	// Available variables: .TenantID, .EntityName, .RecordID, .Event.
	// Default: "{{.TenantID}}.{{.EntityName}}.{{.RecordID}}.{{.WorkflowName}}"
	IDTemplate string
}
