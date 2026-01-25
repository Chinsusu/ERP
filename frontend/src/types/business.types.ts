// ============================================
// MATERIAL TYPES
// ============================================

export interface Material {
    id: string
    material_code: string
    name: string
    description?: string
    material_type: 'RAW_MATERIAL' | 'PACKAGING' | 'SEMI_FINISHED'
    category_id?: string
    category?: Category

    // Specifications
    inci_name?: string
    cas_number?: string
    origin_country?: string
    is_organic: boolean
    is_natural: boolean

    // Storage
    storage_conditions?: string
    requires_cold_storage: boolean
    shelf_life_days?: number

    // Safety
    is_hazardous: boolean
    allergen_info?: string
    safety_data_sheet_url?: string

    // Cost
    standard_cost?: number
    currency: string

    // Inventory
    base_uom_id: string
    base_uom?: UnitOfMeasure
    min_stock_quantity?: number
    max_stock_quantity?: number
    reorder_point?: number

    // Status
    is_active: boolean

    // Audit
    created_by: string
    created_at: string
    updated_at: string
}

export interface Category {
    id: string
    code: string
    name: string
    parent_id?: string
    category_type: 'MATERIAL' | 'PRODUCT'
    description?: string
    is_active: boolean
    children?: Category[]
}

export interface UnitOfMeasure {
    id: string
    code: string
    name: string
    symbol?: string
    uom_type: 'WEIGHT' | 'VOLUME' | 'LENGTH' | 'QUANTITY' | 'TIME'
    base_unit_id?: string
    conversion_factor?: number
    is_active: boolean
}

export interface Product {
    id: string
    product_code: string
    name: string
    description?: string
    category_id?: string
    category?: Category
    product_line?: string
    volume?: number
    volume_unit?: string
    standard_price?: number
    barcode?: string
    sku?: string
    shelf_life_months?: number
    is_active: boolean
    created_at: string
    updated_at: string
}

// ============================================
// SUPPLIER TYPES
// ============================================

export interface Supplier {
    id: string
    supplier_code: string
    name: string
    legal_name?: string
    supplier_type: 'MANUFACTURER' | 'DISTRIBUTOR' | 'TRADER'
    tax_id?: string

    // Contact
    address?: string
    city?: string
    country?: string
    phone?: string
    email?: string
    website?: string

    // Business
    payment_terms?: number
    credit_limit?: number
    currency: string

    // Status
    status: 'ACTIVE' | 'INACTIVE' | 'PENDING' | 'BLOCKED'
    is_approved: boolean
    approved_date?: string

    // Quality
    quality_rating?: number
    delivery_rating?: number
    overall_rating?: number

    // Contacts & Certifications
    contacts?: SupplierContact[]
    certifications?: SupplierCertification[]

    created_at: string
    updated_at: string
}

export interface SupplierContact {
    id: string
    supplier_id: string
    name: string
    position?: string
    department?: string
    phone?: string
    email?: string
    is_primary: boolean
}

export interface SupplierCertification {
    id: string
    supplier_id: string
    certification_type: string
    certificate_number: string
    issuing_body: string
    issue_date: string
    expiry_date: string
    document_url?: string
    status: 'VALID' | 'EXPIRED' | 'PENDING'
}

export interface SupplierEvaluation {
    id: string
    supplier_id: string
    evaluation_date: string
    evaluator_id: string
    quality_score: number
    delivery_score: number
    price_score: number
    service_score: number
    overall_score: number
    comments?: string
}

// ============================================
// PROCUREMENT TYPES
// ============================================

export interface PurchaseRequisition {
    id: string
    pr_number: string
    title: string
    description?: string
    requester_id: string
    requester_name?: string
    department_id?: string
    department_name?: string

    status: 'DRAFT' | 'SUBMITTED' | 'APPROVED' | 'REJECTED' | 'CONVERTED' | 'CANCELLED'
    priority: 'LOW' | 'NORMAL' | 'HIGH' | 'URGENT'

    required_date?: string

    items: PRItem[]

    total_amount?: number
    currency: string

    submitted_at?: string
    approved_at?: string
    approved_by?: string
    rejection_reason?: string

    created_at: string
    updated_at: string
}

export interface PRItem {
    id: string
    pr_id: string
    line_number: number
    material_id: string
    material?: Material
    description?: string
    quantity: number
    uom_id: string
    uom?: UnitOfMeasure
    estimated_price?: number
    preferred_supplier_id?: string
    notes?: string
}

export interface PurchaseOrder {
    id: string
    po_number: string
    pr_id?: string
    supplier_id: string
    supplier?: Supplier

    status: 'DRAFT' | 'PENDING_APPROVAL' | 'APPROVED' | 'CONFIRMED' | 'PARTIALLY_RECEIVED' | 'RECEIVED' | 'CANCELLED'

    order_date: string
    expected_delivery_date?: string
    actual_delivery_date?: string

    items: POItem[]

    subtotal: number
    tax_amount: number
    total_amount: number
    currency: string

    payment_terms?: number
    shipping_address?: string
    notes?: string

    approved_by?: string
    approved_at?: string
    confirmed_at?: string

    created_at: string
    updated_at: string
}

export interface POItem {
    id: string
    po_id: string
    line_number: number
    material_id: string
    material?: Material
    description?: string
    quantity: number
    received_quantity: number
    uom_id: string
    uom?: UnitOfMeasure
    unit_price: number
    total_price: number
    tax_rate?: number
}

// ============================================
// WMS TYPES
// ============================================

export interface Stock {
    id: string
    material_id: string
    material?: Material
    warehouse_id: string
    location_id?: string

    quantity: number
    reserved_quantity: number
    available_quantity: number
    uom_id: string
    uom?: UnitOfMeasure

    last_movement_date?: string
}

export interface Lot {
    id: string
    lot_number: string
    material_id: string
    material?: Material

    quantity: number
    available_quantity: number
    uom_id: string
    uom?: UnitOfMeasure

    manufacturing_date?: string
    expiry_date?: string
    received_date: string

    supplier_lot_number?: string
    supplier_id?: string

    status: 'AVAILABLE' | 'RESERVED' | 'QUARANTINE' | 'EXPIRED' | 'CONSUMED'

    warehouse_id: string
    location_id?: string

    qc_status?: 'PENDING' | 'PASSED' | 'FAILED'

    created_at: string
}

export interface GoodsReceiptNote {
    id: string
    grn_number: string
    po_id?: string
    po_number?: string
    supplier_id: string
    supplier?: Supplier

    status: 'DRAFT' | 'PENDING_QC' | 'QC_PASSED' | 'QC_FAILED' | 'COMPLETED' | 'CANCELLED'

    receipt_date: string
    warehouse_id: string

    items: GRNItem[]

    received_by: string
    notes?: string

    created_at: string
    updated_at: string
}

export interface GRNItem {
    id: string
    grn_id: string
    line_number: number
    po_item_id?: string
    material_id: string
    material?: Material

    expected_quantity?: number
    received_quantity: number
    accepted_quantity: number
    rejected_quantity: number
    uom_id: string
    uom?: UnitOfMeasure

    lot_number?: string
    manufacturing_date?: string
    expiry_date?: string

    qc_status?: 'PENDING' | 'PASSED' | 'FAILED'
    rejection_reason?: string
}

export interface StockMovement {
    id: string
    movement_type: 'RECEIPT' | 'ISSUE' | 'TRANSFER' | 'ADJUSTMENT' | 'RETURN'
    reference_type: string
    reference_id: string
    reference_number: string

    material_id: string
    lot_id?: string

    quantity: number
    uom_id: string

    from_warehouse_id?: string
    from_location_id?: string
    to_warehouse_id?: string
    to_location_id?: string

    movement_date: string
    performed_by: string
    notes?: string
}

// ============================================
// MANUFACTURING TYPES
// ============================================

export interface BillOfMaterials {
    id: string
    bom_code: string
    product_id: string
    product?: Product
    version: string

    status: 'DRAFT' | 'ACTIVE' | 'OBSOLETE'

    yield_quantity: number
    yield_uom_id: string
    yield_uom?: UnitOfMeasure

    items: BOMItem[]

    // Encrypted formula (permission required)
    formula_details?: string

    effective_date?: string
    expiry_date?: string

    created_by: string
    approved_by?: string
    approved_at?: string

    created_at: string
    updated_at: string
}

export interface BOMItem {
    id: string
    bom_id: string
    line_number: number
    material_id: string
    material?: Material

    quantity: number
    uom_id: string
    uom?: UnitOfMeasure

    waste_percentage?: number
    is_critical: boolean
    notes?: string
}

export interface WorkOrder {
    id: string
    wo_number: string
    bom_id: string
    bom?: BillOfMaterials
    product_id: string
    product?: Product

    status: 'DRAFT' | 'PLANNED' | 'RELEASED' | 'IN_PROGRESS' | 'COMPLETED' | 'CANCELLED'
    priority: 'LOW' | 'NORMAL' | 'HIGH' | 'URGENT'

    planned_quantity: number
    completed_quantity: number
    rejected_quantity: number
    uom_id: string
    uom?: UnitOfMeasure

    planned_start_date: string
    planned_end_date: string
    actual_start_date?: string
    actual_end_date?: string

    material_issues: MaterialIssue[]

    assigned_to?: string
    notes?: string

    created_at: string
    updated_at: string
}

export interface MaterialIssue {
    id: string
    work_order_id: string
    material_id: string
    material?: Material
    lot_id?: string
    lot?: Lot

    required_quantity: number
    issued_quantity: number
    uom_id: string

    issued_at?: string
    issued_by?: string
}

export interface QCInspection {
    id: string
    inspection_number: string
    reference_type: 'GRN' | 'WORK_ORDER' | 'FINISHED_GOODS'
    reference_id: string

    status: 'PENDING' | 'IN_PROGRESS' | 'PASSED' | 'FAILED'

    material_id?: string
    product_id?: string
    lot_id?: string

    checkpoints: QCCheckpoint[]

    inspector_id: string
    inspection_date: string

    overall_result?: 'PASS' | 'FAIL' | 'CONDITIONAL'
    comments?: string

    created_at: string
}

export interface QCCheckpoint {
    id: string
    inspection_id: string
    checkpoint_name: string
    checkpoint_type: 'VISUAL' | 'MEASUREMENT' | 'TEST'

    specification?: string
    min_value?: number
    max_value?: number

    actual_value?: string
    result: 'PASS' | 'FAIL' | 'PENDING'
    notes?: string
}

export interface TraceabilityRecord {
    id: string
    record_type: 'RECEIPT' | 'ISSUE' | 'PRODUCTION' | 'SHIPMENT'

    lot_id?: string
    lot_number?: string
    material_id?: string
    product_id?: string

    source_reference?: string
    target_reference?: string

    quantity: number
    uom_id: string

    event_date: string

    parent_records?: TraceabilityRecord[]
    child_records?: TraceabilityRecord[]
}
