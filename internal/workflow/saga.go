package workflow

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/workflow"
)

// Step is a forward action paired with its compensation.
type Step struct {
	Name       string
	Action     func(ctx workflow.Context) error
	Compensate func(ctx workflow.Context) error
}

// Saga manages a sequence of steps with automatic compensation on failure.
// Usage inside a Temporal workflow function:
//
//	s := workflow.NewSaga()
//	s.AddStep(reserveInventoryStep)
//	s.AddStep(postGLStep)
//	if err := s.Execute(ctx); err != nil {
//	    return err  // compensations already run
//	}
type Saga struct {
	steps []Step
}

// NewSaga creates an empty Saga.
func NewSaga() *Saga { return &Saga{} }

// AddStep appends a step. Steps execute in addition order; compensations run in reverse.
func (s *Saga) AddStep(step Step) *Saga {
	s.steps = append(s.steps, step)
	return s
}

// Execute runs steps in order. On first failure, runs compensations for all
// previously succeeded steps in reverse order. Returns the original error
// (compensation errors are logged but do not mask the original).
func (s *Saga) Execute(ctx workflow.Context) error {
	succeeded := make([]Step, 0, len(s.steps))

	for _, step := range s.steps {
		if err := step.Action(ctx); err != nil {
			compErr := s.compensate(ctx, succeeded)
			if compErr != nil {
				// Wrap both errors so the caller can see compensation failures.
				return fmt.Errorf("saga step %q failed: %w (compensation error: %v)", step.Name, err, compErr)
			}
			return fmt.Errorf("saga step %q failed: %w", step.Name, err)
		}
		succeeded = append(succeeded, step)
	}
	return nil
}

func (s *Saga) compensate(ctx workflow.Context, succeeded []Step) error {
	var firstErr error
	for i := len(succeeded) - 1; i >= 0; i-- {
		step := succeeded[i]
		if step.Compensate == nil {
			continue
		}
		if err := step.Compensate(ctx); err != nil {
			if firstErr == nil {
				firstErr = fmt.Errorf("compensation %q: %w", step.Name, err)
			}
			// Continue running remaining compensations even if one fails.
		}
	}
	return firstErr
}

// ActivityStep is a convenience constructor for saga steps that delegate to
// Temporal activities (the common case).
func ActivityStep(name string, actFn any, args ...any) Step {
	return Step{
		Name: name,
		Action: func(ctx workflow.Context) error {
			return workflow.ExecuteActivity(ctx, actFn, args...).Get(ctx, nil)
		},
	}
}

// GoContextSaga is a non-Temporal saga for use outside workflow functions
// (e.g. in HTTP handlers for multi-step provisioning operations).
type GoContextSaga struct {
	steps []goStep
}

type goStep struct {
	name       string
	action     func(ctx context.Context) error
	compensate func(ctx context.Context) error
}

// NewGoSaga creates a saga for use with a plain context.Context.
func NewGoSaga() *GoContextSaga { return &GoContextSaga{} }

// AddStep appends a step to the Go saga.
func (s *GoContextSaga) AddStep(name string, action, compensate func(ctx context.Context) error) *GoContextSaga {
	s.steps = append(s.steps, goStep{name: name, action: action, compensate: compensate})
	return s
}

// Execute runs steps in order; compensates in reverse on failure.
func (s *GoContextSaga) Execute(ctx context.Context) error {
	succeeded := make([]goStep, 0, len(s.steps))
	for _, step := range s.steps {
		if err := step.action(ctx); err != nil {
			_ = s.compensate(ctx, succeeded)
			return fmt.Errorf("saga step %q: %w", step.name, err)
		}
		succeeded = append(succeeded, step)
	}
	return nil
}

func (s *GoContextSaga) compensate(ctx context.Context, succeeded []goStep) error {
	var firstErr error
	for i := len(succeeded) - 1; i >= 0; i-- {
		step := succeeded[i]
		if step.compensate == nil {
			continue
		}
		if err := step.compensate(ctx); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
