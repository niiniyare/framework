# AwoERP Advanced Stock Management Module
## Comprehensive Documentation Guide

**Module:** `awoerp/stock`
**Version:** 1.0.0
**Platform:** AwoERP — Multi-Tenant SaaS ERP for East African Businesses
**Audience:** System Administrators, Warehouse Managers, Operations Teams, Developers
**Last Updated:** June 2026

---

## Table of Contents

1. [Introduction & Module Overview](#1-introduction--module-overview)
2. [Core Concepts & Terminology](#2-core-concepts--terminology)
3. [System Architecture](#3-system-architecture)
4. [Initial Setup & Configuration](#4-initial-setup--configuration)
   - 4.1 [Tenant-Level Stock Settings](#41-tenant-level-stock-settings)
   - 4.2 [Warehouse Hierarchy Setup](#42-warehouse-hierarchy-setup)
   - 4.3 [Item Master Configuration](#43-item-master-configuration)
   - 4.4 [Unit of Measure (UOM) Setup](#44-unit-of-measure-uom-setup)
   - 4.5 [Item Groups & Categories](#45-item-groups--categories)
   - 4.6 [Stock Valuation Methods](#46-stock-valuation-methods)
5. [Warehouse Management](#5-warehouse-management)
   - 5.1 [Warehouse Types](#51-warehouse-types)
   - 5.2 [Storage Locations & Bins](#52-storage-locations--bins)
   - 5.3 [Warehouse Permissions](#53-warehouse-permissions)
   - 5.4 [Virtual Warehouses](#54-virtual-warehouses)
6. [Item Master & Variants](#6-item-master--variants)
   - 6.1 [Item Attributes & Variants](#61-item-attributes--variants)
   - 6.2 [Item Pricing](#62-item-pricing)
   - 6.3 [Reorder Rules](#63-reorder-rules)
   - 6.4 [Item Bundling & Kitting](#64-item-bundling--kitting)
7. [Stock Transactions](#7-stock-transactions)
   - 7.1 [Stock Entry Types](#71-stock-entry-types)
   - 7.2 [Purchase Receipts](#72-purchase-receipts)
   - 7.3 [Delivery Notes](#73-delivery-notes)
   - 7.4 [Stock Transfers](#74-stock-transfers)
   - 7.5 [Material Issues & Returns](#75-material-issues--returns)
   - 7.6 [Manufacturing Entries](#76-manufacturing-entries)
8. [Inventory Valuation & Costing](#8-inventory-valuation--costing)
   - 8.1 [FIFO Valuation](#81-fifo-valuation)
   - 8.2 [AVCO (Weighted Average Cost)](#82-avco-weighted-average-cost)
   - 8.3 [Standard Costing](#83-standard-costing)
   - 8.4 [Landed Costs](#84-landed-costs)
   - 8.5 [Cost Centre Allocation](#85-cost-centre-allocation)
9. [Batch & Serial Number Tracking](#9-batch--serial-number-tracking)
   - 9.1 [Batch Management](#91-batch-management)
   - 9.2 [Serial Number Tracking](#92-serial-number-tracking)
   - 9.3 [Expiry Date Management](#93-expiry-date-management)
   - 9.4 [Batch-Level Costing](#94-batch-level-costing)
10. [Stock Reconciliation & Physical Count](#10-stock-reconciliation--physical-count)
    - 10.1 [Cycle Counting](#101-cycle-counting)
    - 10.2 [Full Physical Inventory](#102-full-physical-inventory)
    - 10.3 [Variance Analysis](#103-variance-analysis)
    - 10.4 [Reconciliation Approval Workflows](#104-reconciliation-approval-workflows)
11. [Demand Planning & Replenishment](#11-demand-planning--replenishment)
    - 11.1 [Reorder Point Planning](#111-reorder-point-planning)
    - 11.2 [Min/Max Replenishment](#112-minmax-replenishment)
    - 11.3 [Material Requirement Planning (MRP)](#113-material-requirement-planning-mrp)
    - 11.4 [Automated Purchase Requests](#114-automated-purchase-requests)
12. [Multi-Location & Multi-Currency](#12-multi-location--multi-currency)
    - 12.1 [Multi-Warehouse Operations](#121-multi-warehouse-operations)
    - 12.2 [Inter-Company Transfers](#122-inter-company-transfers)
    - 12.3 [Foreign Currency Inventory](#123-foreign-currency-inventory)
13. [Quality Control Integration](#13-quality-control-integration)
    - 13.1 [Inspection on Receipt](#131-inspection-on-receipt)
    - 13.2 [Quarantine Warehouses](#132-quarantine-warehouses)
    - 13.3 [Quality Hold & Release](#133-quality-hold--release)
14. [Petroleum & Forecourt Stock (Industry Extension)](#14-petroleum--forecourt-stock-industry-extension)
    - 14.1 [Wetstock Management](#141-wetstock-management)
    - 14.2 [Tank Dip Readings](#142-tank-dip-readings)
    - 14.3 [Meter-Based Dispensing](#143-meter-based-dispensing)
    - 14.4 [Variance & Loss Reconciliation](#144-variance--loss-reconciliation)
    - 14.5 [Environmental Compliance Tracking](#145-environmental-compliance-tracking)
15. [Stock Reports & Analytics](#15-stock-reports--analytics)
    - 15.1 [Standard Reports](#151-standard-reports)
    - 15.2 [Custom Report Builder](#152-custom-report-builder)
    - 15.3 [KPI Dashboard](#153-kpi-dashboard)
    - 15.4 [Anomaly Detection Alerts](#154-anomaly-detection-alerts)
16. [API Reference & Integration](#16-api-reference--integration)
    - 16.1 [REST API Endpoints](#161-rest-api-endpoints)
    - 16.2 [Webhook Events](#162-webhook-events)
    - 16.3 [Third-Party Integrations](#163-third-party-integrations)
17. [Feature Flags & Permissions](#17-feature-flags--permissions)
    - 17.1 [Stock Module Feature Flags](#171-stock-module-feature-flags)
    - 17.2 [Role-Based Access Control](#172-role-based-access-control)
18. [Audit Trail & Compliance](#18-audit-trail--compliance)
19. [Troubleshooting & FAQs](#19-troubleshooting--faqs)
20. [Glossary](#20-glossary)

---

## 1. Introduction & Module Overview

The **AwoERP Advanced Stock Management Module** (`awoerp/stock`) is a purpose-built inventory and warehousing engine designed for the operational realities of East African businesses. Taking inspiration from battle-tested ERP patterns (ERPNext's warehouse-centric design and QuickBooks' accessibility for SMEs), AwoERP extends these conventions with multi-tenant isolation, schema-driven UI rendering, and industry-specific extensions for petroleum retail, general trade, restaurants, and aviation.

### 1.1 What This Module Does

The stock module governs the entire lifecycle of physical goods from procurement through consumption or sale:

- Tracks real-time quantities across an unlimited number of warehouses and sub-locations
- Maintains accurate inventory valuation using FIFO, AVCO, or Standard Cost methods
- Provides batch and serial number traceability with full forwards/backwards trace
- Enforces configurable approval workflows for high-value or sensitive transactions
- Emits structured, permission-gated JSON consumed by the amis-ui web client, Flutter mobile apps, and any third-party integration
- Surfaces anomaly alerts for suspicious stock movements, velocity spikes, or valuation outliers

### 1.2 Design Philosophy

AwoERP's stock module is built on three principles:

**Server-Authoritative Truth.** The Go backend is the sole source of fact for stock quantities, valuations, and transaction history. The frontend never computes stock — it renders what the server returns. This eliminates client-side drift common in hybrid-sync architectures.

**Tenant Isolation by Default.** Every table carrying stock data is partitioned by `tenant_id`. No query crosses tenant boundaries. Feature flags, valuation methods, warehouse trees, and approval workflows are independently configured per tenant.

**Traceable, Reversible Ledger.** Stock transactions write to an append-only ledger. There are no in-place quantity updates. Corrections are made through reversal entries, not deletions. This satisfies KRA audit requirements and IFRS-aligned valuation reporting.

### 1.3 Module Boundaries

The Stock module interfaces with the following AwoERP modules:

| Interface | Direction | Description |
|---|---|---|
| **Purchasing** | Inbound | Purchase Orders trigger receipt entries |
| **Sales** | Outbound | Sales Orders trigger delivery/issue entries |
| **Finance / GL** | Outbound | Every stock movement posts to the ledger via account mappings |
| **Manufacturing** | Bidirectional | BOM explosions consume stock; finished goods receipts add stock |
| **HR / Shifts** | Inbound | Forecourt shifts drive wetstock period boundaries |
| **Governance** | Outbound | Anomaly events published to governance module for review |

---

## 2. Core Concepts & Terminology

Understanding these foundational concepts is essential before configuring the module.

### Stock Ledger Entry (SLE)

The atomic unit of stock accounting. Every movement — receipt, issue, transfer, adjustment — creates one or more SLEs. An SLE records:

```
sle_id          UUID (primary key)
tenant_id       UUID
posting_datetime TIMESTAMPTZ
item_code       VARCHAR
warehouse_id    UUID
qty_change      DECIMAL(18,6)   -- positive = in, negative = out
valuation_rate  DECIMAL(18,6)
stock_value     DECIMAL(18,6)   -- qty_change × valuation_rate
voucher_type    ENUM
voucher_no      VARCHAR
batch_no        VARCHAR (nullable)
serial_no       VARCHAR (nullable)
created_by      UUID
```

SLEs are **never deleted or updated**. Corrections create offsetting entries.

### Valuation Rate

The per-unit cost of an item at the time of a transaction. How this is computed depends on the item's configured valuation method (FIFO / AVCO / Standard).

### Warehouse

A logical or physical location where stock is held. Warehouses form a tree (a root company warehouse may contain site warehouses, which may contain sub-stores and bins). Quantities are tracked at leaf-node warehouses only; parent nodes aggregate for reporting.

### Item

A unique SKU tracked in the system. Items are either **Stocked** (quantities tracked), **Non-Stocked** (services, expenses), or **Fixed Assets** (tracked separately).

### Batch

A group of units of the same item received in the same consignment, sharing manufacturing date, expiry date, and supplier lot number. Batches enable FEFO dispensing and recall traceability.

### Serial Number

A unique identifier for a single unit of an item. Used for high-value equipment, cylinders, and warranty-tracked goods.

### Bin

The leaf-level physical location within a warehouse (shelf, rack, tank, bay). Stock can optionally be tracked at bin level below the warehouse level.

---

## 3. System Architecture

### 3.1 Data Flow

```
External Events (Purchase, Sale, Manufacture, Adjustment)
        │
        ▼
  Transaction Service (Go/Fiber)
        │
        ├── Validates item, warehouse, UOM, batch
        ├── Computes valuation rate (FIFO / AVCO / Standard)
        ├── Writes Stock Ledger Entries (PostgreSQL)
        ├── Updates bin quantities (Redis cache → PostgreSQL)
        ├── Posts GL entries to Finance module
        ├── Publishes domain events (Temporal workflows)
        └── Emits anomaly signals (Governance module)
        │
        ▼
  Query Service (Go/Fiber)
        │
        ├── Aggregates SLEs for stock balance queries
        ├── Applies permission gates (role + feature flags)
        └── Emits amis-compatible JSON for UI rendering
```

### 3.2 PostgreSQL Schema Overview

The core tables in the `stock` schema:

```sql
-- Item master
stock.items
stock.item_variants
stock.item_attributes
stock.item_attribute_values
stock.item_uoms

-- Warehouse tree
stock.warehouses
stock.bins

-- Batches & serials
stock.batches
stock.serial_numbers

-- Ledger
stock.stock_ledger_entries      -- append-only, partitioned by posting_datetime
stock.stock_balance             -- materialised view, refreshed on SLE insert

-- Valuation
stock.item_valuation_rates      -- AVCO running average
stock.fifo_queue                -- FIFO cost layers per item/warehouse

-- Transactions (vouchers)
stock.stock_entries
stock.stock_entry_details
stock.purchase_receipts
stock.purchase_receipt_items
stock.delivery_notes
stock.delivery_note_items

-- Replenishment
stock.reorder_rules
stock.material_requests
stock.material_request_items

-- Quality
stock.quality_inspections
stock.quality_inspection_readings
```

### 3.3 Redis Caching Strategy

AwoERP uses Redis for hot-path stock balance reads to avoid expensive SLE aggregation on every page load:

| Cache Key Pattern | TTL | Invalidation |
|---|---|---|
| `stock:balance:{tenant}:{item}:{warehouse}` | 60s | On every SLE insert for that item/warehouse |
| `stock:valuation_rate:{tenant}:{item}` | 300s | On AVCO recalculation or SLE insert |
| `stock:reorder_alerts:{tenant}` | 600s | On SLE insert or reorder rule change |
| `stock:batch_expiry:{tenant}:{batch}` | 3600s | On batch update |

The Redis cache is a read-ahead optimisation only. PostgreSQL is the system of record. Cache misses fall through to the database transparently.

---

## 4. Initial Setup & Configuration

This section walks through the complete setup sequence for a new tenant enabling the Stock module. Follow the steps in order — dependencies exist between them.

### 4.1 Tenant-Level Stock Settings

Navigate to **Settings → Stock Settings** to configure the global stock behaviour for your tenant.

#### General Settings

| Setting | Default | Description |
|---|---|---|
| **Stock Module Enabled** | `false` | Master toggle for all stock functionality |
| **Allow Negative Stock** | `false` | Whether to permit stock balances to go below zero. Recommended: disabled for petroleum. |
| **Auto-Post to GL** | `true` | Automatically create GL journal entries on stock movements |
| **Default Valuation Method** | `AVCO` | Can be overridden per item. Options: FIFO, AVCO, Standard |
| **Fiscal Year Start** | From Finance settings | Determines period for stock aging and valuation reports |
| **Stock Adjustment Account** | Required | GL account for stock write-offs and adjustments |
| **Stock In Transit Account** | Optional | Intermediate GL account during inter-warehouse transfers |

#### Precision Settings

| Setting | Default | Description |
|---|---|---|
| **Qty Decimal Places** | `3` | Precision for quantity fields (e.g., litres at 3dp) |
| **Rate Decimal Places** | `4` | Precision for unit rates/costs |
| **Amount Decimal Places** | `2` | Precision for total amounts |

#### Transaction Settings

| Setting | Default | Description |
|---|---|---|
| **Require Batch for Tracked Items** | `true` | Enforce batch entry on receipts for batch-tracked items |
| **Allow Backdated Entries** | `false` | Enable posting stock entries to past dates (requires special role) |
| **Freeze Stock Entries Before** | `null` | Date before which no stock entries can be created or amended |
| **Default Purchase Warehouse** | Required | Goods received land here if not overridden on the transaction |
| **Default Sales Warehouse** | Required | Goods dispatched from here if not overridden on the transaction |

#### Configuration Example (JSON representation emitted by settings API)

```json
{
  "tenant_id": "anika-global",
  "stock_settings": {
    "module_enabled": true,
    "allow_negative_stock": false,
    "auto_post_gl": true,
    "default_valuation_method": "AVCO",
    "qty_precision": 3,
    "rate_precision": 4,
    "amount_precision": 2,
    "require_batch_for_tracked_items": true,
    "allow_backdated_entries": false,
    "freeze_entries_before": null,
    "default_purchase_warehouse_id": "wh-maanzoni-main",
    "default_sales_warehouse_id": "wh-maanzoni-main",
    "stock_adjustment_account": "5300-STOCK-ADJ",
    "stock_in_transit_account": "1215-STOCK-TRANSIT"
  }
}
```

---

### 4.2 Warehouse Hierarchy Setup

Warehouses in AwoERP are organised as a **tree using materialised paths**, enabling efficient subtree queries. You must set up at least one non-root warehouse before recording any stock.

#### Step 1 — Create the Root Warehouse (Company Node)

The root warehouse is a virtual aggregation node representing the company. It holds no physical stock.

```
Navigate: Stock → Warehouses → New Warehouse

Fields:
  Name:          Anika Global Limited
  Type:          Root
  Parent:        (none)
  Is Group:      Yes
  Company:       Anika Global Limited (from Finance)
  Currency:      KES
```

#### Step 2 — Create Site/Branch Warehouses

```
Navigate: Stock → Warehouses → New Warehouse

Fields:
  Name:          Shell Maanzoni Service Station
  Type:          Site
  Parent:        Anika Global Limited
  Is Group:      Yes
  Address:       Mombasa Road, Nairobi
  Currency:      KES
```

#### Step 3 — Create Operational Warehouses

These are the leaf warehouses where actual stock transactions occur.

```
Stores / Dry Goods:
  Name:        Maanzoni Dry Store
  Type:        Stores
  Parent:      Shell Maanzoni Service Station
  Is Group:    No
  Responsible: [Store Manager Role]

Wetstock (Petroleum):
  Name:        Maanzoni Forecourt Tanks
  Type:        Petroleum Tank Farm
  Parent:      Shell Maanzoni Service Station
  Is Group:    No
  Responsible: [Forecourt Manager Role]

Lubricants:
  Name:        Maanzoni Lubes Bay
  Type:        Stores
  Parent:      Shell Maanzoni Service Station
  Is Group:    No
```

#### Warehouse Tree Example

```
Anika Global Limited  (root, virtual)
└── Shell Maanzoni Service Station  (site, virtual)
    ├── Maanzoni Forecourt Tanks     (leaf — petrol/diesel tanks)
    ├── Maanzoni Lubes Bay           (leaf — lubricants)
    ├── Maanzoni Dry Store           (leaf — FMCG, consumables)
    └── Maanzoni QC Hold             (leaf — quarantine)
```

#### Warehouse Fields Reference

| Field | Type | Required | Description |
|---|---|---|---|
| `name` | VARCHAR(140) | ✓ | Display name |
| `warehouse_code` | VARCHAR(20) | Auto | Short code for reports |
| `type` | ENUM | ✓ | `root`, `site`, `stores`, `tank`, `transit`, `virtual`, `qc_hold` |
| `parent_id` | UUID | ✓ (non-root) | Parent warehouse |
| `path` | LTREE | Auto | Materialised path (e.g., `anika.maanzoni.drystore`) |
| `is_group` | BOOL | ✓ | If true, cannot have stock — aggregation only |
| `company_id` | UUID | ✓ | Associated company |
| `currency` | VARCHAR(3) | ✓ | Operating currency |
| `responsible_id` | UUID | Optional | User or role responsible |
| `disabled` | BOOL | `false` | Soft-delete; existing stock still visible |

---

### 4.3 Item Master Configuration

Items represent every SKU your business tracks. Getting the item master right is the most important setup step.

#### Item Types

| Type | GL Impact | Qty Tracked | Description |
|---|---|---|---|
| **Stock Item** | Inventory account | Yes | Physical goods (fuel, lubricants, FMCG) |
| **Non-Stock Item** | Expense account | No | Services, labour, one-time purchases |
| **Fixed Asset** | Asset account | No (separate module) | Equipment, vehicles |
| **Kit / Bundle** | Per component | Per component | A sellable pack that resolves to components |

#### Creating a Stock Item

```
Navigate: Stock → Item Master → New Item

--- Identity ---
Item Code:        PETROL-ULP-95       (unique per tenant; auto-generate option available)
Item Name:        Unleaded Petrol 95
Description:      Unleaded petrol, RON 95, for resale
Item Group:       Petroleum Products
UOM:              Litre

--- Stock Settings ---
Is Stock Item:    Yes
Valuation Method: AVCO               (override from tenant default)
Default Warehouse: Maanzoni Forecourt Tanks

--- Tracking ---
Has Batch No:     Yes
Has Serial No:    No
Track Expiry:     No                 (not applicable for petroleum)

--- Pricing ---
Standard Buy Rate:   KES 168.00/L    (reference only — actual from purchase)
Standard Sell Rate:  KES 185.50/L

--- Accounts ---
Asset Account:    1210-INVENTORY-PETROLEUM
Income Account:   4100-FUEL-SALES
COGS Account:     5100-FUEL-COGS

--- Reorder ---
Reorder Level:    15,000 L
Reorder Qty:      30,000 L
```

#### Item Fields Reference

| Field | Type | Description |
|---|---|---|
| `item_code` | VARCHAR(140) | Unique identifier |
| `item_name` | VARCHAR(280) | Full display name |
| `item_group_id` | UUID | Hierarchical category |
| `stock_uom` | VARCHAR(20) | Base unit of measure for stock |
| `is_stock_item` | BOOL | Whether to track quantities |
| `valuation_method` | ENUM | FIFO / AVCO / Standard |
| `standard_rate` | DECIMAL | Standard cost (used in Standard Costing) |
| `has_batch_no` | BOOL | Require batch on transactions |
| `has_serial_no` | BOOL | Require serial number |
| `track_expiry` | BOOL | Enable expiry date management |
| `min_order_qty` | DECIMAL | Minimum purchase quantity |
| `safety_stock` | DECIMAL | Buffer below reorder point |
| `lead_time_days` | INT | Supplier lead time for MRP |
| `asset_account` | VARCHAR | GL account for stock value |
| `income_account` | VARCHAR | GL account for sales |
| `cogs_account` | VARCHAR | GL account for COGS |

---

### 4.4 Unit of Measure (UOM) Setup

AwoERP supports multi-UOM items — you can purchase in one unit and sell in another. The system converts using defined factors.

#### Default UOMs (Pre-seeded)

The following UOMs are available out of the box:

```
Volume:    Litre (L), Kilolitre (KL), Gallon (GAL), Millilitre (ML)
Weight:    Kilogram (KG), Gram (G), Tonne (MT), Pound (LB)
Count:     Each (EA), Piece (PCS), Dozen (DOZ), Carton (CTN), Box (BOX)
Area:      Square Metre (SQM), Square Foot (SQFT)
Length:    Metre (M), Centimetre (CM), Foot (FT)
Time:      Hour (HR), Day (DAY)
```

#### Adding a UOM Conversion

```
Navigate: Stock → UOM → Conversions → New

Example: 1 Kilolitre = 1000 Litres
  From UOM:  KL
  To UOM:    L
  Factor:    1000

Example: 1 Carton = 12 Each
  From UOM:  CTN
  To UOM:    EA
  Factor:    12
```

#### Item-Specific UOM Conversions

If an item has non-standard pack sizes, define conversions at the item level:

```
Navigate: Item Master → [Item] → UOM Conversions tab

Item:      LUBRICANT-CASTROL-20W50-4L
  Purchase UOM:  CTN  (1 carton = 6 × 4L bottles)
  Stock UOM:     EA   (track as individual 4L cans)
  Sales UOM:     EA
  Factor:        6    (1 CTN = 6 EA)
```

---

### 4.5 Item Groups & Categories

Item groups form a hierarchy used for reporting, GL account defaulting, and access control.

#### Recommended Group Hierarchy (Petroleum Station)

```
All Items
├── Petroleum Products
│   ├── Fuels
│   │   ├── Petrol (ULP 93, ULP 95)
│   │   └── Diesel (AGO, HSD)
│   └── Gases (LPG)
├── Lubricants & Chemicals
│   ├── Engine Oils
│   ├── Transmission Fluids
│   └── Greases
├── Auto Accessories
│   ├── Tyres & Batteries
│   └── Wiper Blades
├── FMCG & Convenience
│   ├── Snacks & Beverages
│   └── Tobacco Products
└── Consumables (Internal)
    └── Cleaning & Maintenance
```

Each group can carry default GL accounts, which items inherit unless overridden.

---

### 4.6 Stock Valuation Methods

This is one of the most consequential configuration choices. Choose carefully — changing valuation methods mid-year requires a stock revaluation entry and Finance sign-off.

#### AVCO (Average Cost) — Recommended Default

The weighted average cost is recalculated every time stock is received. The formula:

```
New AVCO = (Current Stock Value + Incoming Stock Value)
           ÷ (Current Qty + Incoming Qty)
```

**Best for:** High-turnover items where lot-specific tracking is impractical — petrol, diesel, lubricants in bulk, commodities.

**Example:**

```
Opening:   1,000 L @ KES 160.00 = KES 160,000
Receipt 1: 5,000 L @ KES 168.00 = KES 840,000
           ─────────────────────────────────────
Total:     6,000 L               = KES 1,000,000
New AVCO:  KES 1,000,000 ÷ 6,000 = KES 166.67/L

Next sale of 500L posts COGS at KES 166.67 × 500 = KES 83,333.33
```

#### FIFO (First In, First Out)

Issues consume the oldest cost layers first. A queue of (qty, cost) layers is maintained per item/warehouse.

**Best for:** Perishables, pharmaceuticals, goods with significant price volatility, any item where regulators require FIFO (KRA audit preference).

**Example:**

```
Layer 1:  2,000 L @ KES 158.00  (oldest)
Layer 2:  3,000 L @ KES 165.00
Layer 3:  4,000 L @ KES 168.00  (newest)

Issue of 2,500 L:
  → Consumes Layer 1 fully:  2,000 L × 158.00 = 316,000
  → Consumes Layer 2 partial:  500 L × 165.00 =  82,500
  → COGS = KES 398,500
  → Remaining Layer 2: 2,500 L @ 165.00
```

#### Standard Costing

A predetermined standard cost is set per item. Variances between actual purchase price and standard cost are posted to a **Purchase Price Variance (PPV)** account.

**Best for:** Manufacturing environments with stable input costs; useful for budgeting and performance tracking. Less common in petroleum.

| Valuation Method | Pros | Cons | Recommended For |
|---|---|---|---|
| AVCO | Simple, stable COGS | Masks price movements | Fuel, bulk commodities |
| FIFO | Accurate, auditor-preferred | Complex queuing, more SLEs | Perishables, FMCG |
| Standard | Easy budgeting | Requires variance management | Manufacturing |

---

## 5. Warehouse Management

### 5.1 Warehouse Types

AwoERP defines the following warehouse types, each with specific behaviours:

| Type | Code | Stock Allowed | GL Posting | Notes |
|---|---|---|---|---|
| Root | `root` | No | No | Company-level aggregation node |
| Site | `site` | No | No | Branch/location aggregation |
| Stores | `stores` | Yes | Yes | General goods storage |
| Petroleum Tank | `tank` | Yes | Yes | Wetstock — dip-calibrated |
| Transit | `transit` | Yes | Optional | In-transit between locations |
| QC Hold | `qc_hold` | Yes | Optional | Quarantine pending inspection |
| Virtual | `virtual` | Yes | No | Theoretical/planning nodes |
| Scrap | `scrap` | Yes | Yes | Write-off destination |

### 5.2 Storage Locations & Bins

Within a leaf warehouse, you can optionally track stock at **bin** level (shelf, rack, tank compartment, bay position).

#### Enabling Bin-Level Tracking

```
Navigate: Settings → Stock Settings → Bin Tracking

Enable Bin Tracking:    Yes
Mandatory on Receipts:  No   (guide but don't enforce)
Mandatory on Issues:    No
```

#### Creating Bins

```
Navigate: Stock → Warehouses → [Warehouse] → Bins tab → Add Bin

Warehouse:      Maanzoni Dry Store
Bin Code:       A-01-01         (Row A, Shelf 01, Position 01)
Bin Name:       Aisle A, Shelf 1, Slot 1
Capacity:       200 KG
Item Restriction: (optional — restrict bin to specific items)
```

#### Bin Transfer

Moving stock within a warehouse between bins does not create a warehouse-level stock movement — only a bin quantity update. No GL entry is generated.

```
Navigate: Stock → Bin Transfer → New

From Warehouse:  Maanzoni Dry Store
From Bin:        A-01-01
To Bin:          B-02-03
Item:            LUBRICANT-CASTROL-20W50-4L
Qty:             24 EA
Reason:          Racking reorganisation
```

### 5.3 Warehouse Permissions

Access to warehouses is controlled through the role-permission system. A user can be granted Read, Write, or No Access per warehouse. This is enforced at the query layer — users only see stock balances for warehouses they have access to.

```
Navigate: Settings → Roles → [Role] → Warehouse Access tab

Role:              Forecourt Attendant
Warehouse Access:
  Maanzoni Forecourt Tanks    → Read
  Maanzoni Lubes Bay          → Read
  Maanzoni Dry Store          → No Access
  Maanzoni QC Hold            → No Access
```

### 5.4 Virtual Warehouses

Virtual warehouses are used for planning and reporting scenarios where no physical stock movement occurs. Common uses:

- **In-Transit:** Represents stock dispatched from supplier but not yet received. The purchase order creates a receipt into In-Transit; a physical arrival entry moves it to the actual store.
- **Customer Consignment:** Stock held at a customer site on consignment — technically still your asset.
- **Reserved Stock:** Soft-reserves stock against a sales order without removing it from the warehouse balance.

---

## 6. Item Master & Variants

### 6.1 Item Attributes & Variants

For items that come in multiple sizes, colours, or specifications, use Item Variants to avoid creating separate item codes for each combination.

#### Defining Attributes

```
Navigate: Stock → Item Attributes → New

Attribute:    Size
Values:       1L, 4L, 20L, 200L

Attribute:    Grade
Values:       10W40, 15W40, 20W50, 5W30
```

#### Creating a Template Item with Variants

```
Navigate: Stock → Item Master → New Item

Item Code:     CASTROL-ENGINE-OIL
Is Template:   Yes
Attributes:
  - Grade:     (All values)
  - Size:      1L, 4L, 20L
```

AwoERP auto-generates variant combinations:
```
CASTROL-ENGINE-OIL-10W40-1L
CASTROL-ENGINE-OIL-10W40-4L
CASTROL-ENGINE-OIL-10W40-20L
CASTROL-ENGINE-OIL-15W40-1L
... (and so on)
```

Each variant is a full stock item with independent stock ledger, valuation, and reorder rules.

### 6.2 Item Pricing

AwoERP supports multi-price-list pricing, customer/supplier-specific pricing, and quantity-break pricing.

#### Price Lists

```
Navigate: Stock → Price Lists

System Price Lists (auto-created):
  Standard Buying    — default purchase prices
  Standard Selling   — default sale prices

Custom Price Lists:
  Shell Fleet Pricing — discounted rates for fleet accounts
  Wholesale Pricing   — volume-based rates
```

#### Setting Item Prices

```
Navigate: Stock → Item Prices → New

Item:         PETROL-ULP-95
Price List:   Standard Selling
UOM:          Litre
Rate:         KES 185.50
Min Qty:      0
Valid From:   2026-01-01
Valid To:     (blank = always valid)
```

#### Pricing Priority (Cascade)

When a sales transaction is created, the system applies the first matching price rule:

1. Customer-specific price list
2. Customer group price list
3. Campaign/promotional price (date-bounded)
4. Quantity-break pricing
5. Standard selling price list
6. Item standard rate (fallback)

### 6.3 Reorder Rules

Reorder rules define the automatic replenishment trigger for each item/warehouse combination.

```
Navigate: Stock → Reorder Rules → New

Item:                 PETROL-ULP-95
Warehouse:            Maanzoni Forecourt Tanks
Reorder Point:        15,000 L     (trigger when stock falls below this)
Reorder Qty:          30,000 L     (quantity to request)
Min Order Qty:        20,000 L     (cannot order less than this)
Max Stock Level:      80,000 L     (system won't suggest ordering above this)
Preferred Supplier:   Vivo Energy Kenya Limited
Lead Time Days:       1
Replenishment Method: Create Material Request
```

The `stock_replenishment_check` job runs hourly and evaluates all reorder rules. When stock falls at or below the reorder point, it creates a **Material Request** (draft purchase requisition) and optionally sends an alert to the purchasing team.

### 6.4 Item Bundling & Kitting

Kits allow you to sell or transfer a predefined combination of items as a single line item, with the system automatically exploding components.

```
Navigate: Stock → Items → New Item

Item Code:      CAR-CARE-KIT-001
Item Name:      Basic Car Care Kit
Is Kit:         Yes
Kit Components:
  - CASTROL-ENGINE-OIL-10W40-1L   × 1
  - OIL-FILTER-GENERIC             × 1
  - WINDSCREEN-WASHER-500ML        × 1
  - AIR-FRESHENER-VANILLA          × 1

Sell at Kit Price:  KES 2,500.00
```

When the kit is added to a delivery note, the system explodes it into individual component lines, each deducted from stock independently.

---

## 7. Stock Transactions

### 7.1 Stock Entry Types

AwoERP uses a unified **Stock Entry** document for internal movements, with `purpose` controlling the entry type. External-facing entries (purchases, sales) use their own document types but write the same SLEs.

| Entry Type | Source Warehouse | Target Warehouse | Qty Sign | GL Impact |
|---|---|---|---|---|
| Material Receipt | — | Destination | +QTY | Dr Inventory, Cr Purchase Clearing |
| Material Issue | Source | — | -QTY | Dr COGS/Expense, Cr Inventory |
| Material Transfer | Source | Destination | ±QTY | Within same inventory account |
| Manufacture | RM Warehouse | FG Warehouse | ±QTY | Complex BOM costing |
| Repack | Source | Destination | ±QTY | UOM change, same value |
| Send to Sub-Contractor | Source | — | -QTY | Dr WIP, Cr Inventory |
| Receive from Sub-Contractor | — | Destination | +QTY | Dr Inventory, Cr WIP |
| Stock Reconciliation | — | Destination | ±QTY | Adjustment account |

### 7.2 Purchase Receipts

Purchase Receipts record the physical arrival of goods against a Purchase Order.

#### Workflow

```
Purchase Order (Purchasing module)
        │
        ▼ (goods physically arrive)
Purchase Receipt (Stock module)
  ├── Validates PO reference and open quantities
  ├── Creates SLEs in target warehouse
  ├── Updates AVCO or FIFO queue
  ├── Creates GL: Dr Inventory / Cr Goods Receipt Note (GRN) Clearing
  └── Updates PO received quantity
        │
        ▼ (Finance matches GRN to Supplier Invoice)
Purchase Invoice (Finance module)
  └── Clears GRN, posts Cr Accounts Payable
```

#### Creating a Purchase Receipt

```
Navigate: Stock → Purchase Receipts → New
(Or: From Purchase Order → Create → Purchase Receipt)

Supplier:           Vivo Energy Kenya Limited
PO Reference:       PO-2026-00147
Receipt Date:       2026-06-22
Receive To Warehouse: Maanzoni Forecourt Tanks

Items Tab:
  Item Code:        PETROL-ULP-95
  Accepted Qty:     30,000 L
  Rejected Qty:     0
  Rate:             KES 168.00
  Amount:           KES 5,040,000.00
  Batch No:         VIVO-2026-06-A

  Item Code:        DIESEL-AGO
  Accepted Qty:     20,000 L
  Rejected Qty:     0
  Rate:             KES 155.00
  Amount:           KES 3,100,000.00
  Batch No:         VIVO-2026-06-B
```

#### Quality Inspection on Receipt

If QC is enabled for an item, the system blocks submission until a Quality Inspection is created and accepted:

```
Items Tab → [Row] → Inspection Required: Yes
        │
        ▼
Quality Inspection created automatically
  Readings entered by QC team
  QC Result: Accepted / Rejected
        │
        ▼ (if accepted)
Purchase Receipt can be submitted
```

### 7.3 Delivery Notes

Delivery Notes record the physical dispatch of goods against a Sales Order or direct sale.

```
Navigate: Stock → Delivery Notes → New

Customer:          Fleet Customer — ABC Logistics Ltd
Delivery Date:     2026-06-22
Dispatch From:     Maanzoni Forecourt Tanks

Items Tab:
  Item Code:       DIESEL-AGO
  Qty:             500 L
  UOM:             Litre
  Rate:            KES 160.00
  Amount:          KES 80,000.00
  Batch No:        VIVO-2026-06-B       (system auto-suggests FEFO batch)

GL Impact (auto-posted):
  Dr  5100-DIESEL-COGS       KES 77,500.00  (at AVCO rate of 155.00)
  Cr  1211-INVENTORY-DIESEL  KES 77,500.00
```

### 7.4 Stock Transfers

Inter-warehouse transfers move stock between two warehouses within the same company.

#### One-Step Transfer (Immediate)

```
Navigate: Stock → Stock Entry → New

Purpose:          Material Transfer
From Warehouse:   Maanzoni Forecourt Tanks
To Warehouse:     Maanzoni Lubes Bay
Date:             2026-06-22

Items:
  Item:       SHELL-HELIX-10W40-4L   Qty: 24   UOM: EA
```

A one-step transfer immediately reduces source stock and increases destination stock in the same transaction.

#### Two-Step Transfer (In-Transit)

For transfers between physical sites where goods are in a vehicle for a period:

**Step 1 — Outward Entry:**
```
Purpose:          Material Transfer (Outward)
From Warehouse:   Site A - Stores
To Warehouse:     In-Transit Warehouse
```

**Step 2 — Inward Entry (at destination):**
```
Purpose:          Material Transfer (Inward)
From Warehouse:   In-Transit Warehouse
To Warehouse:     Site B - Stores
```

Only after Step 2 does the destination site show the stock.

### 7.5 Material Issues & Returns

#### Internal Issue (for consumption)

```
Navigate: Stock → Stock Entry → New

Purpose:          Material Issue
From Warehouse:   Maanzoni Dry Store
Cost Centre:      Shell Maanzoni - Forecourt Operations

Items:
  Item:    CLEANING-CLOTHS-10PK    Qty: 5   EA
  Item:    PUMP-LUBRICANT-500ML    Qty: 2   EA
```

#### Customer Return / Sales Return

```
Navigate: Stock → Delivery Notes → [Original DN] → Create → Return

Return to Warehouse:   Maanzoni Dry Store
Return Reason:         Wrong grade delivered
Return Qty:            10 L     (partial return allowed)
```

A return entry creates **positive** SLEs (stock comes back in) and reverses GL entries. The returned stock inherits the original batch.

### 7.6 Manufacturing Entries

When manufacturing is enabled, stock entries handle raw material consumption and finished goods production.

```
Work Order → [Manufacture trigger]

Auto-creates two entries:
  1. Material Issue (RM Consumption)
     From: Raw Material Warehouse
     Items: All BOM components × manufactured qty

  2. Material Receipt (FG Production)
     To: Finished Goods Warehouse
     Item: Finished product at BOM rolled-up cost
```

---

## 8. Inventory Valuation & Costing

### 8.1 FIFO Valuation

FIFO maintains a **cost layer queue** per item/warehouse. Layers are created on every receipt and consumed in order on every issue.

#### FIFO Queue Structure

```sql
-- stock.fifo_queue
fifo_id          UUID
tenant_id        UUID
item_code        VARCHAR
warehouse_id     UUID
incoming_rate    DECIMAL(18,6)
qty_remaining    DECIMAL(18,6)
posting_datetime TIMESTAMPTZ
voucher_no       VARCHAR
```

#### FIFO Consumption Algorithm (Go pseudocode)

```go
func ConsumeFIFO(ctx context.Context, item string, warehouseID uuid.UUID, qtyToConsume float64) ([]FIFOLayer, float64, error) {
    layers, err := repo.GetFIFOLayers(ctx, item, warehouseID) // ordered by posting_datetime ASC
    if err != nil { return nil, 0, err }

    var consumed []FIFOLayer
    var totalCost float64
    remaining := qtyToConsume

    for _, layer := range layers {
        if remaining <= 0 { break }
        take := math.Min(remaining, layer.QtyRemaining)
        consumed = append(consumed, FIFOLayer{Rate: layer.IncomingRate, Qty: take})
        totalCost += take * layer.IncomingRate
        remaining -= take
        layer.QtyRemaining -= take
        repo.UpdateFIFOLayer(ctx, layer)
    }

    if remaining > 0 && !settings.AllowNegativeStock {
        return nil, 0, ErrInsufficientStock
    }

    return consumed, totalCost / qtyToConsume, nil // returns weighted rate of consumed layers
}
```

### 8.2 AVCO (Weighted Average Cost)

AVCO maintains a single **running average rate** per item/warehouse, recalculated on every inbound movement.

#### AVCO Recalculation

```go
func RecalculateAVCO(ctx context.Context, item string, warehouseID uuid.UUID, incomingQty, incomingRate float64) (float64, error) {
    balance, err := repo.GetStockBalance(ctx, item, warehouseID)
    if err != nil { return 0, err }

    currentValue := balance.Qty * balance.ValuationRate
    incomingValue := incomingQty * incomingRate
    newQty := balance.Qty + incomingQty

    if newQty == 0 { return 0, nil }

    newRate := (currentValue + incomingValue) / newQty
    return newRate, nil
}
```

**Important:** AVCO can produce unexpected results when stock goes negative (allowed in some configurations). A negative balance multiplied by a positive AVCO rate creates a negative inventory value. Enabling negative stock with AVCO requires Finance sign-off.

### 8.3 Standard Costing

Standard costing posts all receipts at the **standard rate** defined on the item, regardless of actual purchase price. The variance is posted to a PPV account.

```
Standard Rate:     KES 165.00/L
Actual Purchase:   KES 168.00/L
Qty Received:      10,000 L

GL Entries:
  Dr  1210-INVENTORY         KES 1,650,000   (at standard rate)
  Dr  5310-PPV-ADVERSE        KES    30,000   (adverse variance)
  Cr  2010-GRN-CLEARING      KES 1,680,000   (at actual rate)
```

### 8.4 Landed Costs

Landed cost vouchers distribute freight, customs, clearing, and insurance charges across the items in a purchase receipt, adjusting their valuation rates.

```
Navigate: Stock → Landed Cost Vouchers → New

Purchase Receipt:  GRN-2026-00089
Applicable Charges:
  Freight (Mombasa to Nairobi):   KES  45,000
  Customs Duty:                   KES 180,000
  Port Charges:                   KES  12,000
  Insurance:                      KES   8,000
  ─────────────────────────────────────────────
  Total Landed Costs:             KES 245,000

Distribution Method:  By Amount   (alternatives: By Qty, By Weight)

Allocation:
  PETROL-ULP-95  (KES 5,040,000 = 62.5%)   → KES 153,125
  DIESEL-AGO     (KES 3,100,000 = 37.5%)   → KES  91,875
```

After posting the landed cost voucher, the valuation rate of each item is adjusted upward. All subsequent COGS postings use the landed rate.

### 8.5 Cost Centre Allocation

When stock is issued for internal consumption, the COGS or expense posting must be attributed to a cost centre. AwoERP supports:

- Direct allocation to a single cost centre
- Percentage-based split across multiple cost centres
- Auto-allocation based on the requesting department

---

## 9. Batch & Serial Number Tracking

### 9.1 Batch Management

Batches group units of the same item from the same production run or supplier lot. Every movement of a batch-tracked item must reference a batch.

#### Creating a Batch

Batches can be created manually or automatically from a Purchase Receipt or Manufacturing entry.

```
Navigate: Stock → Batches → New

Item:              PETROL-ULP-95
Batch ID:          VIVO-2026-06-A        (auto or manual)
Manufacturing Date: 2026-05-30
Expiry Date:       2028-05-29           (24-month shelf life)
Supplier Lot No:   LOT-VEK-240501
Qty on Creation:   30,000 L
Reference:         GRN-2026-00089
Notes:             RON 95 tested and certified
```

#### Batch Traceability

AwoERP provides both forwards and backwards trace from any batch:

**Backwards Trace** — Where did this batch come from?
```
Navigate: Stock → Batch → [VIVO-2026-06-A] → Trace tab

Origin:
  Purchase Receipt   GRN-2026-00089     30,000 L    2026-06-22
  Supplier:          Vivo Energy Kenya
  Supplier Lot:      LOT-VEK-240501
```

**Forwards Trace** — Where did units of this batch go?
```
Movements:
  Delivery Note   DN-2026-01205    Customer: Fleet XYZ       500 L   2026-06-22
  Delivery Note   DN-2026-01210    Customer: Walk-in          10 L   2026-06-22
  Delivery Note   DN-2026-01215    Customer: Cash             20 L   2026-06-22
  Remaining:                                                29,470 L
```

#### Batch Recall Procedure

In the event of a product recall:

1. Navigate to **Stock → Batch → [Batch ID] → Actions → Place on Hold**
2. System blocks any new outgoing transactions for this batch
3. **Batch → Forwards Trace** identifies all customers who received units
4. Create a **Stock Entry (Return to Supplier)** to send back unsold stock
5. Generate the **Batch Recall Report** for regulatory submission

### 9.2 Serial Number Tracking

Serial numbers track individual units. Each serial number has one current location and a full movement history.

```
Navigate: Stock → Serial Numbers → New

Item:             LPG-CYLINDER-12KG
Serial No:        CYL-KE-2024-001847
Purchase Date:    2024-03-15
Warranty Expiry:  2027-03-14
Current Location: Maanzoni Dry Store
Status:           Active
```

Every transaction involving the item requires specifying the serial number. The system enforces:
- A serial number cannot appear twice in the same stock movement
- A serial number in status `Delivered` cannot be received unless it's a return
- Warranties are tracked and expiry alerts generated automatically

### 9.3 Expiry Date Management

For perishable goods and pharmaceutical items:

```
Navigate: Settings → Stock Settings → Expiry Management

Notify Before Expiry (Days):  90, 30, 7   (multiple alerts)
Auto-Block Expired Batches:   Yes
FEFO Enforced:                Yes          (issue oldest expiry first)
Allow Override with Reason:   Yes          (manager approval required)
```

#### FEFO (First Expired, First Out)

When FEFO is enabled, the system automatically sorts suggested batches by expiry date (soonest-to-expire first) when creating delivery notes and stock issues. This is advisory unless FEFO enforcement is turned on, in which case the system rejects attempts to use a newer batch when an older batch has stock.

### 9.4 Batch-Level Costing

Even with AVCO as the item-level valuation method, individual batches retain their actual received cost. This supports:

- **Batch-level profitability analysis** — actual margin per batch
- **Recall cost calculation** — exact value of recalled stock
- **Return pricing** — credit customer/supplier at original cost

---

## 10. Stock Reconciliation & Physical Count

### 10.1 Cycle Counting

Cycle counting is a continuous audit process where different subsets of inventory are counted on a rotating schedule without halting operations.

#### Setting Up Cycle Count Schedule

```
Navigate: Stock → Cycle Count Settings

Schedule:
  A-Items (top 20% by value):   Count Weekly
  B-Items (middle 30% by value): Count Monthly
  C-Items (bottom 50% by value): Count Quarterly

Auto-Assign Counters:  Yes
Print Count Sheets:    Yes
```

#### Executing a Cycle Count

```
Navigate: Stock → Physical Inventory → New Cycle Count

Warehouse:     Maanzoni Dry Store
Count Date:    2026-06-22
Items:         (auto-populated based on schedule; manual items can be added)
Assigned To:   Store Supervisor

Count Sheet columns:
  Item Code | Item Name | System Qty | Counted Qty | Variance | Variance %
```

The counter enters physical counts. System quantity is hidden by default (blind count mode) to prevent anchoring bias. Enable **Show System Qty** if your process requires it.

### 10.2 Full Physical Inventory

A full count stops stock movements for all items in the warehouse during the count period.

```
Navigate: Stock → Physical Inventory → New Full Count

Warehouse:          Maanzoni Dry Store
Count Start Date:   2026-06-30 08:00
Count End Date:     2026-06-30 18:00
Freeze Transactions: Yes    (no stock entries allowed during count window)

Pre-Count Actions (system auto-completes):
  ✓ Print count sheets for all active items
  ✓ Post all pending stock entries
  ✓ Clear staging areas (ensure no in-process stock)
```

### 10.3 Variance Analysis

After counts are entered, the system produces a variance report:

```
Stock Reconciliation Variance Report — Maanzoni Dry Store — 2026-06-30

Item                    System Qty  Count Qty  Variance  Value (KES)  Variance %
─────────────────────────────────────────────────────────────────────────────────
CASTROL-HELIX-10W40-4L    120 EA      116 EA    -4 EA     -3,600.00   -3.33%
SHELL-ROTELLA-20W50-4L     48 EA       48 EA     0 EA          0.00    0.00%
WIPER-BLADE-21IN           30 EA       31 EA    +1 EA       +450.00   +3.33%
AIR-FRESHENER-VANILLA      80 EA       75 EA    -5 EA       -625.00   -6.25%
─────────────────────────────────────────────────────────────────────────────────
Total Negative Variance:                                  -4,225.00
Total Positive Variance:                                    +450.00
Net Variance:                                             -3,775.00
```

Any variance above the configured threshold (e.g., ±2% of value) triggers an alert to the Governance module and requires manager-level approval before the reconciliation can be posted.

### 10.4 Reconciliation Approval Workflows

```
Count Sheet Submitted
        │
        ▼
Variance within threshold?
  YES → Auto-approve, post Stock Reconciliation Entry
  NO  → Requires manager review
          │
          ├── Manager: Approve → post entry, log reason
          ├── Manager: Investigate → create investigation task, hold posting
          └── Manager: Reject → discard count, schedule recount
```

The approval workflow uses Temporal for state persistence, ensuring the process survives server restarts and can time out (escalate) after a configurable period.

---

## 11. Demand Planning & Replenishment

### 11.1 Reorder Point Planning

The simplest replenishment strategy: when stock falls at or below the reorder point, create a purchase request.

```
Reorder Point  = Safety Stock + (Average Daily Usage × Lead Time Days)

Example:
  Average Daily Usage of PETROL-ULP-95:  8,000 L/day
  Lead Time Days:                         1 day
  Safety Stock:                           5,000 L
  ────────────────────────────────────────────────────
  Reorder Point = 5,000 + (8,000 × 1) = 13,000 L
```

### 11.2 Min/Max Replenishment

Min/Max defines a floor (min) and ceiling (max) stock level. The system orders up to Max whenever stock drops below Min.

```
Navigate: Stock → Reorder Rules → [Item/Warehouse] → Method: Min/Max

Min Stock:    15,000 L    (reorder trigger)
Max Stock:    80,000 L    (order up to this level)
Order Qty:    80,000 - current_stock  (variable)
```

### 11.3 Material Requirement Planning (MRP)

MRP works backward from sales orders and demand forecasts to determine what raw materials need to be purchased and when.

```
Navigate: Stock → MRP → Run MRP

Planning Horizon:  30 days
Consider:
  ✓ Sales Orders (confirmed demand)
  ✓ Demand Forecast (statistical forecast)
  ✓ Safety Stock
  ✓ Existing Purchase Orders (supply in pipeline)
  ✓ Open Production Orders

Output:
  → Planned Purchase Orders (draft, needs buyer review)
  → Planned Production Orders (draft, needs planner review)
```

MRP integrates with the Finance module to respect budget limits — planned orders that would exceed department budgets are flagged for approval rather than auto-submitted.

### 11.4 Automated Purchase Requests

When the replenishment engine creates a Material Request, the Purchasing module converts it to a Purchase Order through the following flow:

```
Replenishment Engine → Material Request (auto-draft)
                              │
                     Purchasing Team Reviews
                              │
                    Convert to Request for Quote
                              │
                    Supplier Quotes Received
                              │
                    Purchase Order Created
                              │
                    Goods Received (Purchase Receipt)
                              │
                    Stock Updated
```

---

## 12. Multi-Location & Multi-Currency

### 12.1 Multi-Warehouse Operations

AwoERP supports consolidated reporting across all warehouses in a tenant. From the Stock Dashboard:

```
Stock Dashboard — Anika Global Limited
───────────────────────────────────────────────────────
Warehouse                    Total Stock Value (KES)
───────────────────────────────────────────────────────
Shell Maanzoni Service Station
  ├── Maanzoni Forecourt Tanks       48,920,000
  ├── Maanzoni Lubes Bay              2,345,000
  └── Maanzoni Dry Store               890,000
                                  ──────────────
Maanzoni Total:                    52,155,000

[Future Site - Thika Road]
  ├── Thika Forecourt Tanks              0
  └── Thika Dry Store                    0
                                  ──────────────
Thika Total:                              0
───────────────────────────────────────────────────────
Grand Total:                       52,155,000
```

### 12.2 Inter-Company Transfers

Where a tenant operates multiple legal entities (e.g., Anika Global Limited operating both a Nairobi and Mombasa service station as separate companies), the Inter-Company Transfer feature handles the accounting on both sides.

```
Navigate: Stock → Inter-Company Transfer → New

From Company:    Anika Global Nairobi Ltd
From Warehouse:  Nairobi Central Stores
To Company:      Anika Global Mombasa Ltd
To Warehouse:    Mombasa Branch Stores

Transfer Price Method:
  - At Cost (inter-company at cost)
  - At Selling Price + x% markup
  - Custom rate
```

The system auto-creates:
- A Delivery Note in the selling company
- A Purchase Receipt in the buying company
- Inter-company receivable/payable entries in both GL ledgers

### 12.3 Foreign Currency Inventory

For items purchased in USD (e.g., imported lubricant brands), AwoERP maintains both the foreign currency cost and the KES equivalent.

```
Purchase Receipt:
  Item:           CASTROL-EDGE-5W30-1L
  Supplier:       Castrol International (USD)
  Foreign Rate:   USD 8.50/unit
  Exchange Rate:  KES 129.50/USD (from Finance rate table)
  KES Rate:       KES 1,100.75/unit

Valuation:
  Inventory posted in KES at rate on posting date
  Forex revaluation at period-end adjusts for rate movements
```

---

## 13. Quality Control Integration

### 13.1 Inspection on Receipt

Quality inspections can be triggered automatically when goods are received, based on:

- Item setting: `Inspection Required on Purchase: Yes`
- Item Group setting: all items in group require inspection
- Supplier setting: all receipts from this supplier require inspection

```
Navigate: Stock → Quality Inspections → [Auto-created from GRN]

Purchase Receipt:   GRN-2026-00089
Item:               PETROL-ULP-95
Batch:              VIVO-2026-06-A
Sample Size:        100 ML
Inspector:          Lab Technician

Readings:
  Parameter         Reading    Min     Max    Status
  RON               95.2       94.0    97.0   ✓ PASS
  Density @ 15°C    0.751      0.720   0.760  ✓ PASS
  Colour            Clear      Clear   —      ✓ PASS
  Water Content     0.01%      0       0.05%  ✓ PASS

Overall Result:     ACCEPTED
```

If the result is **REJECTED**, the received quantity is moved to the **QC Hold Warehouse** and a non-conformance report is raised. The Purchase Receipt cannot be fully accepted until QC is resolved.

### 13.2 Quarantine Warehouses

The QC Hold warehouse type isolates suspect stock:

- Transactions **to** QC Hold: automatically created on QC rejection
- Transactions **from** QC Hold: require QC Manager role
- QC Hold stock does not appear in available-to-promise calculations
- Stock in QC Hold is still valued on the balance sheet (separate line)

### 13.3 Quality Hold & Release

```
Navigate: Stock → Quality Hold → [Record] → Actions

Actions available:
  Release to Stores       — inspection passed, transfer to normal warehouse
  Return to Supplier      — raise Return to Vendor (RTV) entry
  Rework / Reprocess      — route to Manufacturing for correction
  Write Off               — if scrapped, create Stock Entry (Scrap)
  Accept Under Waiver     — override with documented justification (QC Manager only)
```

---

## 14. Petroleum & Forecourt Stock (Industry Extension)

This section covers AwoERP's purpose-built extensions for petroleum retail operations — a critical feature for Shell Maanzoni and similar service station operators.

### 14.1 Wetstock Management

Wetstock refers to liquid fuel held in underground storage tanks (USTs). Unlike packaged goods, wetstock quantities are subject to:

- **Thermal expansion/contraction** (temperature-compensated volumes)
- **Evaporation losses**
- **Meter vs. tank measurement discrepancies**
- **Regulatory dip tolerance limits** (EPRA guidelines)

AwoERP's Wetstock module maintains a parallel measurement layer on top of standard stock, reconciling meter readings, dip readings, and delivery records every shift.

### 14.2 Tank Dip Readings

Dip readings are the physical measurement of fuel depth in a tank, converted to volume using a calibration chart.

#### Tank Configuration

```
Navigate: Stock → Petroleum → Tanks → New

Tank ID:           T-001
Tank Name:         Petrol Tank 1 (ULP 95)
Product:           PETROL-ULP-95
Warehouse:         Maanzoni Forecourt Tanks
Capacity (L):      60,000
Dip Calibration:   [Upload calibration chart CSV]
  (Maps depth in mm to volume in litres based on tank geometry)
Safe Fill Level:   57,000 L   (95% of capacity)
Low Level Alert:   15,000 L
Dead Stock Level:   2,000 L   (cannot be pumped, always remains)
```

#### Recording a Dip Reading

```
Navigate: Stock → Petroleum → Dip Readings → New

Tank:              T-001 (Petrol ULP 95)
Shift:             [Links to current open shift]
Reading Date/Time: 2026-06-22 06:00
Dip MM:            1,847 mm
Calculated Volume: 28,450 L   (from calibration chart lookup)
Temperature °C:    24.5       (for temperature compensation if required)
Observed By:       [Attendant name]
Witnessed By:      [Supervisor name]

Water Dip (mm):    12         (any water accumulation at tank bottom)
Water Volume (L):  85         (to be excluded from product volume)
Net Product Volume: 28,365 L
```

### 14.3 Meter-Based Dispensing

Each pump has one or more meters that record cumulative volume dispensed.

#### Meter Types Tracked

| Meter Type | Description | Primary Use |
|---|---|---|
| Electronic Volume Meter | Precision litres dispensed by pump controller | COGS and SLE qty |
| Electronic Cash Meter | Cumulative KES value at pump price | Revenue reconciliation |
| Manual Mechanical Meter | Backup dials readable without power | Cross-validation and fraud detection |

#### Shift Meter Readings

```
Navigate: Stock → Petroleum → Meter Readings → New

Shift:             Morning Shift 2026-06-22 (06:00–14:00)
Pump:              Pump 3 (Diesel)

Opening Readings (from previous shift close):
  Electronic Volume:    1,284,532.450 L
  Electronic Cash:      KES 199,102,029.75
  Manual Mechanical:    1,284,530 L

Closing Readings:
  Electronic Volume:    1,285,682.450 L
  Electronic Cash:      KES 199,280,229.75
  Manual Mechanical:    1,285,680 L

Computed Dispensed:
  Electronic Volume:    1,150.000 L
  Electronic Cash:      KES 178,200.00
  Manual Mechanical:    1,150 L

Cross-Validation:
  Volume vs Manual:     0.000 L variance   ✓ PASS
  Volume vs Cash:       KES 178,200 ÷ 155.00 = 1,149.677 L   Δ = 0.323 L  ✓ within tolerance
```

**Cross-validation rules:**

- Electronic Volume vs Manual Mechanical: tolerance ≤ 2 L per shift
- Electronic Volume vs Electronic Cash (back-calculated at pump price): tolerance ≤ 5 L per shift
- Any variance exceeding tolerance triggers an automatic alert and flags the shift for investigation

### 14.4 Variance & Loss Reconciliation

At shift close (or daily close), the wetstock reconciliation engine computes:

```
Wetstock Reconciliation — T-001 Petrol ULP 95 — 2026-06-22 Morning Shift
──────────────────────────────────────────────────────────────────────────

Opening Stock (Dip):              28,365 L   (06:00 dip reading)
Add: Deliveries Received:         30,000 L   (GRN-2026-00089)
Less: Meter Sales:                 1,150 L   (sum of pump meters for this tank)
Theoretical Closing Stock:        57,215 L

Actual Closing Stock (Dip):       57,150 L   (14:00 dip reading)

Wetstock Variance:                   -65 L   (LOSS)
Variance as % of Throughput:       -5.65%   ⚠ ABOVE EPRA THRESHOLD (0.5%)

Variance Breakdown:
  Allowable Evaporation/Handling:   -12 L   (0.1% of opening + receipts)
  Unexplained Variance:             -53 L   ⚠ REQUIRES INVESTIGATION

GL Posting:
  Dr  5400-WETSTOCK-LOSSES   KES 8,904.00   (53 L × KES 168.00 AVCO)
  Cr  1211-INVENTORY-PETROL  KES 8,904.00
```

Variance reasons can be coded:
- `EV` — Evaporation (within allowance)
- `TH` — Thermal contraction
- `DM` — Dip measurement error
- `SP` — Spillage (documented incident)
- `TF` — Theft (triggers security investigation)
- `UK` — Unknown (triggers mandatory review)

### 14.5 Environmental Compliance Tracking

EPRA (Energy and Petroleum Regulatory Authority) and NEMA (National Environment Management Authority) requirements:

```
Navigate: Stock → Petroleum → Compliance Tracking

Tracked Parameters:
  Daily tank ullage reports            (EPRA requirement)
  Monthly water accumulation readings
  Annual tank integrity tests
  Vapour recovery unit efficiency
  Spill incident log (NEMA Form ENV-08)
  Environmental bond status

Reports Available:
  EPRA Monthly Stock Return
  NEMA Environmental Compliance Summary
  Tank Integrity Certificate tracker
```

---

## 15. Stock Reports & Analytics

### 15.1 Standard Reports

| Report | Description | Filters Available |
|---|---|---|
| **Stock Balance** | Current qty and value per item/warehouse | Warehouse, Item Group, Date, Below Reorder |
| **Stock Ledger** | All SLEs in a period | Item, Warehouse, Date Range, Voucher Type |
| **Stock Ageing** | How long stock has been held | Warehouse, Item Group, Ageing Brackets |
| **Item Price List** | All prices across price lists | Price List, Item Group, Supplier/Customer |
| **Batch-Wise Balance** | Stock by batch with expiry | Item, Warehouse, Expiry Date Range |
| **Serial Number Status** | Current location of serialised items | Item, Status, Warranty Status |
| **Purchase Receipt Register** | All GRNs in a period | Supplier, Warehouse, Date Range |
| **Delivery Note Register** | All DNs in a period | Customer, Warehouse, Date Range |
| **Stock Reorder Report** | Items below reorder point | Warehouse, Item Group |
| **Wetstock Reconciliation** | Petroleum shift/daily reconciliation | Site, Tank, Date Range |
| **Variance Analysis** | Count vs system variances | Warehouse, Item, Date |
| **COGS Report** | Cost of goods sold by period | GL Period, Item Group |
| **Inventory Turnover** | Turns per year by item | Period, Item Group, Warehouse |
| **Dead Stock Report** | Items with no movement | Days without Movement, Warehouse |

### 15.2 Custom Report Builder

```
Navigate: Reports → Custom Report Builder → Stock

Drag-and-drop columns from:
  Items table:         item_code, item_name, item_group, uom
  SLE table:           posting_datetime, qty_change, valuation_rate
  Warehouse table:     warehouse_name, warehouse_type
  Batch table:         batch_id, expiry_date, mfg_date
  Voucher fields:      voucher_type, voucher_no, remarks

Group By:             Item Group / Warehouse / Month / Supplier
Sort By:              Any column, ASC/DESC
Filters:              Dynamic filter builder
Scheduled Export:     Daily/Weekly/Monthly → Email or Drive
```

### 15.3 KPI Dashboard

The Stock KPI Dashboard provides real-time operational visibility:

```
Navigate: Stock → Dashboard

KPI Cards:
  Total Inventory Value         KES 52,155,000
  Items Below Reorder Point     3 items          ⚠
  Batches Expiring in 30 Days   2 batches        ⚠
  Open Material Requests        5
  Pending Purchase Receipts     1
  Stock Entries Pending Approval 0

Charts:
  Inventory Value Trend (90 days)
  Top 10 Items by Value
  Warehouse-wise Stock Distribution
  Daily Wetstock Variance (last 30 days)
  Inventory Turnover by Category
```

### 15.4 Anomaly Detection Alerts

AwoERP's Governance module receives stock anomaly signals and surfaces them for review:

| Anomaly Type | Trigger Condition | Severity |
|---|---|---|
| **Velocity Spike** | Item moves > 3σ above its 30-day average movement rate | High |
| **Valuation Outlier** | Purchase rate deviates > 20% from rolling 90-day average | Medium |
| **Wetstock Over-Threshold** | Daily variance > EPRA 0.5% tolerance | High |
| **Serial Number Conflict** | Same serial number in two locations simultaneously | Critical |
| **Expiry Override** | Batch past expiry date issued with override | High |
| **Backdated Entry** | Stock entry posted to a date > 7 days in the past | Medium |
| **Negative Stock** | Stock balance goes negative (if allowed) | High |
| **Reorder Ignored** | Item below reorder for > lead_time_days × 2 without action | Medium |

Each alert is routed to the relevant role (Warehouse Manager, Finance Controller, Station Manager) and requires a documented disposition.

---

## 16. API Reference & Integration

### 16.1 REST API Endpoints

All endpoints are tenant-scoped via subdomain routing (`{tenant}.awoerp.com`) and require session authentication.

#### Stock Balance

```
GET /api/v1/stock/balance

Query Parameters:
  item_code     VARCHAR    Filter by item
  warehouse_id  UUID       Filter by warehouse (includes children)
  as_of_date    DATE       Balance as at this date (default: today)
  include_zero  BOOL       Include items with zero balance

Response:
{
  "data": [
    {
      "item_code": "PETROL-ULP-95",
      "item_name": "Unleaded Petrol 95",
      "warehouse_id": "wh-maanzoni-tanks",
      "warehouse_name": "Maanzoni Forecourt Tanks",
      "qty": 57150.000,
      "uom": "L",
      "valuation_rate": 166.67,
      "stock_value": 9524670.50,
      "currency": "KES",
      "last_updated": "2026-06-22T14:00:00Z"
    }
  ],
  "total_value": 52155000.00,
  "currency": "KES"
}
```

#### Stock Ledger Query

```
GET /api/v1/stock/ledger

Query Parameters:
  item_code       VARCHAR
  warehouse_id    UUID
  from_date       DATE
  to_date         DATE
  voucher_type    ENUM
  page            INT     (default: 1)
  per_page        INT     (default: 50, max: 500)

Response:
{
  "data": [
    {
      "sle_id": "uuid",
      "posting_datetime": "2026-06-22T06:00:00+03:00",
      "item_code": "PETROL-ULP-95",
      "warehouse_name": "Maanzoni Forecourt Tanks",
      "qty_change": 30000.000,
      "valuation_rate": 168.00,
      "stock_value": 5040000.00,
      "balance_qty": 57150.000,
      "balance_value": 9524670.50,
      "voucher_type": "purchase_receipt",
      "voucher_no": "GRN-2026-00089",
      "batch_no": "VIVO-2026-06-A"
    }
  ],
  "pagination": { "page": 1, "per_page": 50, "total": 1, "pages": 1 }
}
```

#### Create Stock Entry

```
POST /api/v1/stock/entries

Request Body:
{
  "purpose": "material_transfer",
  "posting_date": "2026-06-22",
  "from_warehouse_id": "wh-maanzoni-tanks",
  "to_warehouse_id": "wh-maanzoni-lubes",
  "items": [
    {
      "item_code": "SHELL-HELIX-10W40-4L",
      "qty": 24,
      "uom": "EA",
      "batch_no": "SH-2025-KE-0441"
    }
  ],
  "remarks": "Transfer for lube bay display"
}

Response:
{
  "entry_no": "STE-2026-00234",
  "status": "submitted",
  "sle_ids": ["uuid1", "uuid2"],
  "gl_entry_ids": ["uuid3", "uuid4"]
}
```

#### Batch Trace

```
GET /api/v1/stock/batches/{batch_id}/trace

Response:
{
  "batch_id": "VIVO-2026-06-A",
  "item_code": "PETROL-ULP-95",
  "origin": { ... },
  "movements": [ ... ],
  "current_balance": { "qty": 57150.0, "warehouse": "Maanzoni Forecourt Tanks" }
}
```

### 16.2 Webhook Events

AwoERP publishes the following stock events to configured webhook endpoints:

| Event | Trigger | Payload |
|---|---|---|
| `stock.entry.submitted` | Any stock entry is submitted | entry_no, type, items, warehouses |
| `stock.balance.low` | Item drops below reorder point | item, warehouse, qty, reorder_point |
| `stock.balance.critical` | Item drops below safety stock | item, warehouse, qty |
| `stock.batch.expiry_soon` | Batch within notification window | batch_id, item, expiry_date, qty |
| `stock.batch.expired` | Batch expiry date passed | batch_id, item, qty_remaining |
| `stock.anomaly.detected` | Anomaly signal raised | anomaly_type, severity, details |
| `stock.reconciliation.variance` | Post-count variance above threshold | warehouse, variance_value, variance_pct |
| `stock.wetstock.variance` | Wetstock daily variance above EPRA threshold | tank, variance_litres, variance_pct |

### 16.3 Third-Party Integrations

#### POS Systems (Forecourt PTS-2 / Wayne / Gilbarco)

The forecourt pump controller integration pushes meter readings directly into AwoERP at shift close via a locally-installed agent:

```
PTS-2 Controller → Local Agent (TCP/IP) → AwoERP Forecourt API
                                              (POST /api/v1/forecourt/meter-readings/batch)
```

#### Barcode & RFID Scanners

AwoERP's mobile stock app (React Native) communicates with Bluetooth scanners. Scanned barcodes resolve to item codes via:

```
GET /api/v1/stock/items/barcode/{barcode}
→ Returns item details, current batch suggestion, reorder status
```

#### EFT / Fleet Card Processors

Fleet card transactions (Shell Card, etc.) are imported and reconciled against Delivery Notes:

```
Navigate: Stock → Integrations → Fleet Card Import

Import Format:  CSV / API
Match by:       Transaction ID, Vehicle Registration, Product Code
Auto-create:    Delivery Notes on successful match
Exceptions:     Unmatched transactions routed to manual review queue
```

---

## 17. Feature Flags & Permissions

### 17.1 Stock Module Feature Flags

Feature flags allow granular enablement of sub-features per tenant. All flags default to `false` (opt-in).

```
Navigate: Settings → Feature Flags → Stock

Flag Key                              Default  Description
────────────────────────────────────────────────────────────────────────
stock.module.enabled                  false    Master stock enable
stock.bin_tracking                    false    Sub-warehouse bin locations
stock.batch_tracking                  false    Batch management
stock.serial_tracking                 false    Serial number tracking
stock.expiry_management               false    Expiry date alerts and FEFO
stock.landed_costs                    false    Landed cost vouchers
stock.quality_inspection              false    QC inspection on receipt/delivery
stock.mrp                             false    Material Requirement Planning
stock.forecourt.wetstock              false    Petroleum wetstock module
stock.forecourt.meter_readings        false    Pump meter tracking
stock.inter_company_transfer          false    Cross-entity transfers
stock.negative_stock_override         false    Allow negative stock (Finance approval)
stock.backdated_entries               false    Allow past-date stock entries
stock.anomaly_detection               true     Send anomaly signals to Governance
stock.auto_reorder                    false    Auto-create Material Requests
```

### 17.2 Role-Based Access Control

#### Standard Stock Roles

| Role | Description | Typical User |
|---|---|---|
| `stock.manager` | Full access to all stock functions | Warehouse Manager |
| `stock.supervisor` | Approve transactions, view all, create entries | Senior Attendant / Supervisor |
| `stock.attendant` | Create entries, view own warehouse | Pump Attendant, Store Clerk |
| `stock.auditor` | Read-only access to all stock data + reports | Internal Auditor |
| `stock.purchasing` | Manage purchase receipts, landed costs | Procurement Officer |
| `stock.qc` | Create and approve quality inspections | QC Inspector |
| `stock.readonly` | View balances and ledger only | Management / Finance |

#### Permission Matrix

| Action | Manager | Supervisor | Attendant | Auditor | Readonly |
|---|---|---|---|---|---|
| View Stock Balance | ✓ | ✓ | Own WH | ✓ | ✓ |
| Create Stock Entry | ✓ | ✓ | ✓ | ✗ | ✗ |
| Submit Stock Entry | ✓ | ✓ | Limited | ✗ | ✗ |
| Amend Submitted Entry | ✓ | ✗ | ✗ | ✗ | ✗ |
| Post Reconciliation | ✓ | Approve | Submit | ✗ | ✗ |
| Override Expiry Block | ✓ | ✓ | ✗ | ✗ | ✗ |
| Backdate Entry | ✓ | ✗ | ✗ | ✗ | ✗ |
| Delete Entry | ✓ | ✗ | ✗ | ✗ | ✗ |
| View Anomaly Alerts | ✓ | ✓ | ✗ | ✓ | ✗ |
| Manage Warehouse Config | ✓ | ✗ | ✗ | ✗ | ✗ |

---

## 18. Audit Trail & Compliance

### 18.1 Immutable Stock Ledger

Every SLE is **append-only**. The Go application layer enforces this — there are no `UPDATE` or `DELETE` paths for SLEs in the data access layer. Corrections always create new SLEs:

- **Quantity correction:** Reversal entry (negative of original) + new correct entry
- **Valuation correction:** Stock revaluation entry
- **Wrong warehouse:** Transfer entry (not a deletion and re-entry)

### 18.2 Generalised Audit Log

All stock document state changes (creation, submission, amendment, cancellation) are recorded in the generalised audit log:

```sql
-- audit.event_log
event_id         UUID
tenant_id        UUID
occurred_at      TIMESTAMPTZ
entity_type      VARCHAR     -- 'stock_entry', 'purchase_receipt', etc.
entity_id        UUID
action           VARCHAR     -- 'create', 'submit', 'amend', 'cancel', 'approve'
actor_id         UUID
actor_ip         INET
before_state     JSONB       -- serialised document before change
after_state      JSONB       -- serialised document after change
sensitivity      ENUM        -- 'low', 'medium', 'high', 'critical'
```

High-sensitivity events (reconciliation postings, backdated entries, expiry overrides) are retained for 7 years per KRA record-keeping requirements.

### 18.3 KRA Compliance Notes

AwoERP's stock module is designed with KRA audit requirements in mind:

- Stock values reconcile to GL at every period close (no off-balance stock)
- All entries carry posting dates and cannot be antedated beyond the configured freeze date
- FIFO queue is preserved for full reconstruction of historical cost flows
- Wetstock records comply with EPRA's daily reconciliation requirements
- Purchase receipts link to supplier KRA PIN and invoice numbers for input VAT claim support

---

## 19. Troubleshooting & FAQs

### Q: Stock balance shows negative even though I have stock

**Cause:** Most commonly, a stock entry is awaiting submission (draft status). Draft entries do not affect stock balances.

**Resolution:** Navigate to **Stock → Stock Entries** and look for any entries in Draft status for the affected item/warehouse. Submit or discard them.

### Q: AVCO rate jumped unexpectedly after a return

**Cause:** When stock is returned from a customer and the return rate is different from the current AVCO, this recalculates the average. Returns are treated as an inflow with the return value.

**Resolution:** This is correct behaviour. If the return rate was entered incorrectly, amend the return entry. Check that the return rate matches the original delivery rate.

### Q: Purchase Receipt is stuck in "Pending QC Inspection"

**Cause:** The item or item group has `inspection_required_on_purchase = true`. A Quality Inspection must be created, readings entered, and the result set to Accepted before the receipt can be submitted.

**Resolution:** Navigate to **Stock → Quality Inspections**, find the inspection linked to this receipt, enter readings, and save with result = Accepted.

### Q: Wetstock variance is always slightly negative

**Cause:** Minor losses from evaporation, thermal contraction, and measurement imprecision are normal and expected, particularly for petrol. Allowable evaporation is typically 0.1–0.15% of throughput.

**Resolution:** Review the **Allowable Loss Configuration** under `Stock → Petroleum → Settings`. If actual variance is consistently within the allowable band, no action is required. Adjust the evaporation factor if needed based on tank size, product type, and climate.

### Q: Item below reorder point but no Material Request was auto-created

**Cause:** The `stock.auto_reorder` feature flag is disabled for your tenant, OR the replenishment scheduler has not run since the stock dropped below the threshold (runs hourly), OR the item has an open/pending Material Request already.

**Resolution:** Check feature flags. You can manually trigger replenishment via **Stock → Replenishment → Run Now**. Check for existing open MRs for the item.

### Q: Serial number shows in wrong warehouse

**Cause:** A stock entry was submitted with an incorrect destination warehouse, or a delivery note was created against the wrong source warehouse.

**Resolution:** Create an amending stock transfer to move the serial number to the correct warehouse. The transfer will update the serial number's `current_warehouse` field. If the original entry is wrong, amend it (if permitted) or create a reversal + corrective entry.

### Q: GL and stock values do not reconcile at period close

**Cause:** Stock transactions processed without the `auto_post_gl` setting enabled, OR GL accounts were changed after transactions were posted.

**Resolution:** Run the **Stock-GL Reconciliation Report** under `Reports → Finance → Stock GL Reconciliation`. It identifies specific stock entries that lack corresponding GL entries. A reconciliation job can re-post missing GL entries (requires Finance Controller role).

---

## 20. Glossary

| Term | Definition |
|---|---|
| **AVCO** | Average Cost — a valuation method where the per-unit cost is the weighted average of all units in stock |
| **Batch** | A group of inventory units received together, sharing an origin, expiry date, and lot number |
| **Bin** | A physical sub-location within a warehouse (rack, shelf, tank compartment) |
| **COGS** | Cost of Goods Sold — the inventory cost recognised when goods are sold or consumed |
| **Cycle Count** | A partial physical count covering a subset of items on a rotating schedule |
| **EPRA** | Energy and Petroleum Regulatory Authority — Kenya's petroleum sector regulator |
| **FEFO** | First Expired, First Out — issue the batch with the earliest expiry date first |
| **FIFO** | First In, First Out — issue or value stock in the order it was received |
| **GL** | General Ledger — the central accounting record into which stock movements post |
| **GRN** | Goods Received Note — a Purchase Receipt document |
| **KRA** | Kenya Revenue Authority |
| **Landed Cost** | Total cost of goods including purchase price, freight, customs, insurance, and handling |
| **Material Request** | An internal request to purchase or transfer stock |
| **MRP** | Material Requirement Planning — computing what to order based on demand and current stock |
| **NEMA** | National Environment Management Authority — Kenya's environmental regulator |
| **PPV** | Purchase Price Variance — difference between standard cost and actual purchase price |
| **QC Hold** | A quarantine warehouse type where suspect stock is held pending inspection |
| **Reorder Point** | The stock level at which a replenishment order should be triggered |
| **Safety Stock** | A buffer quantity held to absorb demand or supply uncertainty |
| **SLE** | Stock Ledger Entry — the atomic record of a single stock movement |
| **Standard Cost** | A predetermined per-unit cost used in Standard Costing method |
| **UOM** | Unit of Measure — the quantity unit (Litre, Kilogram, Each, etc.) |
| **UST** | Underground Storage Tank — fuel storage tank at a service station |
| **Valuation Rate** | The per-unit monetary value assigned to inventory at a point in time |
| **Wetstock** | Liquid fuel held in bulk storage tanks at a service station |
| **WIP** | Work in Progress — inventory that is partially through a manufacturing process |

---

*This document is maintained as part of the AwoERP platform documentation suite. For the developer API reference, integration guides, and release changelog, refer to the `docs/` directory in the AwoERP monorepo. For support, contact the AwoERP platform team.*

*© 2026 Anika Global Limited / AwoERP Platform. All rights reserved.*
