// Package hr is the Awo built-in HR module.
// Provides: Employee, LeaveType, LeaveAllocation, LeaveRequest,
// AttendanceRecord, PayrollRun, Payslip, DisciplinaryWarning.
package hr

import (
	"awo.so/framework/internal/core"
	"awo.so/framework/modules/hr/entities"
)

// Register adds all HR module EntityDefinitions to the system registry.
func Register(reg *core.EntityRegistry) {
	reg.MustRegister(entities.EmployeeDefinition)
	reg.MustRegister(entities.LeaveTypeDefinition)
	reg.MustRegister(entities.LeaveAllocationDefinition)
	reg.MustRegister(entities.LeaveRequestDefinition)
	reg.MustRegister(entities.AttendanceRecordDefinition)
	reg.MustRegister(entities.PayrollRunDefinition)
	reg.MustRegister(entities.PayslipDefinition)
	reg.MustRegister(entities.DisciplinaryWarningDefinition)
}
