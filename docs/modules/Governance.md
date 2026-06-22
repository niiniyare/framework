# AwoERP Governance Module
## Comprehensive Documentation Guide

**Module:** `awoerp/governance`
**Version:** 1.0.0
**Platform:** AwoERP — Multi-Tenant SaaS ERP for East African Businesses
**Audience:** Board Members, Compliance Officers, Internal Auditors, System Administrators, Risk Managers, Finance Controllers
**Inspired By:** SAP GRC (Governance, Risk & Compliance), ERPNext Audit & Compliance patterns
**Last Updated:** June 2026

---

## Table of Contents

1. [Introduction & Module Philosophy](#1-introduction--module-philosophy)
2. [Core Concepts & Terminology](#2-core-concepts--terminology)
3. [System Architecture](#3-system-architecture)
4. [Initial Setup & Configuration](#4-initial-setup--configuration)
   - 4.1 [Tenant Governance Settings](#41-tenant-governance-settings)
   - 4.2 [Organisational Hierarchy](#42-organisational-hierarchy)
   - 4.3 [Governance Roles & Responsibilities](#43-governance-roles--responsibilities)
   - 4.4 [Policy Registry Setup](#44-policy-registry-setup)
   - 4.5 [Notification & Escalation Configuration](#45-notification--escalation-configuration)
5. [Risk Management](#5-risk-management)
   - 5.1 [Risk Register](#51-risk-register)
   - 5.2 [Risk Categories & Taxonomy](#52-risk-categories--taxonomy)
   - 5.3 [Risk Assessment & Scoring](#53-risk-assessment--scoring)
   - 5.4 [Risk Appetite & Tolerance](#54-risk-appetite--tolerance)
   - 5.5 [Risk Mitigation Controls](#55-risk-mitigation-controls)
   - 5.6 [Risk Heat Map](#56-risk-heat-map)
   - 5.7 [Key Risk Indicators (KRIs)](#57-key-risk-indicators-kris)
6. [Internal Controls Framework](#6-internal-controls-framework)
   - 6.1 [Control Library](#61-control-library)
   - 6.2 [Control Types & Classifications](#62-control-types--classifications)
   - 6.3 [Control Ownership & Accountability](#63-control-ownership--accountability)
   - 6.4 [Control Testing & Effectiveness](#64-control-testing--effectiveness)
   - 6.5 [Segregation of Duties (SoD)](#65-segregation-of-duties-sod)
   - 6.6 [Compensating Controls](#66-compensating-controls)
7. [Policy & Procedure Management](#7-policy--procedure-management)
   - 7.1 [Policy Lifecycle](#71-policy-lifecycle)
   - 7.2 [Policy Versioning & Change Control](#72-policy-versioning--change-control)
   - 7.3 [Policy Acknowledgement & Training](#73-policy-acknowledgement--training)
   - 7.4 [Regulatory Mapping](#74-regulatory-mapping)
8. [Approval Workflows](#8-approval-workflows)
   - 8.1 [Workflow Engine Architecture](#81-workflow-engine-architecture)
   - 8.2 [Standard Approval Chains](#82-standard-approval-chains)
   - 8.3 [Delegation of Authority (DoA)](#83-delegation-of-authority-doa)
   - 8.4 [Parallel & Sequential Approvals](#84-parallel--sequential-approvals)
   - 8.5 [Escalation & Time-Out Rules](#85-escalation--time-out-rules)
   - 8.6 [Mobile Approvals](#86-mobile-approvals)
9. [Anomaly Detection Engine](#9-anomaly-detection-engine)
   - 9.1 [Detection Algorithms](#91-detection-algorithms)
   - 9.2 [Statistical Outlier Detection](#92-statistical-outlier-detection)
   - 9.3 [Velocity & Acceleration Analysis](#93-velocity--acceleration-analysis)
   - 9.4 [Cross-Tenant Spread Detection](#94-cross-tenant-spread-detection)
   - 9.5 [Amount Sequence Analysis](#95-amount-sequence-analysis)
   - 9.6 [Behavioural Baseline Monitoring](#96-behavioural-baseline-monitoring)
   - 9.7 [Alert Triage & Disposition](#97-alert-triage--disposition)
10. [Audit Management](#10-audit-management)
    - 10.1 [Audit Plan & Universe](#101-audit-plan--universe)
    - 10.2 [Audit Engagements](#102-audit-engagements)
    - 10.3 [Fieldwork & Evidence Management](#103-fieldwork--evidence-management)
    - 10.4 [Audit Findings & Recommendations](#104-audit-findings--recommendations)
    - 10.5 [Management Action Plans (MAPs)](#105-management-action-plans-maps)
    - 10.6 [Follow-Up & Issue Tracking](#106-follow-up--issue-tracking)
11. [Compliance Management](#11-compliance-management)
    - 11.1 [Regulatory Universe](#111-regulatory-universe)
    - 11.2 [Compliance Obligations Register](#112-compliance-obligations-register)
    - 11.3 [KRA Tax Compliance](#113-kra-tax-compliance)
    - 11.4 [EPRA Regulatory Compliance (Petroleum)](#114-epra-regulatory-compliance-petroleum)
    - 11.5 [NEMA Environmental Compliance](#115-nema-environmental-compliance)
    - 11.6 [Compliance Calendar & Deadlines](#116-compliance-calendar--deadlines)
    - 11.7 [Compliance Evidence Vault](#117-compliance-evidence-vault)
12. [Incident Management](#12-incident-management)
    - 12.1 [Incident Types & Classification](#121-incident-types--classification)
    - 12.2 [Incident Reporting & Intake](#122-incident-reporting--intake)
    - 12.3 [Investigation Workflow](#123-investigation-workflow)
    - 12.4 [Root Cause Analysis](#124-root-cause-analysis)
    - 12.5 [Corrective & Preventive Actions (CAPA)](#125-corrective--preventive-actions-capa)
    - 12.6 [Whistleblower Portal](#126-whistleblower-portal)
13. [Reconciliation Health & Financial Controls](#13-reconciliation-health--financial-controls)
    - 13.1 [Reconciliation Subsystem](#131-reconciliation-subsystem)
    - 13.2 [GL Reconciliation Controls](#132-gl-reconciliation-controls)
    - 13.3 [Stock-GL Reconciliation](#133-stock-gl-reconciliation)
    - 13.4 [Cash & Bank Reconciliation](#134-cash--bank-reconciliation)
    - 13.5 [Reconciliation Health Scoring](#135-reconciliation-health-scoring)
14. [TTL-Based Report Cache & Data Freshness](#14-ttl-based-report-cache--data-freshness)
15. [Governance Dashboards & Reports](#15-governance-dashboards--reports)
    - 15.1 [Board-Level Dashboard](#151-board-level-dashboard)
    - 15.2 [Risk Dashboard](#152-risk-dashboard)
    - 15.3 [Compliance Status Report](#153-compliance-status-report)
    - 15.4 [Audit Progress Tracker](#154-audit-progress-tracker)
    - 15.5 [Anomaly Alert Summary](#155-anomaly-alert-summary)
    - 15.6 [Control Effectiveness Report](#156-control-effectiveness-report)
16. [Generalised Audit Log](#16-generalised-audit-log)
    - 16.1 [Audit Log Architecture](#161-audit-log-architecture)
    - 16.2 [Sensitivity Classification](#162-sensitivity-classification)
    - 16.3 [Log Retention & Archival](#163-log-retention--archival)
    - 16.4 [Tamper-Evidence](#164-tamper-evidence)
17. [Feature Flags & Permissions](#17-feature-flags--permissions)
    - 17.1 [Governance Feature Flags](#171-governance-feature-flags)
    - 17.2 [Role-Based Access Control](#172-role-based-access-control)
18. [API Reference & Integration](#18-api-reference--integration)
    - 18.1 [REST API Endpoints](#181-rest-api-endpoints)
    - 18.2 [Webhook Events](#182-webhook-events)
    - 18.3 [External GRC Integration](#183-external-grc-integration)
19. [Regulatory Context — Kenya & East Africa](#19-regulatory-context--kenya--east-africa)
20. [Troubleshooting & FAQs](#20-troubleshooting--faqs)
21. [Glossary](#21-glossary)

---

## 1. Introduction & Module Philosophy

### 1.1 What Is the AwoERP Governance Module?

The **AwoERP Governance Module** (`awoerp/governance`) is the oversight and accountability spine of the AwoERP platform. It brings together Risk Management, Internal Controls, Compliance Tracking, Approval Workflows, Anomaly Detection, Audit Management, and Incident Handling into a single, integrated operational layer — all backed by a tamper-evident audit log and configurable per-tenant feature flags.

Drawing inspiration from **SAP GRC** (Governance, Risk & Compliance Suite) for its structural rigour and enterprise-grade controls framework, and from **ERPNext**'s pragmatic, accessible approach to workflow and compliance tooling for growing businesses, AwoERP's Governance module is calibrated specifically for **East African business realities**: KRA tax obligations, EPRA and NEMA regulatory requirements, Kenyan company law, and the operational environment of multi-site SMEs scaling toward enterprise.

### 1.2 Why Governance Is a First-Class Module

Most ERP platforms bolt on governance as an afterthought — a collection of approval buttons and a report here and there. AwoERP treats governance differently:

**Every transaction in AwoERP produces a governance signal.** When a purchase order is submitted, a stock reconciliation is posted, a cash event occurs at the forecourt, or a user changes a price — the Governance module receives a structured event. It evaluates that event against risk thresholds, control requirements, and compliance obligations, and takes the appropriate automated or human-in-the-loop action.

This is the **SAP GRC principle of continuous controls monitoring (CCM)** applied at the application layer, not bolted on after the fact.

### 1.3 Design Principles

**Separation of Duty by Design.** The governance module is intentionally structurally separate from operational modules. A Finance Manager can approve a payment but cannot modify the audit trail of that approval. A Warehouse Supervisor can submit a reconciliation but cannot approve their own submission.

**Signal-First Architecture.** Every governance concern — a risk alert, an anomaly signal, a compliance deadline — is modelled as a structured **event** written to the `governance.events` table. Downstream consumers (dashboards, email notifications, Temporal workflows) read from this event stream. This decouples detection from response.

**Evidence-Centric Compliance.** Every compliance obligation must have documentary evidence attached. The system tracks evidence completeness and freshness, and refuses to mark an obligation as satisfied without it.

**Human Authority, System Memory.** The system never makes final decisions about risk disposition, control effectiveness, or compliance status — those remain human judgements. But it remembers everything, enforces deadlines, and makes it very difficult to ignore or lose a governance obligation.

---

## 2. Core Concepts & Terminology

### Governance Event

A structured signal emitted by any AwoERP module indicating something happened that governance cares about. Every event has:

```
event_id          UUID
tenant_id         UUID
occurred_at       TIMESTAMPTZ
source_module     VARCHAR         -- 'stock', 'finance', 'sales', 'hr', etc.
event_type        VARCHAR         -- 'anomaly.velocity_spike', 'workflow.approval_timeout', etc.
entity_type       VARCHAR         -- document type that triggered the event
entity_id         UUID
severity          ENUM            -- 'info', 'low', 'medium', 'high', 'critical'
payload           JSONB           -- event-specific structured data
status            ENUM            -- 'open', 'acknowledged', 'investigating', 'resolved', 'dismissed'
assigned_to       UUID (nullable)
resolved_at       TIMESTAMPTZ (nullable)
resolution_notes  TEXT
```

### Risk

A potential event or condition that, if it occurs, would have a negative impact on business objectives. Risks are assessed on **likelihood** and **impact**, and assigned to an **owner** who is accountable for monitoring and mitigation.

### Control

A policy, procedure, or system feature that reduces the likelihood or impact of a risk. Controls are either **preventive** (stop the risk from occurring) or **detective** (identify when the risk has occurred).

### Compliance Obligation

A specific, measurable requirement imposed by a law, regulation, contract, or internal policy that the organisation must fulfil by a defined deadline. Obligations require evidence and have defined responsible parties.

### Anomaly

A data point or pattern that deviates significantly from established norms — detected automatically by the anomaly engine and surfaced as a governance event for human review.

### Audit Finding

A documented observation resulting from an internal audit engagement, classified by severity, assigned to a control owner, and tracked through remediation.

### Reconciliation Health Score

A composite metric (0–100) reflecting how current, complete, and balanced all active reconciliations are for a tenant. Continuously computed and surfaced to Finance and Board dashboards.

---

## 3. System Architecture

### 3.1 Module Interaction Map

```
 ┌─────────────────────────────────────────────────────────────────┐
 │                    AwoERP Operational Modules                    │
 │   Finance │ Stock │ Sales │ Purchasing │ HR │ Forecourt │ Mfg   │
 └───────────────────────────┬─────────────────────────────────────┘
                             │ Domain Events (structured signals)
                             ▼
 ┌─────────────────────────────────────────────────────────────────┐
 │                   governance.events (append-only)                │
 └───────┬─────────────────────────────────────────────────────────┘
         │
         ├──► Anomaly Detection Engine  → anomaly alerts
         ├──► Workflow Engine (Temporal) → approval routing
         ├──► Compliance Checker        → obligation due-date tracking
         ├──► Reconciliation Monitor    → health score computation
         ├──► KRI Evaluator             → threshold breach detection
         └──► Notification Dispatcher   → email / SMS / in-app alerts
```

### 3.2 Temporal Workflow Orchestration

All long-running governance processes — approval chains, audit engagements, incident investigations, CAPA tracking — are orchestrated by **Temporal**. This provides:

- **Durability:** A workflow survives server restarts, deployments, and crashes. An approval pending for 7 days is not lost on a deploy.
- **Timeouts & Escalation:** Built-in timer signals trigger escalations exactly on schedule.
- **Auditability:** Temporal's event history provides a complete record of every step in every workflow.
- **Retries:** Notification delivery failures retry automatically with exponential backoff.

### 3.3 PostgreSQL Schema Overview

```sql
-- Core governance
governance.events
governance.event_assignments

-- Risk management
governance.risks
governance.risk_categories
governance.risk_assessments
governance.risk_controls          -- links risks to controls
governance.kri_definitions
governance.kri_readings

-- Controls
governance.controls
governance.control_tests
governance.control_test_results
governance.sod_rules
governance.sod_violations

-- Policy
governance.policies
governance.policy_versions
governance.policy_acknowledgements
governance.regulatory_obligations

-- Approval workflows
governance.workflow_definitions
governance.workflow_instances
governance.workflow_steps
governance.delegation_of_authority

-- Audit
governance.audit_plans
governance.audit_engagements
governance.audit_fieldwork
governance.audit_findings
governance.management_action_plans
governance.map_evidence

-- Compliance
governance.compliance_obligations
governance.compliance_evidence
governance.compliance_calendar

-- Incidents
governance.incidents
governance.incident_capa
governance.whistleblower_reports

-- Reconciliation
governance.reconciliation_definitions
governance.reconciliation_runs
governance.reconciliation_health_scores

-- Audit log
audit.event_log                   -- generalised, all-module audit trail
audit.log_sensitivity_config      -- configuration-table-driven classification
```

---

## 4. Initial Setup & Configuration

### 4.1 Tenant Governance Settings

Navigate to **Settings → Governance Settings** to configure the global governance posture for your tenant.

#### General Settings

| Setting | Default | Description |
|---|---|---|
| **Governance Module Enabled** | `false` | Master toggle |
| **Governance Framework** | `Custom` | Options: COSO, ISO 31000, King IV, Custom |
| **Default Risk Scoring Method** | `5×5 Matrix` | Likelihood × Impact matrix dimensions |
| **Require Board Approval Above (KES)** | `5,000,000` | Transactions above this trigger Board workflow |
| **Require Finance Approval Above (KES)** | `500,000` | Transactions above this require Finance Controller sign-off |
| **Governance Contact Email** | Required | Receives critical alerts and compliance digest |
| **Audit Committee Chair** | Required | User assigned as audit committee head |

#### Anomaly Detection Settings

| Setting | Default | Description |
|---|---|---|
| **Anomaly Detection Enabled** | `true` | Enable continuous monitoring |
| **Statistical Baseline Window (Days)** | `90` | Window for computing normal behaviour |
| **Sigma Threshold for Alerts** | `3.0` | Standard deviations above mean to flag |
| **Velocity Lookback (Days)** | `30` | Window for velocity/acceleration analysis |
| **Auto-Escalate Critical Anomalies** | `true` | Skip triage for critical-severity anomalies |
| **Cross-Tenant Comparison** | `false` | Use anonymised cross-tenant data for spread detection |

#### Workflow Settings

| Setting | Default | Description |
|---|---|---|
| **Default Approval Timeout (Hours)** | `24` | Hours before an approval step auto-escalates |
| **Allow Self-Approval** | `false` | Whether a user can approve their own submission |
| **Require Comments on Rejection** | `true` | Mandatory rejection reason |
| **Delegation Allowed** | `true` | Allow users to delegate approval authority |
| **Max Delegation Chain** | `2` | Prevent circular or excessively long delegation chains |

---

### 4.2 Organisational Hierarchy

The governance module uses the organisational hierarchy to route approvals, assign control ownership, and scope risk assessments. This is distinct from the Finance cost centre tree — it represents management accountability lines.

```
Navigate: Governance → Organisation → Hierarchy

Example: Anika Global Limited

  Board of Directors
  └── Managing Director (MD)
      ├── Finance Controller
      │   ├── Accounts Officer
      │   └── Tax & Compliance Officer
      ├── Operations Manager
      │   ├── Station Manager — Shell Maanzoni
      │   │   ├── Forecourt Supervisor
      │   │   └── Store Supervisor
      │   └── [Future Site Manager]
      └── IT / Systems Administrator
```

Each node in the hierarchy carries:
- `responsible_user_id` — the individual accountable
- `deputy_user_id` — backup approver during absence
- `escalation_path` — ordered list of escalation targets
- `cost_centre_id` — link to Finance cost centre (optional)

---

### 4.3 Governance Roles & Responsibilities

| Role | Code | Description |
|---|---|---|
| **Board Member** | `gov.board` | View-only access to board dashboards, approve high-value items |
| **Audit Committee Chair** | `gov.audit_chair` | Approve audit plans, receive findings, accept MAPs |
| **Compliance Officer** | `gov.compliance` | Manage obligations, evidence, regulatory universe |
| **Risk Manager** | `gov.risk` | Maintain risk register, assess risks, define KRIs |
| **Internal Auditor** | `gov.auditor` | Execute audit engagements, raise findings |
| **Control Owner** | `gov.control_owner` | Accountable for specific controls; tests and attests |
| **Approver** | `gov.approver` | Approves transactions within defined authority limits |
| **Governance Admin** | `gov.admin` | Configure the governance module (settings, workflows, SoD rules) |
| **Whistleblower Reporter** | `gov.whistleblower` | Anonymous incident submission (no login required) |

---

### 4.4 Policy Registry Setup

Before other governance functions can operate, define your organisation's foundational policies. These become the reference framework for controls, compliance obligations, and audit scope.

```
Navigate: Governance → Policy Registry → New Policy

Policy Code:   POL-FIN-001
Policy Title:  Financial Authority and Approval Policy
Category:      Financial Controls
Owner:         Finance Controller
Effective:     2026-01-01
Review Due:    2027-01-01
Status:        Active

Description:   Establishes the financial authority levels within Anika Global Limited,
               defining monetary thresholds for purchase approvals, payment
               authorisations, contract commitments, and budget adjustments.

Applicable To: All departments
Regulatory Map: Companies Act (Kenya) Cap 486, KRA VAT Regulations 2017

Linked Controls:
  CTRL-001  Purchase Order Approval by Finance for transactions > KES 100,000
  CTRL-002  Dual signatory requirement on bank payments > KES 500,000
  CTRL-003  Board resolution required for capital expenditure > KES 2,000,000
```

---

### 4.5 Notification & Escalation Configuration

```
Navigate: Governance → Settings → Notifications

Notification Channels:
  Email:     Enabled (SMTP via Finance settings)
  SMS:       Enabled (Africa's Talking integration)
  In-App:    Enabled (always on)
  WhatsApp:  Optional (via Africa's Talking WhatsApp API)

Notification Rules:
  Critical Anomaly Detected   → Immediate → MD, Finance Controller, Station Manager
  Approval Pending > 4 hours  → Reminder → Approver
  Approval Pending > 8 hours  → Escalate → Approver's Supervisor
  Compliance Deadline in 7d   → Warning  → Compliance Officer, Obligation Owner
  Compliance Deadline in 1d   → Urgent   → Compliance Officer, MD
  Compliance Overdue          → Critical → All above + Board Chair
  Audit Finding Raised        → Normal   → Finding Owner, Audit Committee Chair
  MAP Overdue                 → Urgent   → MAP Owner, Internal Auditor, MD
  KRI Threshold Breached      → Warning  → Risk Owner, Risk Manager
  Risk Score Changed (>2pts)  → Normal   → Risk Owner, Risk Manager
```

---

## 5. Risk Management

### 5.1 Risk Register

The Risk Register is the central repository of all identified risks. It is a living document — risks are continuously added, assessed, updated, and retired as the business environment changes.

```
Navigate: Governance → Risk Register

List View Columns:
  Risk ID | Risk Title | Category | Inherent Score | Control Status | Residual Score | Owner | Last Reviewed | Status
```

### 5.2 Risk Categories & Taxonomy

AwoERP structures risks in a four-tier taxonomy:

#### Tier 1 — Risk Universe

| Category Code | Category | Description |
|---|---|---|
| `FIN` | Financial Risk | Credit, liquidity, market, valuation risks |
| `OPS` | Operational Risk | Process failures, people, systems, external events |
| `COMP` | Compliance & Regulatory | Laws, regulations, contractual obligations |
| `STRAT` | Strategic Risk | Market changes, competitive threats, business model |
| `TECH` | Technology & Cyber | IT failures, data breaches, system outages |
| `REP` | Reputational Risk | Brand damage, stakeholder trust, media |
| `ENV` | Environmental & ESG | Environmental incidents, climate, sustainability |
| `FRAUD` | Fraud & Integrity | Internal fraud, external theft, corruption |

#### Example Risk Hierarchy (Petroleum Retail)

```
FRAUD
└── FRAUD-01  Internal Theft
    ├── FRAUD-01-01  Pump Meter Manipulation
    ├── FRAUD-01-02  Short-dispensing by Attendants
    ├── FRAUD-01-03  Fictitious Customer Receipts
    └── FRAUD-01-04  Fuel Diversion to Personal Vehicles

OPS
└── OPS-02  Supply Chain Disruption
    ├── OPS-02-01  Fuel Delivery Delay (Vivo Energy supply chain)
    └── OPS-02-02  Transport Strike Action

COMP
└── COMP-01  Regulatory Non-Compliance
    ├── COMP-01-01  EPRA Licence Conditions Breach
    ├── COMP-01-02  KRA VAT Filing Penalty
    └── COMP-01-03  NEMA Environmental Violation
```

### 5.3 Risk Assessment & Scoring

Every risk is scored on two dimensions: **Likelihood** and **Impact**. The default matrix is 5×5.

#### Likelihood Scale

| Score | Level | Description | Frequency Indicator |
|---|---|---|---|
| 1 | Rare | May happen only in exceptional circumstances | < once in 5 years |
| 2 | Unlikely | Could happen at some time | once in 2–5 years |
| 3 | Possible | Might happen at some time | once per year |
| 4 | Likely | Will probably happen | monthly |
| 5 | Almost Certain | Expected to happen in most circumstances | weekly or more |

#### Impact Scale

| Score | Level | Financial Impact (KES) | Operational | Reputational |
|---|---|---|---|---|
| 1 | Negligible | < 10,000 | Minor inconvenience | Internal only |
| 2 | Minor | 10,000 – 100,000 | Short-term disruption | Limited external |
| 3 | Moderate | 100,000 – 1,000,000 | Significant disruption | Local media |
| 4 | Major | 1,000,000 – 10,000,000 | Extended operations impact | National attention |
| 5 | Catastrophic | > 10,000,000 | Business continuity threat | Severe long-term damage |

#### Risk Score Calculation

```
Risk Score = Likelihood × Impact

Score Bands:
  1–4:    Low        (green)
  5–9:    Medium     (amber)
  10–16:  High       (orange)
  17–25:  Critical   (red)
```

#### Full Risk Record

```
Navigate: Governance → Risk Register → New Risk

Risk ID:           RISK-FRAUD-001
Title:             Pump Meter Manipulation by Attendant
Category:          FRAUD-01-01
Description:       An attendant manually interferes with the pump meter or
                   exploits a malfunction to under-record dispensed volume,
                   pocketing the cash difference.
Risk Owner:        Station Manager — Shell Maanzoni
Business Unit:     Forecourt Operations

--- Inherent Risk (before controls) ---
Likelihood:        4 (Likely)
Impact:            4 (Major)
Inherent Score:    16  (HIGH)

--- Residual Risk (after controls) ---
Likelihood:        2 (Unlikely)     ← after controls applied
Impact:            4 (Major)        ← impact unchanged by controls
Residual Score:    8  (MEDIUM)

Linked Controls:
  CTRL-FORE-001  Triple-meter cross-validation per shift
  CTRL-FORE-002  CCTV coverage of all pump islands
  CTRL-FORE-003  Supervisor sign-off on all shift reconciliations
  CTRL-FORE-004  Unannounced dip reading spot checks

KRI:               Daily wetstock variance % (threshold: 0.5% per EPRA)
Last Assessed:     2026-06-01
Next Review:       2026-09-01
Risk Status:       Open / Active
```

### 5.4 Risk Appetite & Tolerance

Risk appetite is the amount of risk the organisation is willing to accept in pursuit of its objectives. Tolerance is the acceptable variance around the appetite.

```
Navigate: Governance → Risk Appetite Statement

Board-Approved Risk Appetite (2026):

Financial Risk:       Appetite LOW    — we do not accept risks that could materially
                                        impair our cash position or KES 500K+

Operational Risk:     Appetite MEDIUM — we accept operational risks with mitigation
                                        controls in place; target residual ≤ HIGH

Compliance Risk:      Appetite ZERO   — we do not accept regulatory violations;
                                        all compliance obligations must be met

Fraud Risk:           Appetite ZERO   — no tolerance for internal fraud or theft

Strategic Risk:       Appetite HIGH   — we accept strategic risk as part of growth;
                                        managed at Board level
```

Any risk whose residual score exceeds the appetite band automatically escalates to the MD and is flagged on the Board Dashboard.

### 5.5 Risk Mitigation Controls

Each risk must have at least one mitigating control. The relationship between risks and controls is many-to-many — a control can mitigate multiple risks, and a risk can be mitigated by multiple controls.

```
Navigate: Governance → Risk Register → [Risk] → Controls tab

Risk:    RISK-FRAUD-001 (Pump Meter Manipulation)

Mitigation Plan:
  Control 1:  CTRL-FORE-001  Triple-meter cross-validation
              Type: Detective | Frequency: Every Shift | Owner: Forecourt Supervisor
              Effectiveness: Effective | Last Tested: 2026-05-15

  Control 2:  CTRL-FORE-002  CCTV coverage
              Type: Preventive & Detective | Frequency: Continuous | Owner: Station Manager
              Effectiveness: Mostly Effective | Last Tested: 2026-04-01

Residual Risk after controls: Score 8 (MEDIUM) — within appetite
Treatment Decision: ACCEPT with monitoring (KRI tracked daily)
```

### 5.6 Risk Heat Map

The Risk Heat Map provides a visual 5×5 grid plotting all active risks by likelihood and impact — both inherent and residual. Accessible from the Risk Dashboard.

```
IMPACT →        1-Negligible  2-Minor   3-Moderate  4-Major  5-Catastrophic
                ─────────────────────────────────────────────────────────────
5-Almost Certain     LOW       MED        HIGH       CRIT      CRIT
4-Likely             LOW       MED        HIGH       HIGH★     CRIT
3-Possible           LOW       LOW        MED        HIGH      HIGH
2-Unlikely           LOW       LOW        MED        MED★      HIGH
1-Rare               LOW       LOW        LOW        LOW       MED

★ = active risks currently plotted at this position
```

### 5.7 Key Risk Indicators (KRIs)

KRIs are leading metrics that signal increasing risk exposure before a risk event occurs.

#### Defining a KRI

```
Navigate: Governance → KRIs → New KRI

KRI ID:           KRI-FORE-001
KRI Name:         Daily Wetstock Variance Percentage
Linked Risk:      RISK-FRAUD-001 (Pump Meter Manipulation)
Data Source:      Stock module — wetstock reconciliation
Formula:          (Theoretical Closing Stock - Actual Closing Stock) / Throughput × 100

Thresholds:
  Green  (In Control):    < 0.3%
  Amber  (Watch):         0.3% – 0.5%
  Red    (Breach):        > 0.5%   (EPRA tolerance limit)

Measurement Frequency:  Daily
Responsible:            Station Manager
Reporting To:           Operations Manager

Alert Actions on Red:
  → Create governance event (severity: HIGH)
  → Notify Station Manager, Operations Manager immediately
  → Escalate to MD if unresolved within 24 hours
  → Create incident record for investigation
```

#### KRI Reading Entry

KRI readings are fed automatically from source modules where possible, or entered manually:

```
Navigate: Governance → KRIs → [KRI-FORE-001] → Readings tab

Date        Reading    Status    Entered By           Source
────────────────────────────────────────────────────────────
2026-06-22  0.23%      GREEN     System (auto)        Stock module SLE
2026-06-21  0.41%      AMBER     System (auto)        Stock module SLE
2026-06-20  0.58%      RED  ⚠   System (auto)        Stock module SLE
2026-06-19  0.19%      GREEN     System (auto)        Stock module SLE
```

---

## 6. Internal Controls Framework

### 6.1 Control Library

The Control Library is the master catalogue of all internal controls operating across the organisation. Inspired by SAP GRC's control repository design, every control is a first-class entity — documented, owned, tested, and rated.

```
Navigate: Governance → Control Library

Control ID:      CTRL-FIN-001
Control Name:    Three-Way Purchase Order Match
Category:        Financial — Procurement
Type:            Preventive
Frequency:       Per Transaction
Automation Level: Semi-Automated (system-prompted, human confirmed)
Owner:           Finance Controller
Risk(s) Mitigated:
  RISK-FIN-002  Fraudulent Supplier Invoices
  RISK-FIN-003  Overpayment to Suppliers

Description:
  Before any supplier invoice is approved for payment, the Finance
  Officer verifies that:
    1. A valid Purchase Order exists
    2. A Goods Receipt Note (GRN) confirms delivery
    3. The Supplier Invoice matches PO and GRN on: item, quantity, price
  Discrepancies are escalated to Finance Controller before payment is released.

Evidence Required:  PO document, GRN document, Supplier Invoice, match confirmation
Testing Procedure:  Sample 25% of invoices monthly; verify 3-way match completion
Last Test Date:     2026-05-01
Test Result:        Effective
Next Test Due:      2026-08-01
```

### 6.2 Control Types & Classifications

#### By Objective

| Type | Description | Example |
|---|---|---|
| **Preventive** | Stops the risk event from occurring | Approval required before PO is submitted |
| **Detective** | Identifies when a risk event has occurred | Daily wetstock reconciliation |
| **Corrective** | Remediates the impact after detection | Stock write-off process after shortage |
| **Directive** | Guides people toward correct behaviour | Signing limits policy, authority matrix |

#### By Nature

| Nature | Description |
|---|---|
| **Manual** | Performed entirely by a person without system assistance |
| **Semi-Automated** | System prompts and tracks; human executes and confirms |
| **Fully Automated** | System performs and enforces with no human action (e.g., SoD block) |

#### By Frequency

`Continuous`, `Per Transaction`, `Hourly`, `Daily`, `Weekly`, `Monthly`, `Quarterly`, `Annual`, `Ad Hoc`

#### COSO Framework Mapping (Optional)

Each control can be mapped to a COSO Internal Control category:

| COSO Component | Description |
|---|---|
| Control Environment | Tone, culture, ethics, HR policies |
| Risk Assessment | Risk identification and analysis processes |
| Control Activities | The actual policies, procedures, and controls |
| Information & Communication | Information systems, reporting, whistleblowing |
| Monitoring Activities | Ongoing and periodic evaluation of controls |

### 6.3 Control Ownership & Accountability

Every control has a **Control Owner** — the individual personally accountable for:

1. Ensuring the control is operating as designed
2. Conducting or facilitating periodic control testing
3. Attesting to control effectiveness at period-end
4. Escalating control failures immediately

Control ownership is enforced through the governance module: owners receive automated reminders for testing due dates, effectiveness attestation deadlines, and any anomalies related to their controls.

### 6.4 Control Testing & Effectiveness

Controls are tested periodically to verify they are operating as designed and are effective.

#### Creating a Control Test

```
Navigate: Governance → Control Library → [Control] → Tests tab → New Test

Control:       CTRL-FIN-001 (Three-Way PO Match)
Test Period:   2026-Q2 (April–June 2026)
Test Date:     2026-07-05
Tester:        Internal Auditor
Test Method:   Sampling — 25% of invoices paid in period

Population:    84 invoices paid in Q2
Sample Size:   21 invoices selected (25%)

Results:
  Exceptions Found:    2
  Exception Details:   2 invoices approved without GRN (GRN was posted 1 day late)
  Exception Rate:      9.5%

Effectiveness Rating:   Mostly Effective     (threshold: < 5% exception = Effective)
Recommendations:        Tighten GRN timing policy — GRN must precede invoice by ≥ 1 day
Follow-Up Required:     Yes — retest in Q3
```

#### Effectiveness Rating Scale

| Rating | Exception Rate | Description |
|---|---|---|
| **Effective** | < 5% | Control operating as designed |
| **Mostly Effective** | 5% – 15% | Minor deviations; improvements needed |
| **Partially Effective** | 15% – 30% | Material deviations; remediation required |
| **Ineffective** | > 30% | Control has failed; immediate action required |
| **Not Tested** | — | Testing overdue |

#### Period-End Control Attestation

At the end of each quarter, control owners receive an automated attestation request:

```
Attestation Request — Q2 2026 — Due: 2026-07-15

Control Owner: [Finance Controller]
Controls for Attestation:
  CTRL-FIN-001  Three-Way PO Match           → Attest Effective / Ineffective / Exception
  CTRL-FIN-002  Dual Bank Payment Signatory  → Attest Effective / Ineffective / Exception
  CTRL-FIN-003  Budget Variance Review       → Attest Effective / Ineffective / Exception

Comments field: (mandatory if any exception attested)
Submit By: 2026-07-15
Approver: Audit Committee Chair
```

### 6.5 Segregation of Duties (SoD)

SoD is one of the most important preventive controls in any ERP system. It ensures that no single person can complete a transaction from start to finish without a second person's involvement. AwoERP enforces SoD rules at the application layer.

#### Defining SoD Rules

```
Navigate: Governance → SoD Rules → New Rule

Rule ID:     SOD-001
Name:        Purchase Order Creation and Approval
Description: A user who creates a Purchase Order cannot also approve it.

Conflicting Permissions:
  Permission A:  purchasing.purchase_order.create
  Permission B:  purchasing.purchase_order.approve

Conflict Type:  Transaction-Level SoD

Enforcement:   HARD BLOCK      (system prevents approval by same user who created)
Severity:      Critical

Risk:          RISK-FIN-002  (Fraudulent Supplier Invoices)
Compensating:  (none — hard block)
```

```
Rule ID:     SOD-002
Name:        Payment Initiation and Bank Authorisation
Description: A user who initiates a bank payment cannot also be the sole bank signatory.

Conflicting Permissions:
  Permission A:  finance.payments.initiate
  Permission B:  finance.bank.sole_authorise

Conflict Type:  User-Role SoD

Enforcement:   SOFT WARN + AUDIT LOG   (dual signatory is physical; system records both)
Severity:      High
```

#### SoD Violation Handling

When a hard-block SoD rule is triggered:

```
System Response:
  1. Block the action immediately
  2. Display clear error: "You cannot approve a transaction you created (SoD Rule SOD-001)"
  3. Create governance event: sod.violation.blocked (severity: HIGH)
  4. Log to audit.event_log with before/after state
  5. Notify Governance Admin if violation frequency exceeds threshold
```

When a soft-warn SoD rule is triggered:

```
System Response:
  1. Display warning and require user acknowledgment
  2. Create governance event: sod.violation.acknowledged (severity: MEDIUM)
  3. Require comment/reason for override
  4. Log to audit.event_log with acknowledged_by and reason
  5. Add to SoD violation report for auditor review
```

#### SoD Violation Report

```
Navigate: Governance → Reports → SoD Violations

Period:  2026-Q2

Rule       Type          Violations  Blocked  Overridden  Override Reason
────────────────────────────────────────────────────────────────────────────────
SOD-001    Hard Block         0          0         0        n/a
SOD-002    Soft Warn          3          0         3        "Dual signatory offline"
SOD-005    Hard Block         1          1         0        n/a (blocked)

Action Required:
  SOD-002 overrides reviewed by Audit Committee → 3 valid exceptions documented
  SOD-005 block: user attempted to approve own stock reconciliation → investigated
```

### 6.6 Compensating Controls

Where a primary SoD control cannot be enforced (e.g., a small operation where roles cannot be fully separated), a **compensating control** documents the alternative mitigation:

```
Navigate: Governance → SoD Rules → [SOD-003] → Compensating Control

SOD Rule:        SOD-003  Stock Reconciliation Creation and Approval
Reason Cannot Enforce: Single-person store team at off-peak hours

Compensating Control:
  Name:         Weekly Manager Review of All Stock Reconciliations
  Description:  Station Manager reviews 100% of all reconciliation entries
                posted in the prior week, signs off evidence log.
  Frequency:    Weekly
  Owner:        Station Manager
  Evidence:     Signed review checklist uploaded to compliance vault
  Risk Accepted: Yes — acknowledged by MD and Audit Committee (2026-01-10)
```

---

## 7. Policy & Procedure Management

### 7.1 Policy Lifecycle

AwoERP manages the full lifecycle of organisational policies:

```
DRAFT → REVIEW → APPROVED → ACTIVE → UNDER REVISION → SUPERSEDED/RETIRED

Transitions:
  DRAFT → REVIEW:          Policy author submits for review
  REVIEW → APPROVED:       Policy owner approves; Audit Committee notified
  APPROVED → ACTIVE:       Effective date reached (automatic)
  ACTIVE → UNDER REVISION: Scheduled review or triggered by regulatory change
  UNDER REVISION → ACTIVE: Revised version approved
  ACTIVE → SUPERSEDED:     Replaced by new policy
  ACTIVE → RETIRED:        Policy no longer needed
```

### 7.2 Policy Versioning & Change Control

Every policy change creates a new version. The full version history is retained indefinitely.

```
Navigate: Governance → Policy Registry → [POL-FIN-001] → Versions

Version  Effective Date  Status     Changed By          Change Summary
───────────────────────────────────────────────────────────────────────────────
v1.0     2024-01-01      Superseded Finance Controller  Initial policy
v1.1     2025-03-01      Superseded Finance Controller  Raised Board threshold to 5M
v2.0     2026-01-01      Active     MD + Board          Major revision; added e-approval
```

Changes to policies above a defined materiality level require MD or Board approval and a documented change rationale. Minor changes (grammar, formatting) can be approved by the policy owner alone.

### 7.3 Policy Acknowledgement & Training

```
Navigate: Governance → Policy Registry → [Policy] → Acknowledgements tab

Rollout:
  Target Users:     All Staff (14 users)
  Deadline:         2026-02-28
  Acknowledgement:  User must read and digitally sign
  Training Required: Yes — link to training video

Status:
  Acknowledged:     12 / 14 users  (85.7%)
  Overdue:          2 users        ⚠ reminder sent 2026-02-25
  Exceptions:       0 formal exemptions

Compliance Obligation:
  This policy acknowledgement is linked to COMP-OBL-011 (Staff Policy
  Certification) which is due annually to meet ISO 9001 internal audit requirements.
```

### 7.4 Regulatory Mapping

Each policy and control can be mapped to specific clauses in relevant regulations:

```
Navigate: Governance → Regulatory Universe → [Regulation] → Mapping tab

Regulation:  Kenya Companies Act, Cap 486, 2015
Clause:      Section 153 — Annual Financial Statements
             Section 155 — Director Responsibility for Accounts

Mapped Policies:
  POL-FIN-002  Financial Reporting Policy
  POL-FIN-003  Director Approval of Accounts Policy

Mapped Controls:
  CTRL-FIN-008  Annual audit engagement completion
  CTRL-FIN-009  Board sign-off on financial statements

Compliance Obligation:
  COMP-OBL-001  File Annual Return with Registrar of Companies (due 30 April)
```

---

## 8. Approval Workflows

### 8.1 Workflow Engine Architecture

AwoERP's approval workflows are powered by **Temporal** workflow orchestration. Each workflow definition is a Go struct that defines:

- The sequence of approval steps
- Conditions for routing (amount thresholds, document types, departments)
- Timeout and escalation behaviour
- Parallel vs sequential step execution
- Outcome actions (post document, reject, return for amendment)

```go
// Simplified workflow definition (Go pseudocode)
type ApprovalWorkflow struct {
    WorkflowID   string
    DocumentType string
    Steps        []ApprovalStep
    Conditions   []RoutingCondition
    OnApprove    []Action
    OnReject     []Action
}

type ApprovalStep struct {
    StepName      string
    ApproverRole  string
    TimeoutHours  int
    EscalateTo    string
    AllowDelegate bool
    RequireReason bool    // on rejection
    Parallel      bool    // run concurrently with next step
}
```

### 8.2 Standard Approval Chains

#### Purchase Order Approval

```
PO Amount           Approver Chain
──────────────────────────────────────────────────────────────────
< KES 50,000        Store Supervisor  (1 step)
KES 50K – 500K      Store Supervisor → Finance Controller (2 steps)
KES 500K – 2M       Store Supervisor → Finance Controller → MD (3 steps)
KES 2M – 5M         Store Supervisor → Finance Controller → MD → Board Chair (4 steps)
> KES 5M            Full Board Resolution required (offline; minutes uploaded as evidence)
```

#### Payment Approval (Supplier Invoice)

```
Invoice Amount      Approver Chain
──────────────────────────────────────────────────────────────────
< KES 100,000       Accounts Officer (verify 3-way match)
KES 100K – 500K     Accounts Officer → Finance Controller
KES 500K – 2M       Accounts Officer → Finance Controller → MD
> KES 2M            Finance Controller + MD (parallel, dual signature)
```

#### Stock Reconciliation Approval

```
Variance Level      Approver Chain
──────────────────────────────────────────────────────────────────
< 0.5% of value     Store Supervisor (auto-approve if zero variance)
0.5% – 2% of value  Store Supervisor → Station Manager
2% – 5% of value    Store Supervisor → Station Manager → Finance Controller
> 5% of value       All above + MD; create governance incident
```

### 8.3 Delegation of Authority (DoA)

The DoA matrix formally documents who can approve what, at what financial limit. It is stored in the system and enforces the approval routing.

```
Navigate: Governance → Delegation of Authority → DoA Matrix

Role/Position       Document Type         Limit (KES)   Notes
──────────────────────────────────────────────────────────────────────────────
Attendant           Requisition            10,000        For consumables only
Store Supervisor    Purchase Order         50,000        Stocked items only
Station Manager     Purchase Order        500,000        Any category
Finance Controller  Purchase Order      2,000,000        All categories
Finance Controller  Payment Release     2,000,000        Dual signatory > 500K
Managing Director   Purchase Order      5,000,000        All categories
Managing Director   Payment Release     5,000,000        Dual signatory required
Board of Directors  Capital Expenditure Unlimited        Board resolution required
```

#### Temporary Delegation

```
Navigate: Governance → Delegation of Authority → Temporary Delegation → New

Delegating User:    Finance Controller
Delegate To:        Accounts Officer
Scope:              Purchase Order approval up to KES 200,000
Reason:             Finance Controller on annual leave
Valid From:         2026-07-01 08:00
Valid To:           2026-07-14 18:00
Approved By:        MD
Notifications:      All pending approvals reassigned to delegate
```

The delegation is automatically revoked at the Valid To timestamp, with all pending items reverting to the original approver.

### 8.4 Parallel & Sequential Approvals

#### Sequential (Default)

Steps complete in order. The second approver only receives the request after the first has approved.

```
Step 1: Finance Controller ────► Step 2: MD ────► Document Posted
```

#### Parallel

All designated approvers receive the request simultaneously. The document proceeds when all approve (AND logic) or when any approve (OR logic).

```
Step 2a: Finance Controller ──┐
                               ├──► (Both required) ────► Document Posted
Step 2b: MD ──────────────────┘
```

Parallel approvals are commonly used for the dual-signatory bank payment control.

### 8.5 Escalation & Time-Out Rules

```
Temporal Workflow Timer Logic:

Step Created
    │
    ├──► t + 4 hours: Send reminder to approver
    ├──► t + 8 hours: Send urgent reminder + notify approver's supervisor
    ├──► t + TimeoutHours: Auto-escalate to EscalateTo user
    │        │
    │        └──► EscalateTo t + TimeoutHours: Escalate to next level
    │                 │
    │                 └──► Final escalation: Notify MD
    │
    └──► If no action by MaxEscalationDays: Auto-reject + notify all parties
```

The Temporal workflow persists all timer state — the escalation will fire exactly on schedule even if the server restarts.

### 8.6 Mobile Approvals

For managers who are frequently on-site or travelling, AwoERP supports approval actions from the mobile app:

- Push notifications on iOS and Android when an approval is pending
- Swipe-to-approve or tap-to-reject from the notification itself
- View document summary before deciding (key fields, amount, submitted by, previous approvals)
- Mandatory comment field on rejection, even from mobile
- Biometric authentication required for approvals above a configured threshold

---

## 9. Anomaly Detection Engine

### 9.1 Detection Algorithms

The AwoERP anomaly detection engine runs continuously, evaluating every significant transaction and data point against a suite of detection algorithms. It is implemented as a pipeline of functional options — each detector is independently configured and can be enabled or disabled per tenant via feature flags.

```go
// Functional options pattern (Go)
type AnomalyDetector func(ctx context.Context, event DomainEvent) ([]AnomalySignal, error)

var DefaultDetectors = []AnomalyDetector{
    StatisticalOutlierDetector(SigmaThreshold(3.0), WindowDays(90)),
    VelocityAccelerationDetector(LookbackDays(30), SpikeThreshold(2.5)),
    AmountSequenceDetector(BenfordLawEnabled(true)),
    BehaviouralBaselineDetector(BaselineDays(60)),
    CrossTenantSpreadDetector(MinTenants(10), Percentile(95)),
    ReconciliationHealthDetector(HealthThreshold(70)),
}
```

### 9.2 Statistical Outlier Detection

Detects individual data points that fall outside the expected range based on a rolling statistical baseline.

**Algorithm:**

```
For each measurement type M (e.g., daily fuel sales volume):

1. Compute rolling mean μ and standard deviation σ over last 90 days
2. For each new observation x:
   z-score = (x - μ) / σ
   IF abs(z-score) > 3.0:
     → Raise anomaly signal
     → Severity scaled by z-score:
         3.0–4.0σ  → MEDIUM
         4.0–5.0σ  → HIGH
         > 5.0σ    → CRITICAL
```

**Example — Fuel Sales Volume Anomaly:**

```
Baseline (90-day):  Mean = 8,200 L/day, σ = 850 L/day
Today's Sales:      2,400 L

z-score = (2,400 - 8,200) / 850 = -6.82σ

Signal:  ANOMALY-STAT-001
Type:    statistical_outlier.below_lower_bound
Severity: CRITICAL
Context: Daily diesel sales at Shell Maanzoni: 2,400 L vs expected 8,200 L (–68%)
Action:  Immediate alert to Station Manager, Operations Manager, MD
```

### 9.3 Velocity & Acceleration Analysis

Detects when the rate of change (velocity) or the change in rate (acceleration) of a metric is abnormal — catching gradual drift before it becomes a large deviation.

**Velocity:** Rate of change over a rolling period.
**Acceleration:** Change in velocity — is the rate itself speeding up or slowing down?

```
Daily Wetstock Loss — 30-Day Velocity Analysis:

Day 1–10:  0.15% average loss/day   (baseline)
Day 11–20: 0.22% average loss/day   (velocity +0.07%/period)
Day 21–30: 0.38% average loss/day   (velocity +0.16%/period — acceleration detected)

Signal:  ANOMALY-VEL-001
Type:    velocity.accelerating_loss
Severity: HIGH
Message: Wetstock loss rate has been accelerating over 30 days (+153% velocity).
         Current rate approaching EPRA tolerance threshold.
```

### 9.4 Cross-Tenant Spread Detection

With the user's consent and full anonymisation, AwoERP compares a tenant's metrics to the distribution across similar tenants (same industry, similar size). This surfaces anomalies that are hidden when looking at a single tenant's history alone — for example, a price suddenly 40% above market rate.

```
Signal:  ANOMALY-SPREAD-001
Type:    cross_tenant.price_outlier
Severity: MEDIUM
Context: Purchase price for DIESEL-AGO on GRN-2026-00089 (KES 162.50/L) is 
         in the 98th percentile of peer group rates for this product.
         Peer group median: KES 155.40/L
Note:   Cross-tenant data is fully anonymised; no competitor data is disclosed.
```

This feature is **opt-in** and controlled by the `governance.cross_tenant_spread` feature flag. It is disabled by default.

### 9.5 Amount Sequence Analysis

Uses **Benford's Law** to detect manipulation in financial figures. In naturally occurring datasets, the leading digit of amounts follows a logarithmic distribution (1 appears ~30% of the time, 9 appears ~5%). When amounts are fabricated or manipulated, they deviate from this distribution.

```
Benford's Law Analysis — Supplier Invoices — 2026-Q1

Leading Digit  Expected %  Actual %   Deviation
─────────────────────────────────────────────────
1              30.1%       22.4%       -7.7%
2              17.6%       15.2%       -2.4%
3              12.5%       11.8%       -0.7%
4              9.7%        8.1%        -1.6%
5              7.9%        6.4%        -1.5%
6              6.7%        5.9%        -0.8%
7              5.8%        8.2%        +2.4%
8              5.1%        11.3%       +6.2%  ⚠
9              4.6%        10.7%       +6.1%  ⚠

Anomaly Detected: Significant excess of invoices starting with 8 and 9.
Signal:           ANOMALY-AMOUNT-001
Type:             amount_sequence.benford_deviation
Severity:         HIGH
Context:          Unusual clustering of amounts 8,000–9,999 range — possible
                  fabrication below approval thresholds. Refer to audit.
```

### 9.6 Behavioural Baseline Monitoring

Tracks user and entity behaviour patterns and alerts on deviations.

**User Behaviour:**

```
User: [Attendant ID]
Normal Pattern (60-day baseline):
  Login Times:   06:00–22:00 KE
  Actions/Day:   45–80 transactions
  Peak:          Weekday morning and evening

Anomaly Detected:
  Login at 02:43 EAT (off hours)
  12 stock entries created in 15 minutes
  All entries modified an item not normally touched by this user

Signal:  ANOMALY-BEHAV-001
Type:    behavioural.off_hours_unusual_activity
Severity: HIGH
Action:  Alert IT Admin and Station Manager; session flagged for audit
```

**Entity Behaviour (Supplier):**

```
Supplier: [Supplier ID]
Normal Pattern: 2–4 invoices/month, range KES 200K–800K

Anomaly Detected:
  11 invoices submitted in 3 days
  All invoices for KES 98,000 (just below KES 100,000 Finance approval threshold)

Signal:  ANOMALY-BEHAV-002
Type:    behavioural.threshold_splitting
Severity: CRITICAL
Action:  Hold all pending invoices from supplier; create incident; notify Finance Controller
```

### 9.7 Alert Triage & Disposition

Every anomaly signal creates a Governance Event in `governance.events` and enters a triage workflow.

#### Alert Status Flow

```
OPEN → ACKNOWLEDGED → INVESTIGATING → RESOLVED
                   └─────────────────► DISMISSED (with justification)
```

#### Triage Interface

```
Navigate: Governance → Anomaly Alerts → Triage Queue

Alert ID         Type                      Severity  Age    Assigned To  Status
─────────────────────────────────────────────────────────────────────────────────
ALRT-2026-0412  statistical_outlier        CRITICAL  2h     [unassigned] OPEN  ⚠
ALRT-2026-0411  amount_sequence.benford    HIGH      6h     Auditor      INVESTIGATING
ALRT-2026-0409  velocity.accelerating_loss HIGH      18h    Stn Manager  ACKNOWLEDGED
ALRT-2026-0408  behavioural.off_hours      HIGH      1d     IT Admin     RESOLVED
```

#### Disposition Actions

| Action | When to Use | Outcome |
|---|---|---|
| **Acknowledge** | Alert received, initial review done | Stops escalation timers; assigns owner |
| **Investigate** | Alert requires deeper investigation | Creates investigation task; pauses auto-escalation |
| **Resolve** | Root cause identified; corrective action taken | Closes alert; documents resolution; links CAPA |
| **Dismiss** | Alert is a known false positive | Closes alert; records dismissal reason; feeds back into model tuning |
| **Escalate** | Requires senior attention | Pushes to next authority level immediately |

---

## 10. Audit Management

### 10.1 Audit Plan & Universe

The Audit Universe is the complete inventory of auditable units — every business process, department, system, or location that could theoretically be audited. The Annual Audit Plan selects which units to audit in a given year, based on risk assessment and strategic priorities.

```
Navigate: Governance → Audit → Audit Universe

Unit ID    Unit Name                        Risk Level  Last Audit   Next Scheduled
─────────────────────────────────────────────────────────────────────────────────
AU-001     Forecourt Cash Handling          HIGH        2025-Q4      2026-Q2 ✓
AU-002     Fuel Procurement Process         HIGH        2025-Q3      2026-Q3
AU-003     Payroll & HR                     MEDIUM      2025-Q2      2026-Q4
AU-004     IT Access Controls               HIGH        2025-Q4      2026-Q2 ✓
AU-005     Stock Reconciliation Process     HIGH        2025-Q1      2026-Q1 ✓ (done)
AU-006     Supplier Management              MEDIUM      2025-Q3      2027-Q1
AU-007     Environmental Compliance         MEDIUM      2025-Q4      2026-Q3
```

#### Annual Audit Plan

```
Navigate: Governance → Audit → Audit Plan → New Plan

Plan Year:  2026
Approved By: Audit Committee Chair (2026-01-15)

Planned Engagements:
  Q1:  AU-005  Stock Reconciliation          Status: COMPLETE
  Q2:  AU-001  Forecourt Cash Handling        Status: IN PROGRESS
       AU-004  IT Access Controls             Status: PLANNED
  Q3:  AU-002  Fuel Procurement               Status: PLANNED
       AU-007  Environmental Compliance       Status: PLANNED
  Q4:  AU-003  Payroll & HR                   Status: PLANNED

Unplanned/Advisory:
       Triggered by anomaly or Board request  (budget reserved: 10% of audit hours)
```

### 10.2 Audit Engagements

An Audit Engagement is a formal, scoped audit project with defined objectives, timeline, team, and evidence trail.

```
Navigate: Governance → Audit → Engagements → New Engagement

Engagement ID:     AUD-2026-002
Title:             Forecourt Cash Handling & Shift Reconciliation Audit
Audit Unit:        AU-001
Engagement Type:   Assurance
Lead Auditor:      Internal Auditor
Reviewer:          Audit Committee Chair

Objectives:
  1. Verify that shift cash collections reconcile to meter sales daily
  2. Confirm that all cash events are documented and authorised
  3. Assess adequacy of petty cash controls and float management
  4. Evaluate physical security of cash at station

Scope:
  Period:       2026-01-01 to 2026-03-31
  Location:     Shell Maanzoni Service Station
  Systems:      AwoERP (Forecourt module, Finance module)

Timeline:
  Planning:     2026-06-01 – 2026-06-07
  Fieldwork:    2026-06-08 – 2026-06-22
  Reporting:    2026-06-23 – 2026-06-30
  Response:     2026-07-01 – 2026-07-15
  Final Report: 2026-07-16

Status:         FIELDWORK IN PROGRESS
```

### 10.3 Fieldwork & Evidence Management

During fieldwork, auditors document their testing, attach evidence, and record observations.

```
Navigate: Governance → Audit → [Engagement] → Fieldwork tab → New Workpaper

Workpaper Ref:  WP-AUD-2026-002-03
Title:          Cash Count Reconciliation Testing
Objective:      OBJ-2 — Verify daily cash reconciliations

Procedure:
  Selected 20 shifts from the audit period.
  For each shift, obtained:
    - Shift Reconciliation Report from AwoERP
    - Cash count sheet signed by Supervisor
    - POS/cash register Z-report
  Agreed cash per AwoERP to physical count; traced to bank deposit slip.

Population:     62 shifts in period
Sample:         20 shifts (32.3%)
Sampling Method: Risk-based (weekends, month-end)

Results:
  Exceptions:   3 shifts where AwoERP cash balance did not agree to signed count
  Exception 1:  2026-01-14 — KES 1,200 shortfall; no documentation found
  Exception 2:  2026-02-28 — KES 500 excess; documented as rejected note
  Exception 3:  2026-03-07 — KES 4,800 shortfall; shift supervisor absent

Evidence Attached:
  [Shift_Reconciliation_Sample.pdf]
  [Cash_Count_Sheet_Jan14.pdf]
  [Bank_Deposit_Slips_Q1.pdf]
```

### 10.4 Audit Findings & Recommendations

```
Navigate: Governance → Audit → [Engagement] → Findings tab → New Finding

Finding ID:     FIND-2026-002-01
Title:          Undocumented Cash Shortfalls in Shift Reconciliations
Severity:       HIGH
Control:        CTRL-CASH-003 (Daily Shift Cash Reconciliation)
Control Rating: Partially Effective

Condition:
  3 of 20 sampled shifts (15%) contained cash shortfalls ranging
  from KES 500 to KES 4,800 with insufficient documentation of cause.

Criteria:
  Station Management Policy Section 4.2 requires all cash variances
  exceeding KES 200 to be documented with root cause and supervisor sign-off.

Cause:
  Supervisors are not consistently enforcing the variance documentation
  requirement. The policy appears not to be clearly understood at operational level.

Risk:
  RISK-FRAUD-003 (Misappropriation of cash collections). Undocumented
  shortfalls may conceal systematic theft.

Recommendation:
  1. Conduct mandatory refresher training on cash variance documentation for all shifts
  2. Implement AwoERP alert when a shift is submitted with >KES 200 variance and
     no documentation reason code entered
  3. Station Manager to perform weekly review of all variance documentation

Finding Owner:  Station Manager — Shell Maanzoni
Agreed By:      Station Manager (2026-06-25)
```

### 10.5 Management Action Plans (MAPs)

For each finding, management commits to a specific remediation action with a deadline.

```
Navigate: Governance → Audit → [Finding] → Management Action Plan

Finding:   FIND-2026-002-01
MAP Owner: Station Manager

Actions:

  Action 1:
    Description:  Conduct mandatory cash-handling refresher training
                  for all forecourt and cashier staff
    Responsible:  Station Manager
    Due Date:     2026-07-31
    Evidence:     Training attendance register; signed acknowledgements
    Status:       IN PROGRESS

  Action 2:
    Description:  Configure AwoERP alert for shift submissions with
                  cash variance > KES 200 without reason code
    Responsible:  System Administrator
    Due Date:     2026-07-15
    Evidence:     Screenshot of configured alert; test shift submission
    Status:       COMPLETE ✓ (completed 2026-07-10)

  Action 3:
    Description:  Implement weekly Manager review log for all variance docs
    Responsible:  Station Manager
    Due Date:     2026-08-01
    Evidence:     Weekly review template; first 4 weeks of signed logs
    Status:       NOT STARTED
```

### 10.6 Follow-Up & Issue Tracking

```
Navigate: Governance → Audit → Follow-Up Dashboard

Open Findings Requiring MAP Update:

Finding             Due Date     MAP Owner       Status         Days Overdue
────────────────────────────────────────────────────────────────────────────
FIND-2026-002-01   2026-07-31   Stn Manager     IN PROGRESS    0
FIND-2026-001-03   2026-06-30   Finance Ctrl    OVERDUE        ⚠ 25 days
FIND-2025-004-02   2026-05-15   IT Admin        OVERDUE        ⚠ 58 days

Escalation:
  FIND-2026-001-03:  Escalated to MD (2026-07-05); response pending
  FIND-2025-004-02:  Escalated to Board Chair; discussed at June Board meeting
```

---

## 11. Compliance Management

### 11.1 Regulatory Universe

The Regulatory Universe catalogues every law, regulation, and standard that applies to the organisation.

```
Navigate: Governance → Compliance → Regulatory Universe

Reg ID    Regulation                                          Jurisdiction  Status
──────────────────────────────────────────────────────────────────────────────────
REG-001   Companies Act, Cap 486 (2015)                      Kenya         Active
REG-002   Income Tax Act, Cap 470                            Kenya         Active
REG-003   Value Added Tax Act, 2013 (+ 2020 amendments)      Kenya         Active
REG-004   Employment Act, Cap 226                            Kenya         Active
REG-005   Energy Act, 2019                                   Kenya         Active
REG-006   EPRA Petroleum Retail Licence Conditions           Kenya         Active
REG-007   NEMA Environmental Management Act, 1999            Kenya         Active
REG-008   Data Protection Act, 2019                          Kenya         Active
REG-009   OHSA — Occupational Safety & Health Act, 2007      Kenya         Active
REG-010   KRA Tax Procedures Act, 2015                       Kenya         Active
```

### 11.2 Compliance Obligations Register

Each regulation spawns specific, measurable obligations. The obligations register tracks every one.

```
Navigate: Governance → Compliance → Obligations Register

Obligation ID:   COMP-OBL-003
Title:           Monthly VAT Return Filing
Regulation:      REG-003 (VAT Act 2013)
Clause:          Section 31 — VAT Returns
Description:     File VAT return and pay any VAT liability by the 20th
                 of the month following the tax period.
Frequency:       Monthly
Deadline Rule:   20th of following month (if weekend: next business day)
Owner:           Tax & Compliance Officer
Approver:        Finance Controller
Evidence Required:
  - iTax VAT return submission receipt
  - Payment confirmation slip
  - Signed reconciliation of output and input VAT
Penalty for Non-Compliance: 5% of tax due per month late + KES 10,000 minimum
Severity if Missed: CRITICAL
```

### 11.3 KRA Tax Compliance

```
Navigate: Governance → Compliance → KRA Calendar

Annual KRA Obligations (Anika Global Limited):

Month  Obligation                              Due Date  Status
──────────────────────────────────────────────────────────────────
Jan    VAT Return (December)                  20 Jan    COMPLETE
Jan    PAYE (December payroll)                9 Jan     COMPLETE
Jan    NSSF (December)                        15 Jan    COMPLETE
Jan    NHIF (December)                        15 Jan    COMPLETE
Feb    VAT Return (January)                   20 Feb    COMPLETE
Feb    PAYE (January)                         9 Feb     COMPLETE
Mar    VAT Return (February)                  20 Mar    COMPLETE
Apr    Corporate Income Tax Instalment        20 Apr    COMPLETE
Apr    Annual Company Return (Registrar)      30 Apr    COMPLETE
...
Jun    VAT Return (May)                       20 Jun    COMPLETE
Jun    PAYE (May)                             9 Jun     COMPLETE
Jul    VAT Return (June)                      20 Jul    UPCOMING ← today
Dec    Annual Income Tax Return               31 Dec    PLANNED
```

### 11.4 EPRA Regulatory Compliance (Petroleum)

```
Navigate: Governance → Compliance → EPRA Dashboard

EPRA Licence Conditions — Shell Maanzoni Service Station

Obligation                              Frequency  Status       Next Due
────────────────────────────────────────────────────────────────────────────────
Retail Price Compliance (pump prices    Daily      COMPLIANT    Ongoing
  within EPRA maximum retail price)
Daily Wetstock Returns (submit to EPRA  Daily      COMPLIANT    Tomorrow
  online portal)
Monthly Stock & Sales Return            Monthly    COMPLIANT    5 Jul 2026
Annual Licence Renewal                  Annual     COMPLIANT    31 Dec 2026
Pump Meter Calibration (KEBS certified) Annual     COMPLIANT    Nov 2026
Tank Integrity Test                     5-yearly   DUE          Q3 2027
Vapour Recovery Unit Maintenance        Annual     OVERDUE ⚠    15 Jun 2026
Environmental Management Plan Review    Annual     COMPLIANT    Mar 2027
```

### 11.5 NEMA Environmental Compliance

```
Navigate: Governance → Compliance → NEMA Tracker

Environmental Obligations:

Obligation                         Frequency  Owner           Status
──────────────────────────────────────────────────────────────────────
Environmental Impact Assessment     Once       MD / NEMA       COMPLETE (2019)
Annual Environmental Audit          Annual     Ext. Consultant OVERDUE ⚠ (due Apr 2026)
Spill Incident Reporting            Ad Hoc     Stn Manager     N/A
Hazardous Waste Disposal Records    Monthly    Stn Manager     COMPLIANT
Water effluent discharge licence    Annual     Stn Manager     COMPLIANT
Soil contamination monitoring       Bi-annual  Ext. Consultant COMPLIANT
```

### 11.6 Compliance Calendar & Deadlines

```
Navigate: Governance → Compliance → Calendar View

                         JULY 2026
Mon   Tue   Wed   Thu   Fri   Sat   Sun
                    1     2     3     4     5
  6     7     8     9    10    11    12
                         ↑
                    PAYE deadline (6 Jul)

 13    14    15    16    17    18    19
                    ↑
              VAT deadline (20th approaching)

 20    21    22    23    24    25    26
  ↑
VAT & EPRA Monthly Return (20 Jul)

 27    28    29    30    31
```

Colour coding: Green = complete, Amber = due within 7 days, Red = overdue.

### 11.7 Compliance Evidence Vault

The evidence vault stores all documentary proof of compliance obligation fulfilment.

```
Navigate: Governance → Compliance → Evidence Vault

Evidence ID:       EV-2026-06-003
Obligation:        COMP-OBL-003 (Monthly VAT Return — June 2026)
Uploaded By:       Tax & Compliance Officer
Upload Date:       2026-07-19
Evidence Type:     Tax Return Receipt + Payment Confirmation

Files:
  [iTax_VAT_Return_Jun2026_Acknowledgement.pdf]
  [MPESA_Tax_Payment_Confirmation_Jul2026.pdf]
  [VAT_Reconciliation_Jun2026_Signed.xlsx]

Verified By:       Finance Controller (2026-07-19)
Retention Until:   2033-07-19     (KRA 7-year retention requirement)
Storage:           Encrypted, access-logged, immutable after verification
```

---

## 12. Incident Management

### 12.1 Incident Types & Classification

| Category | Type Code | Examples | Severity Range |
|---|---|---|---|
| **Financial** | `INC-FIN` | Fraud, misappropriation, billing errors, overpayment | MEDIUM–CRITICAL |
| **Operational** | `INC-OPS` | Equipment failure, supply disruption, safety incident | LOW–HIGH |
| **Compliance** | `INC-COMP` | Regulatory breach, licence violation, reporting failure | HIGH–CRITICAL |
| **Security** | `INC-SEC` | Theft, CCTV failure, unauthorized access, data breach | HIGH–CRITICAL |
| **Environmental** | `INC-ENV` | Fuel spill, contamination, vapour release | MEDIUM–CRITICAL |
| **HR / Conduct** | `INC-HR` | Misconduct, harassment, policy violation | MEDIUM–HIGH |
| **IT / Systems** | `INC-IT` | System outage, data loss, cyber incident | MEDIUM–CRITICAL |

### 12.2 Incident Reporting & Intake

#### Standard Incident Report

```
Navigate: Governance → Incidents → New Incident

Incident Title:    Unexplained Cash Shortage — Morning Shift June 22
Category:          INC-FIN
Sub-Type:          Cash misappropriation
Severity:          HIGH
Reported By:       Store Supervisor
Date/Time Occurred: 2026-06-22 14:00 EAT
Date/Time Reported: 2026-06-22 14:45 EAT
Location:          Shell Maanzoni — Cashier Station

Description:
  At shift close on June 22, the cashier balance was KES 8,400 short of the
  expected amount per AwoERP shift reconciliation. No explanation could be provided
  by the attendant. The supervisor has secured the shift report and CCTV footage.

Immediate Action Taken:
  Attendant placed on administrative leave pending investigation.
  CCTV footage preserved for last 24 hours.
  Finance Controller notified at 14:50.

Parties Involved:
  Attendant: [Name] — shift 06:00–14:00
  Supervisor: [Name] — authorised shift close

Evidence Secured:
  [Shift_Reconciliation_Jun22_Morning.pdf]
  [CCTV_Cashier_Jun22_0600_1400.mp4] (flagged — requires IT custody)
```

### 12.3 Investigation Workflow

```
REPORTED → ASSIGNED → INVESTIGATING → FINDINGS DOCUMENTED → CAPA RAISED → CLOSED

Navigate: Governance → Incidents → [INC-2026-034] → Workflow

Step 1: Initial Assessment (within 2 hours of report)
  Assigned To:      Finance Controller + IT Admin
  Assessment:       Credible — evidence preserved; formal investigation required
  Decision:         Proceed to full investigation

Step 2: Investigation (within 5 business days)
  Lead Investigator: Finance Controller
  Support:           Station Manager, IT Admin (CCTV review)

Step 3: Findings Report (within 2 business days of investigation complete)

Step 4: MD Review & Decision (within 2 business days)

Step 5: CAPA Implementation & Monitoring
```

### 12.4 Root Cause Analysis

```
Navigate: Governance → Incidents → [INC-2026-034] → Root Cause

RCA Method:    5 Whys

Why 1: Why was there a KES 8,400 cash shortfall?
       → Cash collected exceeded POS receipts by KES 8,400

Why 2: Why did cash exceed POS receipts?
       → 3 transactions between 11:00–12:00 show cash taken but no POS entry

Why 3: Why were POS entries missing for those cash collections?
       → The POS terminal logged an error at 11:03 and the attendant
         continued collecting cash without logging to paper backup

Why 4: Why didn't the attendant use paper backup?
       → No paper backup protocol was followed or known

Why 5: Why is there no paper backup protocol?
       → The POS failure procedure is documented in the policy but was
         never included in staff onboarding training

Root Cause:    Training gap — POS failure procedure not communicated to staff

Contributing:  No supervisor check of POS-to-cash agreement during shift
               (control gap — supervisor review only at shift end, not mid-shift)
```

### 12.5 Corrective & Preventive Actions (CAPA)

```
Navigate: Governance → Incidents → [INC-2026-034] → CAPA

Corrective Actions (fix the immediate problem):
  CA-001:  Investigate and formally determine whether shortfall was
           theft or error; apply HR outcome accordingly.
           Owner: MD
           Due: 2026-07-05
           Status: IN PROGRESS

Preventive Actions (prevent recurrence):
  PA-001:  Add POS failure procedure to all staff onboarding training.
           Owner: Station Manager
           Due: 2026-07-15
           Evidence: Updated onboarding checklist + training records

  PA-002:  Configure AwoERP alert when POS-to-cash difference > KES 500
           at any point during a shift (not just at close).
           Owner: IT Admin
           Due: 2026-07-10
           Evidence: Screenshot of configured alert + test run

  PA-003:  Supervisors to perform mid-shift cash count and AwoERP
           reconciliation check at 2-hourly intervals.
           Owner: Station Manager (procedure update)
           Due: 2026-07-15
           Evidence: Updated shift procedure; signed acknowledgement by supervisors
```

### 12.6 Whistleblower Portal

AwoERP provides an anonymous whistleblower channel — accessible without login — for reporting concerns about fraud, misconduct, or ethical violations.

```
URL:  https://[tenant].awoerp.com/speak-up
      (also accessible via QR code posted at station)

Features:
  ✓ Fully anonymous (no login, no IP logging in submission record)
  ✓ Encrypted at rest and in transit
  ✓ Assignable reference code so reporter can check status without identifying
  ✓ Multi-language support: English, Swahili
  ✓ Mobile-optimised

Routing:
  All reports route to Audit Committee Chair + MD simultaneously
  If report implicates MD: routes to Board Chair only
  If report implicates Board Chair: routes to external audit firm contact

Triage:
  Initial response within 5 business days
  Formal investigation decision within 10 business days
  Reporter (via reference code) notified of outcome within 45 days
```

---

## 13. Reconciliation Health & Financial Controls

### 13.1 Reconciliation Subsystem

AwoERP's governance module maintains a **continuous reconciliation health subsystem** — monitoring the status, recency, and completeness of all active reconciliations across the tenant. This is surfaced as a real-time score on the Finance and Board dashboards.

### 13.2 GL Reconciliation Controls

```
Navigate: Governance → Reconciliations → GL Reconciliation

Control:  Every subsidiary ledger must reconcile to the GL at period close

Monitored Reconciliations:
  Accounts Payable (AP) subledger → GL Trade Creditors account
  Accounts Receivable (AR) subledger → GL Trade Debtors account
  Fixed Assets register → GL Fixed Asset accounts
  Payroll register → GL Payroll Liability accounts
  Stock module → GL Inventory accounts

Current Status (2026-06-22):
  AP Recon:        ✓ Reconciled (last run: 2026-06-20; difference: 0)
  AR Recon:        ⚠ NOT RECONCILED (last run: 2026-05-31; 22 days ago)
  Fixed Assets:    ✓ Reconciled (last run: 2026-06-01; difference: 0)
  Payroll:         ✓ Reconciled (last run: 2026-06-20; difference: 0)
  Stock-GL:        ✗ DIFFERENCE (last run: 2026-06-22; difference: KES 4,200)
```

### 13.3 Stock-GL Reconciliation

```
Navigate: Governance → Reconciliations → Stock-GL

Period:  2026-06-01 to 2026-06-22

Stock Module — Total Inventory Value:     KES 52,155,000
GL — Inventory Accounts (total):          KES 52,150,800
─────────────────────────────────────────────────────────
Difference:                               KES      4,200  ⚠

Difference Drill-Down:
  Account 1211-INVENTORY-PETROL:  SLE total = 48,920,000 | GL = 48,920,000  ✓
  Account 1212-INVENTORY-DIESEL:  SLE total = 2,345,000  | GL = 2,349,200   ✗ -4,200
  Account 1213-INVENTORY-LUBES:   SLE total =   890,000  | GL =   890,000   ✓

Likely Cause:   One stock entry on 2026-06-20 posted to GL but SLE failed
                due to a constraint error (batch not found). SLE repost required.

Resolution:     System Admin to run `stock:resync-gl --date 2026-06-20 --account 1212`
                Post-fix rerun reconciliation to confirm zero difference.
```

### 13.4 Cash & Bank Reconciliation

```
Navigate: Governance → Reconciliations → Cash & Bank

Bank Account:    Equity Bank Business — 0123456789
Statement Date:  2026-06-22

Bank Statement Closing Balance:    KES 3,847,220.00
AwoERP Bank Ledger Balance:        KES 3,831,450.00
─────────────────────────────────────────────────────
Unreconciled Difference:           KES    15,770.00

Outstanding Items:
  Deposits in Transit:
    2026-06-22  Cash deposit (in transit)   + KES 18,500.00
  Outstanding Cheques:
    CHQ-0341    Vivo Energy (issued 2026-06-20, not yet cleared) - KES  2,730.00

Adjusted Balance (AwoERP):        KES 3,847,220.00   ← agrees to bank ✓
Reconciliation Status:            RECONCILED
Last Reconciled By:               Accounts Officer
Reconciled Date:                  2026-06-22
```

### 13.5 Reconciliation Health Scoring

The Reconciliation Health Score is a composite metric (0–100) updated continuously:

```
Health Score Components:

Component                    Weight   Score   Weighted
───────────────────────────────────────────────────────
All recons run within SLA      30%     90      27.0
Zero unresolved differences    25%     60      15.0   ← stock-GL diff
Recons completed monthly       20%    100      20.0
AR recon currency              15%     40       6.0   ← 22 days overdue
AP recon currency              10%    100      10.0
───────────────────────────────────────────────────────
Total Health Score:                             78 / 100   AMBER

Thresholds:
  90–100:  GREEN  (Healthy)
  70–89:   AMBER  (Attention needed)
  50–69:   ORANGE (Action required)
  < 50:    RED    (Critical — escalate to MD)
```

---

## 14. TTL-Based Report Cache & Data Freshness

### 14.1 Architecture

Governance reports and dashboards aggregate data from multiple modules. To avoid expensive real-time queries on large datasets, AwoERP uses a **TTL-based report cache** backed by Redis.

```
┌──────────────────────┐
│   Dashboard / Report │
│   (amis-ui render)   │
└──────────┬───────────┘
           │ Request report data
           ▼
┌──────────────────────┐     Cache Hit (< TTL)
│   Report Cache Layer │ ──────────────────────► Return cached JSON
│   (Redis)            │
└──────────┬───────────┘
           │ Cache Miss or Expired
           ▼
┌──────────────────────┐
│ Report Compute Engine│  ← aggregates SLEs, workflow states,
│   (Go background)    │     compliance records, risk scores
└──────────┬───────────┘
           │
           ▼
    Store result in Redis
    with TTL + timestamp
    → Return to requester
```

### 14.2 Cache TTL Configuration

| Report Type | TTL | Invalidation Trigger |
|---|---|---|
| Board Dashboard KPIs | 300s (5 min) | Any high-severity governance event |
| Risk Heat Map | 600s (10 min) | Risk assessment update |
| Compliance Status | 600s (10 min) | Obligation status change |
| Reconciliation Health Score | 120s (2 min) | Any reconciliation run completion |
| Anomaly Alert Summary | 60s (1 min) | New anomaly event |
| Audit Finding Status | 300s (5 min) | Finding or MAP update |
| Control Effectiveness | 3600s (1 hr) | Control test result posted |
| KRI Dashboard | 300s (5 min) | KRI reading entered |

### 14.3 Data Freshness Indicators

Every cached report includes a `data_as_of` timestamp and a `freshness_status` indicator:

```json
{
  "report": "board_dashboard",
  "data_as_of": "2026-06-22T13:47:00+03:00",
  "ttl_seconds": 300,
  "freshness_status": "fresh",
  "cache_expires_at": "2026-06-22T13:52:00+03:00",
  "data": { ... }
}
```

The amis-ui frontend renders a subtle "data as of [time]" indicator on cached reports and a "Refresh" button to force a cache bypass.

---

## 15. Governance Dashboards & Reports

### 15.1 Board-Level Dashboard

Designed for directors and the MD — high-level, visual, exception-focused. Shows only what requires attention, not operational detail.

```
Navigate: Governance → Dashboards → Board Dashboard

┌─────────────────────────────────────────────────────────────────────┐
│  ANIKA GLOBAL LIMITED — BOARD GOVERNANCE DASHBOARD                  │
│  As of: 22 June 2026 | Data refreshed: 13:47 EAT                   │
├─────────────────┬──────────────────┬──────────────────┬────────────┤
│ RISK STATUS     │ COMPLIANCE       │ CONTROLS         │ INCIDENTS  │
│                 │                  │                  │            │
│ Critical:  0   │ Obligations: 47  │ Tested:     72% │ Open:    2 │
│ High:      3   │ Compliant:   44  │ Effective:  85% │ Critical: 1│
│ Medium:    8   │ Overdue:      1 ⚠│ Issues:      3 ⚠│ Avg Age: 6d│
│ Low:      12   │ Due 7 days:   3 ⚠│                  │            │
└─────────────────┴──────────────────┴──────────────────┴────────────┘

ITEMS REQUIRING BOARD ATTENTION:
  ⚠ [HIGH]     NEMA Annual Environmental Audit overdue 83 days
  ⚠ [HIGH]     VAT Return for June 2026 due in 28 days
  ⚠ [MEDIUM]   AR Subledger reconciliation overdue 22 days
  ℹ [INFO]     Q2 Internal Audit complete — 2 findings, 1 open MAP

RECONCILIATION HEALTH:  78/100  AMBER
ANOMALY ALERTS (30d):   Open 3  |  Investigating 1  |  Resolved 28
```

### 15.2 Risk Dashboard

```
Navigate: Governance → Dashboards → Risk Dashboard

Summary:
  Total Risks:          23
  By Residual Score:    Critical 0  |  High 3  |  Medium 8  |  Low 12
  Within Appetite:      20 (87%)
  Exceeding Appetite:    3 (13%)  ⚠ — require MD escalation

Heat Map:     [Interactive 5×5 grid — click cell to see risk list]

KRI Status:
  KRI-FORE-001  Wetstock Variance    0.23%   GREEN ✓
  KRI-FIN-002   AP Days Outstanding  42 days  AMBER ⚠ (threshold: 30 days)
  KRI-COMP-001  Compliance Rate      93.6%   GREEN ✓
  KRI-AUD-001   Open Findings        3        AMBER ⚠ (threshold: 5)

Risks Due for Review This Quarter:
  RISK-FRAUD-001  Review due 2026-09-01  (78 days)
  RISK-OPS-003    Review due 2026-07-15  (23 days)  — UPCOMING
```

### 15.3 Compliance Status Report

```
Navigate: Governance → Reports → Compliance Status

Compliance Summary — June 2026

Total Obligations:      47
  ✓ Compliant:         44  (93.6%)
  ⚠ Overdue:            1  (2.1%)   — NEMA Annual Audit
  ↑ Due Within 7 Days:  2  (4.3%)   — VAT, PAYE

Compliance by Category:
  KRA Tax:          12/12  100%  ✓
  EPRA Petroleum:    8/9   88.9% ⚠ (VRU maintenance overdue)
  NEMA Environmental: 5/6  83.3% ⚠ (Annual audit overdue)
  Companies Act:     6/6   100%  ✓
  Employment:        8/8   100%  ✓
  Data Protection:   3/3   100%  ✓
  Internal Policies: 3/3   100%  ✓
```

### 15.4 Audit Progress Tracker

```
Navigate: Governance → Reports → Audit Progress

2026 Audit Plan Progress:

Engagement                       Status          Findings  Open MAPs
────────────────────────────────────────────────────────────────────
AU-005  Stock Reconciliation      COMPLETE         1         0 ✓
AU-001  Forecourt Cash Handling   IN PROGRESS      1         3
AU-004  IT Access Controls        PLANNED (Q2)     —         —
AU-002  Fuel Procurement          PLANNED (Q3)     —         —
AU-007  Environmental Compliance  PLANNED (Q3)     —         —
AU-003  Payroll & HR              PLANNED (Q4)     —         —

Open MAPs by Age:
  0–30 days:   3 MAPs   (acceptable)
  31–60 days:  0 MAPs
  60+ days:    1 MAP ⚠  (FIND-2025-004-02 — escalated to Board)

Audit Coverage: 2/6 units (33%) complete or in progress
```

### 15.5 Anomaly Alert Summary

```
Navigate: Governance → Reports → Anomaly Alert Summary

Rolling 30 Days (2026-05-23 to 2026-06-22):

Total Alerts:     32
  Critical:        0
  High:            7
  Medium:         18
  Low:             7

By Type:
  Statistical Outlier:        12  (8 resolved, 4 dismissed as seasonal)
  Velocity/Acceleration:       6  (5 resolved, 1 investigating)
  Behavioural:                 4  (3 resolved, 1 open)
  Amount Sequence (Benford):   2  (2 open — audit referral pending)
  Cross-Tenant Spread:         4  (4 dismissed — price movement known)
  SoD Violation (soft):        3  (3 documented exceptions)
  Reconciliation Health:       1  (open — stock-GL difference)

Mean Time to Acknowledge:  1.8 hours
Mean Time to Resolve:      28.4 hours
False Positive Rate:        31%  (10 of 32 dismissed)
```

### 15.6 Control Effectiveness Report

```
Navigate: Governance → Reports → Control Effectiveness

Period:  2026-Q2 Testing

Controls Tested:        18
Controls Due But Not Tested:  4  ⚠

Effectiveness Distribution:
  Effective:            12 (67%)
  Mostly Effective:      3 (17%)
  Partially Effective:   2 (11%)
  Ineffective:           1 (5%)   ← CTRL-CASH-003 (Cash Recon)

Controls with Issues:
  CTRL-CASH-003   Partially Effective → Finding FIND-2026-002-01 raised
  CTRL-IT-004     Mostly Effective → recommendation to strengthen access logs
  CTRL-PROC-002   Mostly Effective → training gap identified

Attestations Received:   14/18 Control Owners attested  (Q2 deadline: 2026-07-15)
  Missing Attestations:   4  ⚠ — automated reminders sent
```

---

## 16. Generalised Audit Log

### 16.1 Audit Log Architecture

AwoERP implements a **fully generalised, configuration-table-driven audit log** that captures state changes across every module. Inspired by SAP's change log architecture but implemented as a native PostgreSQL + Go solution, the audit log is the foundation of governance accountability.

Every write operation in AwoERP — create, update, submit, amend, cancel, approve, reject, delete — passes through a common audit middleware that records:

```sql
-- audit.event_log
CREATE TABLE audit.event_log (
  event_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id          UUID NOT NULL,
  occurred_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  entity_type        VARCHAR(100) NOT NULL,  -- 'stock_entry', 'purchase_order', etc.
  entity_id          UUID NOT NULL,
  action             VARCHAR(50) NOT NULL,   -- 'create', 'submit', 'approve', 'cancel', ...
  actor_id           UUID NOT NULL,
  actor_email        VARCHAR(255),
  actor_ip           INET,
  actor_user_agent   TEXT,
  session_id         UUID,
  before_state       JSONB,                  -- full document snapshot before change
  after_state        JSONB,                  -- full document snapshot after change
  diff               JSONB,                  -- computed field-level diff
  sensitivity        sensitivity_level,      -- enum: low, medium, high, critical
  source_module      VARCHAR(50),
  correlation_id     UUID,                   -- ties related events across a workflow
  retention_until    DATE NOT NULL,          -- computed from sensitivity + regulation
  is_tamper_evident  BOOLEAN DEFAULT TRUE,
  row_hash           VARCHAR(64)             -- SHA-256 of the record for tamper detection
);
```

### 16.2 Sensitivity Classification

Sensitivity is determined automatically by the configuration table — not hardcoded. This allows administrators to reclassify events without a code change.

```sql
-- audit.log_sensitivity_config
CREATE TABLE audit.log_sensitivity_config (
  config_id      UUID PRIMARY KEY,
  tenant_id      UUID NOT NULL,
  entity_type    VARCHAR(100),   -- null = matches all
  action         VARCHAR(50),    -- null = matches all
  sensitivity    sensitivity_level NOT NULL,
  retention_days INT NOT NULL,
  requires_mfa   BOOLEAN DEFAULT FALSE,   -- force MFA before this action is permitted
  alert_on_occur BOOLEAN DEFAULT FALSE,   -- trigger governance event on every occurrence
  notes          TEXT
);

-- Example records:
INSERT INTO audit.log_sensitivity_config VALUES
  (gen_random_uuid(), 'tenant-001', 'stock_entry', 'cancel',    'high',     2555, FALSE, TRUE, 'Cancellations must be investigated'),
  (gen_random_uuid(), 'tenant-001', 'payment',     'approve',   'high',     2555, TRUE,  FALSE, 'Payment approvals require MFA'),
  (gen_random_uuid(), 'tenant-001', 'user',        'delete',    'critical', 2555, TRUE,  TRUE,  'User deletion always alerted'),
  (gen_random_uuid(), 'tenant-001', NULL,           'login',     'low',       365, FALSE, FALSE, 'Standard login logs'),
  (gen_random_uuid(), 'tenant-001', 'role',        'assign',    'critical', 2555, TRUE,  TRUE,  'Role changes are critical events');
```

### 16.3 Log Retention & Archival

| Sensitivity | Retention Period | Regulation Basis | Storage |
|---|---|---|---|
| Low | 1 year | General operational records | Hot (PostgreSQL) |
| Medium | 3 years | Standard business records | Hot then warm (compressed) |
| High | 7 years | KRA tax records, Companies Act | Warm (object storage) |
| Critical | 10 years | Financial crime, fraud evidence | Cold (encrypted archive) |

Records approaching their retention boundary are archived to immutable object storage (e.g., AWS S3 Glacier, local Wasabi bucket) before deletion from the operational database. A manifest of archived records is retained permanently.

### 16.4 Tamper-Evidence

Each audit log row carries a `row_hash` — a SHA-256 hash of the record's content fields. A background integrity checker periodically recomputes hashes and alerts if any mismatch is detected (indicating tampering).

```go
// Hash computation (Go)
func ComputeRowHash(e AuditEvent) string {
    data := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%v|%v",
        e.EventID, e.TenantID, e.OccurredAt.UTC().Format(time.RFC3339Nano),
        e.EntityType, e.EntityID, e.Action, e.BeforeState, e.AfterState)
    hash := sha256.Sum256([]byte(data))
    return hex.EncodeToString(hash[:])
}
```

Audit log rows are **never updated or deleted** through application code. Retention-based purges are performed by a privileged background process that writes a purge manifest before deletion.

---

## 17. Feature Flags & Permissions

### 17.1 Governance Feature Flags

```
Navigate: Settings → Feature Flags → Governance

Flag Key                                 Default  Description
──────────────────────────────────────────────────────────────────────────────
governance.module.enabled                false    Master governance enable
governance.risk_management               false    Risk register, KRIs, heat map
governance.internal_controls             false    Control library and testing
governance.sod_enforcement               false    Segregation of duties rules
governance.policy_management             false    Policy registry and acknowledgements
governance.approval_workflows            true     Basic approval routing (always available)
governance.doa_matrix                    false    Delegation of authority matrix
governance.anomaly_detection             true     Continuous anomaly monitoring
governance.anomaly.benford_law           false    Benford's Law amount analysis
governance.anomaly.cross_tenant          false    Cross-tenant spread analysis (opt-in)
governance.audit_management             false    Audit plan, engagements, findings
governance.compliance_management         false    Obligations register, evidence vault
governance.incident_management           false    Incident intake, investigation, CAPA
governance.whistleblower_portal          false    Anonymous reporting portal
governance.reconciliation_health         true     Health score (always on when recons active)
governance.audit_log.full_diff           false    Capture field-level diff in audit log
governance.audit_log.mfa_on_critical     false    Require MFA for critical actions
governance.board_dashboard               false    Board-level governance dashboard
governance.report_cache.enabled          true     TTL-based report cache (always on)
```

### 17.2 Role-Based Access Control

#### Governance Permission Matrix

| Action | Board | Audit Chair | Compliance | Risk Mgr | Auditor | Control Owner | Admin |
|---|---|---|---|---|---|---|---|
| View Board Dashboard | ✓ | ✓ | ✓ | ✓ | ✗ | ✗ | ✓ |
| View Risk Register | ✓ | ✓ | ✓ | ✓ | ✓ | Own | ✓ |
| Edit Risk Register | ✗ | ✗ | ✗ | ✓ | ✗ | ✗ | ✓ |
| Approve Risk Appetite | ✓ | ✓ | ✗ | ✗ | ✗ | ✗ | ✗ |
| View Control Library | ✓ | ✓ | ✓ | ✓ | ✓ | Own | ✓ |
| Create Control Test | ✗ | ✗ | ✗ | ✗ | ✓ | Own | ✓ |
| Attest Control | ✗ | ✗ | ✗ | ✗ | ✗ | Own | ✓ |
| Manage SoD Rules | ✗ | ✓ | ✗ | ✗ | ✗ | ✗ | ✓ |
| View Anomaly Alerts | ✗ | ✓ | ✓ | ✓ | ✓ | ✗ | ✓ |
| Dismiss Anomaly | ✗ | ✓ | ✗ | ✓ | ✓ | ✗ | ✓ |
| Manage Audit Plans | ✗ | ✓ | ✗ | ✗ | ✓ | ✗ | ✓ |
| Raise Audit Finding | ✗ | ✗ | ✗ | ✗ | ✓ | ✗ | ✓ |
| Accept MAP | ✗ | ✓ | ✗ | ✗ | ✗ | ✗ | ✗ |
| Manage Compliance Obligations | ✗ | ✓ | ✓ | ✗ | ✗ | ✗ | ✓ |
| Upload Compliance Evidence | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ | ✓ |
| Report Incident | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| Manage Incident Investigation | ✗ | ✓ | ✗ | ✗ | ✓ | ✗ | ✓ |
| View Audit Log | ✗ | ✓ | ✗ | ✗ | ✓ | ✗ | ✓ |
| Configure Feature Flags | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ |

---

## 18. API Reference & Integration

### 18.1 REST API Endpoints

All endpoints are tenant-scoped and session-authenticated.

#### Governance Events

```
GET    /api/v1/governance/events
       ?severity=high,critical&status=open&page=1&per_page=50

POST   /api/v1/governance/events/{event_id}/acknowledge
POST   /api/v1/governance/events/{event_id}/resolve
POST   /api/v1/governance/events/{event_id}/dismiss
```

#### Risk Register

```
GET    /api/v1/governance/risks
GET    /api/v1/governance/risks/{risk_id}
POST   /api/v1/governance/risks
PUT    /api/v1/governance/risks/{risk_id}/assess
GET    /api/v1/governance/risks/heat-map

GET    /api/v1/governance/kris
GET    /api/v1/governance/kris/{kri_id}/readings
POST   /api/v1/governance/kris/{kri_id}/readings
```

#### Approval Workflows

```
GET    /api/v1/governance/workflows/pending
       ?assigned_to_me=true

POST   /api/v1/governance/workflows/{instance_id}/approve
       Body: { "comments": "Approved — within budget" }

POST   /api/v1/governance/workflows/{instance_id}/reject
       Body: { "comments": "Exceeds department budget; resubmit with MD override" }

POST   /api/v1/governance/workflows/{instance_id}/delegate
       Body: { "delegate_to": "user-uuid", "reason": "On leave" }
```

#### Compliance

```
GET    /api/v1/governance/compliance/obligations
       ?status=overdue&category=kra

GET    /api/v1/governance/compliance/calendar
       ?from=2026-07-01&to=2026-07-31

POST   /api/v1/governance/compliance/evidence
       (multipart/form-data — obligation_id + file upload)
```

#### Audit Log

```
GET    /api/v1/governance/audit-log
       ?entity_type=stock_entry&action=cancel&from=2026-06-01&to=2026-06-30

GET    /api/v1/governance/audit-log/{event_id}

POST   /api/v1/governance/audit-log/integrity-check
       (Triggers hash verification for a date range — gov.admin role only)
```

#### Reconciliation Health

```
GET    /api/v1/governance/reconciliation/health
       → Returns composite health score + component breakdown

GET    /api/v1/governance/reconciliation/stock-gl
       ?period=2026-06

GET    /api/v1/governance/reconciliation/cash-bank
       ?account_id=uuid&statement_date=2026-06-22
```

### 18.2 Webhook Events

| Event | Trigger | Payload |
|---|---|---|
| `governance.anomaly.raised` | New anomaly signal created | anomaly_type, severity, entity, details |
| `governance.workflow.approval_required` | Approval step assigned | workflow_id, document, amount, approver |
| `governance.workflow.approved` | Workflow step approved | workflow_id, approver, comments |
| `governance.workflow.rejected` | Workflow rejected | workflow_id, rejector, reason |
| `governance.workflow.escalated` | Timeout escalation fired | workflow_id, escalated_to, reason |
| `governance.compliance.due_soon` | Obligation due within threshold | obligation_id, due_date, owner |
| `governance.compliance.overdue` | Obligation past due date | obligation_id, days_overdue |
| `governance.incident.raised` | New incident created | incident_id, category, severity |
| `governance.sod.violation` | SoD rule triggered | rule_id, actor, document, blocked/warned |
| `governance.risk.appetite_exceeded` | Residual risk exceeds appetite | risk_id, score, appetite_threshold |
| `governance.reconciliation.health_changed` | Health score crosses threshold | score, previous_score, band |
| `governance.audit.finding_raised` | New audit finding created | finding_id, severity, owner |
| `governance.audit.map_overdue` | MAP action past due date | map_id, finding_id, days_overdue |

### 18.3 External GRC Integration

AwoERP can push governance data to external GRC platforms used by enterprise clients or group holding companies.

#### SAP GRC Integration

```
Navigate: Governance → Integrations → SAP GRC

Supported Sync:
  ↗ Risks         → SAP GRC Risk Management
  ↗ Controls      → SAP GRC Process Control
  ↗ Audit Findings → SAP Audit Management
  ↗ Compliance    → SAP Compliance Management

Auth Method:    OAuth 2.0 (SAP BTP)
Sync Frequency: Hourly (configurable)
Conflict Resolution: AwoERP as master; SAP is read-only recipient
```

#### Microsoft Sentinel / SIEM Integration

```
Navigate: Governance → Integrations → SIEM

Supported Outputs:
  AwoERP audit log stream → Azure Event Hub → Microsoft Sentinel
  Anomaly alerts → Splunk via HTTP Event Collector

Format: JSON (CEF-compatible)
Authentication: API key + IP allowlist
```

---

## 19. Regulatory Context — Kenya & East Africa

### 19.1 Key Kenyan Regulations Impacting ERP Governance

| Regulation | Key Governance Implication | AwoERP Module |
|---|---|---|
| **Companies Act, Cap 486 (2015)** | Annual financial statements; director duties; company secretary records | Compliance, Audit Log |
| **Income Tax Act, Cap 470** | Corporate tax filing, instalment payments, withholding tax | Compliance Calendar |
| **Value Added Tax Act (2013)** | Monthly VAT returns; iTax e-filing; input/output VAT reconciliation | Compliance, Finance |
| **Tax Procedures Act (2015)** | 7-year record retention; KRA audit access rights | Audit Log Retention |
| **Employment Act, Cap 226** | PAYE, NSSF, NHIF remittance; disciplinary procedures | Compliance, HR |
| **Energy Act (2019)** | EPRA licensing; petroleum retail standards; pump calibration | EPRA Compliance |
| **NEMA Act (1999)** | Environmental impact assessments; spill reporting; EMP | Environmental Compliance |
| **Data Protection Act (2019)** | Personal data handling; breach notification within 72 hours | Incident Management, Audit Log |
| **OHSA (2007)** | Workplace safety; incident reporting; training records | Incident Management |

### 19.2 KRA Digital Compliance

KRA's iTax platform and the planned **eTIMS (Electronic Tax Invoice Management System)** integration are central to AwoERP's tax compliance posture:

- All sales invoices generate an eTIMS-compatible electronic receipt
- VAT reconciliation data can be exported directly to iTax-compatible format
- The compliance calendar enforces all KRA deadline dates with pre-alerts
- Audit log records are retained for 7 years in compliance with the Tax Procedures Act

### 19.3 East African Community (EAC) Harmonisation

For businesses operating across multiple EAC member states (Kenya, Uganda, Tanzania, Rwanda, Burundi, South Sudan, DRC), AwoERP's multi-tenant architecture supports:

- Per-entity compliance calendars reflecting each country's tax regime
- Country-specific Chart of Accounts templates (Kenya, Uganda, Tanzania seeded)
- EAC Customs Union documentation for cross-border goods movements

---

## 20. Troubleshooting & FAQs

### Q: An approval notification was sent but the approver says they didn't receive it

**Cause:** Email delivery failure (spam filter, incorrect email address, SMTP quota) or SMS delivery failure.

**Resolution:**
1. Navigate to **Governance → Approval Workflows → [Instance] → Notification Log**. This shows every notification attempted, the channel, and delivery status.
2. If delivery failed, use **Resend Notification** to retry.
3. For persistent email failures, check SMTP settings under **Settings → Email**.
4. The approver can also access all pending approvals directly at **Governance → My Approvals** without relying on notification delivery.

### Q: An anomaly alert was raised for something that is normal for our business (false positive)

**Cause:** The statistical baseline does not yet reflect a legitimate business pattern (e.g., a seasonal spike, a one-time event, a new product line).

**Resolution:**
1. Navigate to **Governance → Anomaly Alerts → [Alert] → Dismiss**
2. Record a detailed dismissal reason (this feeds back into the model's training set)
3. If this is a known recurring pattern, adjust the detection thresholds for this specific KRI/metric under **Governance → Anomaly Settings → Custom Thresholds**
4. After 3+ dismissals of the same type with the same reason, the system will suggest creating a business-rule exception

### Q: A user bypassed SoD by using a colleague's credentials

**Cause:** Credential sharing — a human control failure that no technical SoD rule can fully prevent.

**Resolution:**
1. The audit log records the IP address and user-agent of every login. Cross-reference the anomalous approval with login records to identify shared credential use.
2. Create an incident (**Governance → Incidents → New**) classifying as `INC-SEC` (security) and `INC-HR` (conduct) simultaneously.
3. Reset the compromised user's password and session tokens immediately.
4. Review CTRL-IT-004 (IT Access Controls) and add user awareness training on credential sharing as a CAPA.

### Q: The Reconciliation Health Score dropped suddenly

**Cause:** A previously reconciled account now shows a difference, or a reconciliation SLA has been missed.

**Resolution:**
1. Navigate to **Governance → Reconciliations → Health Score** and click on any component showing red/amber.
2. The drill-down shows which specific reconciliation is contributing to the score drop.
3. Resolve the underlying reconciliation issue (see Stock module or Finance module troubleshooting).
4. Run the reconciliation again — the health score updates within 2 minutes.

### Q: The compliance obligation shows as overdue but we completed it

**Cause:** The obligation was completed but evidence was not uploaded and verified in AwoERP.

**Resolution:**
1. Navigate to **Governance → Compliance → Obligations → [Obligation] → Evidence tab**
2. Upload the completion evidence (iTax receipt, signed document, etc.)
3. Mark the evidence as submitted — it will be routed to the Compliance Officer for verification
4. Once verified, the obligation status automatically updates to Compliant

### Q: The Board Dashboard is showing stale data

**Cause:** The report cache TTL has not expired since the underlying data changed, or a cache invalidation signal was missed.

**Resolution:**
1. Click the **Refresh** button on the dashboard (forces cache bypass for the current session)
2. If data is persistently stale beyond the expected TTL, navigate to **Settings → Cache → Invalidate → Governance Reports** (Governance Admin role required)
3. For persistent issues, check the Redis connection status under **Settings → System Health**

---

## 21. Glossary

| Term | Definition |
|---|---|
| **Anomaly** | A data point or pattern deviating significantly from established norms, detected by the automated anomaly engine |
| **Audit Committee** | A sub-committee of the Board responsible for overseeing internal audit, financial reporting integrity, and compliance |
| **Audit Finding** | A documented observation from an audit engagement, classified by severity and requiring management response |
| **Audit Universe** | The complete inventory of all auditable units across the organisation |
| **CAPA** | Corrective and Preventive Actions — the structured response to an incident or audit finding |
| **CCM** | Continuous Controls Monitoring — automated, ongoing evaluation of controls without human-triggered tests |
| **Compensating Control** | An alternative control that mitigates risk when the primary control cannot be applied |
| **Control** | A policy, procedure, or system feature that reduces the likelihood or impact of a risk |
| **Control Attestation** | A formal declaration by a control owner that a control operated effectively during a period |
| **COSO** | Committee of Sponsoring Organizations — a widely used internal controls framework |
| **DoA** | Delegation of Authority — the formal matrix defining who can approve what within what limits |
| **eTIMS** | KRA's Electronic Tax Invoice Management System |
| **EPRA** | Energy and Petroleum Regulatory Authority — Kenya's petroleum sector regulator |
| **Governance Event** | A structured signal emitted by any module and consumed by the Governance module for monitoring |
| **iTax** | KRA's online tax management portal for filing returns and making payments |
| **KPI** | Key Performance Indicator — a metric for measuring business performance |
| **KRA** | Kenya Revenue Authority — Kenya's national tax collection agency |
| **KRI** | Key Risk Indicator — a leading metric signalling increasing risk exposure |
| **MAP** | Management Action Plan — management's committed response to an audit finding |
| **NEMA** | National Environment Management Authority — Kenya's environmental regulator |
| **Reconciliation Health Score** | A composite 0–100 metric reflecting the currency and completeness of all active reconciliations |
| **Regulatory Universe** | The master catalogue of all laws, regulations, and standards applicable to the organisation |
| **Residual Risk** | The remaining risk level after controls have been applied |
| **Risk Appetite** | The amount of risk an organisation is willing to accept in pursuit of its objectives |
| **Risk Register** | The central repository of all identified risks, with assessments, owners, and controls |
| **SoD** | Segregation of Duties — the control that ensures no single person can execute a complete transaction without oversight |
| **Temporal** | The workflow orchestration engine used by AwoERP for durable, long-running processes |
| **TTL** | Time to Live — the duration for which a cached value remains valid before being refreshed |
| **Whistleblower** | A person who reports concerns about wrongdoing, typically anonymously and with legal protections |

---

*This document is maintained as part of the AwoERP platform documentation suite. For the developer API reference, integration guides, and release changelog, refer to the `docs/` directory in the AwoERP monorepo. For support, contact the AwoERP platform team or your designated compliance consultant.*

*© 2026 Anika Global Limited / AwoERP Platform. All rights reserved.*
