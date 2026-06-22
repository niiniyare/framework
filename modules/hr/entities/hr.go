// Package entities contains all HR module EntityDefinitions.
package entities

import (
	"context"
	"fmt"

	"awo.so/framework/pkg/entitydef"
	"awo.so/framework/pkg/fieldtype"
	"awo.so/framework/pkg/hooks"
	"awo.so/framework/pkg/permissions"
)

// EmployeeDefinition represents a person employed by the organisation.
var EmployeeDefinition = entitydef.New("Employee").
	System().
	Label("Employee").
	Field("employee_number", fieldtype.Data,
		fieldtype.ReadOnly(), fieldtype.WithNamingSeries("EMP-{YYYY}-{SEQ:5}")).
	Field("full_name", fieldtype.Data, fieldtype.Required(), fieldtype.MaxLen(255)).
	Field("date_of_birth", fieldtype.Date).
	Field("gender", fieldtype.Select,
		fieldtype.Choices("Male", "Female", "Other", "PreferNotToSay")).
	Field("national_id", fieldtype.Data, fieldtype.Unique(), fieldtype.MaxLen(50)).
	Field("kra_pin", fieldtype.Data, fieldtype.Sensitive(), fieldtype.MaxLen(20)).
	Field("nssf_number", fieldtype.Data, fieldtype.MaxLen(20)).
	Field("nhif_number", fieldtype.Data, fieldtype.MaxLen(20)).
	Field("department", fieldtype.Link, fieldtype.LinkTo("Department")).
	Field("designation", fieldtype.Data, fieldtype.MaxLen(100)).
	Field("employment_type", fieldtype.Select,
		fieldtype.Required(),
		fieldtype.Default("FullTime"),
		fieldtype.Choices("FullTime", "PartTime", "Contract", "Intern")).
	Field("date_of_joining", fieldtype.Date, fieldtype.Required()).
	Field("date_of_leaving", fieldtype.Date).
	Field("reports_to", fieldtype.Link, fieldtype.LinkTo("Employee")).
	Field("bank_account", fieldtype.Data, fieldtype.Sensitive(), fieldtype.MaxLen(30)).
	Field("bank_name", fieldtype.Data, fieldtype.MaxLen(100)).
	Field("basic_salary", fieldtype.Currency, fieldtype.Default(float64(0))).
	Field("email", fieldtype.Data, fieldtype.MaxLen(255)).
	Field("phone", fieldtype.Data, fieldtype.MaxLen(30)).
	Field("image", fieldtype.AttachImage).
	Field("is_active", fieldtype.Bool, fieldtype.Default(true)).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("hr_manager", permissions.AllActions...)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

// LeaveTypeDefinition defines a category of leave (e.g. Annual, Sick).
var LeaveTypeDefinition = entitydef.New("LeaveType").
	System().
	Label("Leave Type").
	Field("name", fieldtype.Data, fieldtype.Required(), fieldtype.Unique(), fieldtype.MaxLen(100)).
	Field("max_days_allowed", fieldtype.Float, fieldtype.Required(), fieldtype.Default(float64(0))).
	Field("is_carry_forward", fieldtype.Bool, fieldtype.Default(false)).
	Field("max_carry_forward_days", fieldtype.Float, fieldtype.Default(float64(0))).
	Field("is_earned_leave", fieldtype.Bool, fieldtype.Default(false)).
	Field("description", fieldtype.SmallText).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("hr_manager", permissions.AllActions...)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

// LeaveAllocationDefinition grants leave days to an employee for a period.
var LeaveAllocationDefinition = entitydef.New("LeaveAllocation").
	System().
	Label("Leave Allocation").
	Field("employee", fieldtype.Link,
		fieldtype.Required(), fieldtype.LinkTo("Employee"), fieldtype.Immutable()).
	Field("leave_type", fieldtype.Link,
		fieldtype.Required(), fieldtype.LinkTo("LeaveType"), fieldtype.Immutable()).
	Field("from_date", fieldtype.Date, fieldtype.Required()).
	Field("to_date", fieldtype.Date, fieldtype.Required()).
	Field("total_leaves", fieldtype.Float, fieldtype.Required(), fieldtype.MinVal(0)).
	Field("used_leaves", fieldtype.Float, fieldtype.ReadOnly(), fieldtype.Default(float64(0))).
	Field("remaining_leaves", fieldtype.Float, fieldtype.ReadOnly()).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("hr_manager",
		permissions.ActionCreate, permissions.ActionRead,
		permissions.ActionUpdate, permissions.ActionDelete)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

// LeaveRequestDefinition is a request by an employee for leave.
var LeaveRequestDefinition = entitydef.New("LeaveRequest").
	System().
	Label("Leave Request").
	Submittable().
	Field("employee", fieldtype.Link,
		fieldtype.Required(), fieldtype.LinkTo("Employee"), fieldtype.Immutable()).
	Field("leave_type", fieldtype.Link,
		fieldtype.Required(), fieldtype.LinkTo("LeaveType")).
	Field("from_date", fieldtype.Date, fieldtype.Required()).
	Field("to_date", fieldtype.Date, fieldtype.Required()).
	Field("total_days", fieldtype.Float, fieldtype.ReadOnly()).
	Field("reason", fieldtype.SmallText).
	Field("status", fieldtype.Select,
		fieldtype.ReadOnly(),
		fieldtype.Default("Open"),
		fieldtype.Choices("Open", "Approved", "Rejected", "Cancelled")).
	Hook(hooks.BeforeSave, validateLeaveRequestDates).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("hr_manager", permissions.AllActions...)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

func validateLeaveRequestDates(_ context.Context, hctx *hooks.HookContext) error {
	type dataGetter interface{ Get(string) any }
	rec, ok := hctx.Record.(dataGetter)
	if !ok {
		return nil
	}
	from, _ := rec.Get("from_date").(string)
	to, _ := rec.Get("to_date").(string)
	if from != "" && to != "" && to < from {
		return fmt.Errorf("to_date must be on or after from_date")
	}
	return nil
}

// AttendanceRecordDefinition logs daily attendance for an employee.
var AttendanceRecordDefinition = entitydef.New("AttendanceRecord").
	System().
	Label("Attendance Record").
	Field("employee", fieldtype.Link,
		fieldtype.Required(), fieldtype.LinkTo("Employee"), fieldtype.Immutable()).
	Field("attendance_date", fieldtype.Date, fieldtype.Required(), fieldtype.Immutable()).
	Field("status", fieldtype.Select,
		fieldtype.Required(),
		fieldtype.Choices("Present", "Absent", "HalfDay", "OnLeave", "Holiday")).
	Field("in_time", fieldtype.Time).
	Field("out_time", fieldtype.Time).
	Field("working_hours", fieldtype.Float, fieldtype.ReadOnly()).
	Field("remarks", fieldtype.SmallText).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("hr_manager",
		permissions.ActionCreate, permissions.ActionRead, permissions.ActionUpdate)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

// PayrollRunDefinition is the header for a monthly payroll processing run.
var PayrollRunDefinition = entitydef.New("PayrollRun").
	System().
	Label("Payroll Run").
	Submittable().
	Field("run_number", fieldtype.Data,
		fieldtype.ReadOnly(), fieldtype.WithNamingSeries("PR-{YYYY}-{MM}-{SEQ:3}")).
	Field("payroll_month", fieldtype.Select,
		fieldtype.Required(),
		fieldtype.Choices("January", "February", "March", "April", "May", "June",
			"July", "August", "September", "October", "November", "December")).
	Field("payroll_year", fieldtype.Int, fieldtype.Required()).
	Field("payment_date", fieldtype.Date, fieldtype.Required()).
	Field("total_gross", fieldtype.Currency, fieldtype.ReadOnly()).
	Field("total_deductions", fieldtype.Currency, fieldtype.ReadOnly()).
	Field("total_net", fieldtype.Currency, fieldtype.ReadOnly()).
	Field("remarks", fieldtype.SmallText).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("hr_manager",
		permissions.ActionCreate, permissions.ActionRead,
		permissions.ActionUpdate, permissions.ActionSubmit)).
	Allow(permissions.Grant("finance_manager", permissions.ActionRead)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

// PayslipDefinition is the individual salary slip for one employee in a PayrollRun.
var PayslipDefinition = entitydef.New("Payslip").
	System().
	Label("Payslip").
	Field("payroll_run", fieldtype.Link,
		fieldtype.Required(), fieldtype.LinkTo("PayrollRun"), fieldtype.Immutable()).
	Field("employee", fieldtype.Link,
		fieldtype.Required(), fieldtype.LinkTo("Employee"), fieldtype.Immutable()).
	Field("basic_salary", fieldtype.Currency, fieldtype.ReadOnly()).
	Field("gross_pay", fieldtype.Currency, fieldtype.ReadOnly()).
	Field("paye_tax", fieldtype.Currency, fieldtype.ReadOnly()).
	Field("nssf_deduction", fieldtype.Currency, fieldtype.ReadOnly()).
	Field("nhif_deduction", fieldtype.Currency, fieldtype.ReadOnly()).
	Field("other_deductions", fieldtype.Currency, fieldtype.ReadOnly()).
	Field("total_deductions", fieldtype.Currency, fieldtype.ReadOnly()).
	Field("net_pay", fieldtype.Currency, fieldtype.ReadOnly()).
	Field("payment_mode", fieldtype.Select,
		fieldtype.Choices("BankTransfer", "Mpesa", "Cash")).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("hr_manager", permissions.ActionRead)).
	Allow(permissions.Grant("finance_manager", permissions.ActionRead)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

// DisciplinaryWarningDefinition records formal disciplinary actions.
var DisciplinaryWarningDefinition = entitydef.New("DisciplinaryWarning").
	System().
	Label("Disciplinary Warning").
	Field("employee", fieldtype.Link,
		fieldtype.Required(), fieldtype.LinkTo("Employee"), fieldtype.Immutable()).
	Field("warning_date", fieldtype.Date, fieldtype.Required()).
	Field("warning_type", fieldtype.Select,
		fieldtype.Required(),
		fieldtype.Choices("Verbal", "Written", "Final", "Suspension", "Termination")).
	Field("subject", fieldtype.Data, fieldtype.Required(), fieldtype.MaxLen(255)).
	Field("details", fieldtype.LongText).
	Field("issued_by", fieldtype.Link, fieldtype.LinkTo("Employee")).
	Field("acknowledged_on", fieldtype.Date).
	Field("attachment", fieldtype.Attach).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("hr_manager",
		permissions.ActionCreate, permissions.ActionRead, permissions.ActionUpdate)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()
