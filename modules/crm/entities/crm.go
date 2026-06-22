// Package entities contains all CRM module EntityDefinitions.
package entities

import (
	"awo.so/framework/pkg/entitydef"
	"awo.so/framework/pkg/fieldtype"
	"awo.so/framework/pkg/permissions"
)

// CustomerDefinition represents an organisation or individual that buys goods/services.
var CustomerDefinition = entitydef.New("Customer").
	System().
	Label("Customer").
	Field("customer_number", fieldtype.Data,
		fieldtype.ReadOnly(), fieldtype.WithNamingSeries("CUST-{SEQ:6}")).
	Field("name", fieldtype.Data, fieldtype.Required(), fieldtype.Unique(), fieldtype.MaxLen(255)).
	Field("customer_type", fieldtype.Select,
		fieldtype.Required(),
		fieldtype.Default("Company"),
		fieldtype.Choices("Company", "Individual")).
	Field("industry", fieldtype.Select,
		fieldtype.Choices("Agriculture", "Construction", "Education", "Energy",
			"Finance", "Government", "Healthcare", "Hospitality",
			"Manufacturing", "Retail", "Technology", "Transport", "Other")).
	Field("tax_pin", fieldtype.Data, fieldtype.Sensitive(), fieldtype.MaxLen(20)).
	Field("credit_limit", fieldtype.Currency, fieldtype.Default(float64(0))).
	Field("payment_terms_days", fieldtype.Int, fieldtype.Default(30)).
	Field("currency", fieldtype.Data, fieldtype.Default("KES"), fieldtype.MaxLen(10)).
	Field("website", fieldtype.Data, fieldtype.MaxLen(255)).
	Field("notes", fieldtype.SmallText).
	Field("is_disabled", fieldtype.Bool, fieldtype.Default(false)).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("finance_manager",
		permissions.ActionRead, permissions.ActionCreate, permissions.ActionUpdate)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

// ContactDefinition is an individual person associated with a Customer.
var ContactDefinition = entitydef.New("Contact").
	System().
	Label("Contact").
	Field("full_name", fieldtype.Data, fieldtype.Required(), fieldtype.MaxLen(255)).
	Field("customer", fieldtype.Link, fieldtype.LinkTo("Customer")).
	Field("designation", fieldtype.Data, fieldtype.MaxLen(100)).
	Field("email", fieldtype.Data, fieldtype.MaxLen(255)).
	Field("phone", fieldtype.Data, fieldtype.MaxLen(30)).
	Field("is_primary", fieldtype.Bool, fieldtype.Default(false)).
	Field("notes", fieldtype.SmallText).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("finance_manager",
		permissions.ActionRead, permissions.ActionCreate, permissions.ActionUpdate)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

// AddressDefinition stores a physical or postal address for any linked entity.
var AddressDefinition = entitydef.New("Address").
	System().
	Label("Address").
	Field("address_title", fieldtype.Data, fieldtype.Required(), fieldtype.MaxLen(255)).
	Field("address_type", fieldtype.Select,
		fieldtype.Default("Billing"),
		fieldtype.Choices("Billing", "Shipping", "Office", "Other")).
	Field("customer", fieldtype.Link, fieldtype.LinkTo("Customer")).
	Field("line1", fieldtype.Data, fieldtype.Required(), fieldtype.MaxLen(255)).
	Field("line2", fieldtype.Data, fieldtype.MaxLen(255)).
	Field("city", fieldtype.Data, fieldtype.MaxLen(100)).
	Field("county", fieldtype.Data, fieldtype.MaxLen(100)).
	Field("country", fieldtype.Data, fieldtype.Default("Kenya"), fieldtype.MaxLen(100)).
	Field("postal_code", fieldtype.Data, fieldtype.MaxLen(20)).
	Field("is_primary", fieldtype.Bool, fieldtype.Default(false)).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("finance_manager",
		permissions.ActionRead, permissions.ActionCreate, permissions.ActionUpdate)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

// LeadDefinition represents an unqualified sales prospect.
var LeadDefinition = entitydef.New("Lead").
	System().
	Label("Lead").
	Field("lead_number", fieldtype.Data,
		fieldtype.ReadOnly(), fieldtype.WithNamingSeries("LEAD-{YYYY}-{SEQ:5}")).
	Field("lead_name", fieldtype.Data, fieldtype.Required(), fieldtype.MaxLen(255)).
	Field("company_name", fieldtype.Data, fieldtype.MaxLen(255)).
	Field("email", fieldtype.Data, fieldtype.MaxLen(255)).
	Field("phone", fieldtype.Data, fieldtype.MaxLen(30)).
	Field("source", fieldtype.Select,
		fieldtype.Choices("Cold Call", "Email", "Referral", "Social Media",
			"Website", "Exhibition", "Other")).
	Field("status", fieldtype.Select,
		fieldtype.Default("Open"),
		fieldtype.Choices("Open", "Contacted", "Qualified", "Lost", "Converted")).
	Field("notes", fieldtype.SmallText).
	Field("owner", fieldtype.Link, fieldtype.LinkTo("User")).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("finance_manager", permissions.AllActions...)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

// OpportunityDefinition is a qualified sales opportunity converted from a Lead.
var OpportunityDefinition = entitydef.New("Opportunity").
	System().
	Label("Opportunity").
	Submittable().
	Field("opportunity_number", fieldtype.Data,
		fieldtype.ReadOnly(), fieldtype.WithNamingSeries("OPP-{YYYY}-{SEQ:5}")).
	Field("opportunity_name", fieldtype.Data, fieldtype.Required(), fieldtype.MaxLen(255)).
	Field("customer", fieldtype.Link, fieldtype.LinkTo("Customer")).
	Field("lead", fieldtype.Link, fieldtype.LinkTo("Lead")).
	Field("stage", fieldtype.Select,
		fieldtype.Required(),
		fieldtype.Default("Prospecting"),
		fieldtype.Choices("Prospecting", "Qualification", "Proposal",
			"Negotiation", "ClosedWon", "ClosedLost")).
	Field("expected_value", fieldtype.Currency, fieldtype.Default(float64(0))).
	Field("probability_pct", fieldtype.Float,
		fieldtype.Default(float64(0)), fieldtype.MinVal(0), fieldtype.MaxVal(100)).
	Field("expected_close_date", fieldtype.Date).
	Field("owner", fieldtype.Link, fieldtype.LinkTo("User")).
	Field("notes", fieldtype.SmallText).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("finance_manager", permissions.AllActions...)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()
