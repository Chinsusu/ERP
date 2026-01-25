// Mock data for testing frontend without backend

import type {
    Material, Category, UnitOfMeasure,
    Supplier, SupplierContact, SupplierCertification,
    PurchaseRequisition, PurchaseOrder,
    Stock, Lot, GoodsReceiptNote,
    BillOfMaterials, WorkOrder
} from '@/types/business.types'

// ============================================
// COMMON DATA
// ============================================

export const mockCategories: Category[] = [
    { id: '1', code: 'RM', name: 'Raw Materials', category_type: 'MATERIAL', is_active: true },
    { id: '2', code: 'PKG', name: 'Packaging', category_type: 'MATERIAL', is_active: true },
    { id: '3', code: 'ACT', name: 'Active Ingredients', parent_id: '1', category_type: 'MATERIAL', is_active: true },
    { id: '4', code: 'EMU', name: 'Emulsifiers', parent_id: '1', category_type: 'MATERIAL', is_active: true },
]

export const mockUoMs: UnitOfMeasure[] = [
    { id: '1', code: 'KG', name: 'Kilogram', symbol: 'kg', uom_type: 'WEIGHT', is_active: true },
    { id: '2', code: 'G', name: 'Gram', symbol: 'g', uom_type: 'WEIGHT', base_unit_id: '1', conversion_factor: 0.001, is_active: true },
    { id: '3', code: 'L', name: 'Liter', symbol: 'L', uom_type: 'VOLUME', is_active: true },
    { id: '4', code: 'ML', name: 'Milliliter', symbol: 'mL', uom_type: 'VOLUME', base_unit_id: '3', conversion_factor: 0.001, is_active: true },
    { id: '5', code: 'PCS', name: 'Pieces', symbol: 'pcs', uom_type: 'QUANTITY', is_active: true },
]

// ============================================
// MATERIALS
// ============================================

export const mockMaterials: Material[] = [
    {
        id: '1',
        material_code: 'RM-001',
        name: 'Hyaluronic Acid',
        description: 'High molecular weight hyaluronic acid for skincare',
        material_type: 'RAW_MATERIAL',
        category_id: '3',
        category: mockCategories[2],
        inci_name: 'Sodium Hyaluronate',
        cas_number: '9067-32-7',
        origin_country: 'South Korea',
        is_organic: false,
        is_natural: true,
        storage_conditions: 'Store in cool, dry place',
        requires_cold_storage: false,
        shelf_life_days: 730,
        is_hazardous: false,
        standard_cost: 2500000,
        currency: 'VND',
        base_uom_id: '1',
        base_uom: mockUoMs[0],
        min_stock_quantity: 5,
        max_stock_quantity: 50,
        reorder_point: 10,
        is_active: true,
        created_by: 'admin',
        created_at: '2025-01-15T10:00:00Z',
        updated_at: '2025-01-20T14:30:00Z'
    },
    {
        id: '2',
        material_code: 'RM-002',
        name: 'Niacinamide',
        description: 'Vitamin B3 for skin brightening',
        material_type: 'RAW_MATERIAL',
        category_id: '3',
        category: mockCategories[2],
        inci_name: 'Niacinamide',
        cas_number: '98-92-0',
        origin_country: 'China',
        is_organic: false,
        is_natural: false,
        storage_conditions: 'Store below 25°C',
        requires_cold_storage: false,
        shelf_life_days: 730,
        is_hazardous: false,
        standard_cost: 800000,
        currency: 'VND',
        base_uom_id: '1',
        base_uom: mockUoMs[0],
        min_stock_quantity: 10,
        max_stock_quantity: 100,
        reorder_point: 20,
        is_active: true,
        created_by: 'admin',
        created_at: '2025-01-10T09:00:00Z',
        updated_at: '2025-01-18T11:00:00Z'
    },
    {
        id: '3',
        material_code: 'RM-003',
        name: 'Cetearyl Alcohol',
        description: 'Fatty alcohol emulsifier',
        material_type: 'RAW_MATERIAL',
        category_id: '4',
        category: mockCategories[3],
        inci_name: 'Cetearyl Alcohol',
        cas_number: '67762-27-0',
        origin_country: 'Malaysia',
        is_organic: false,
        is_natural: true,
        storage_conditions: 'Store in cool, dry place',
        requires_cold_storage: false,
        shelf_life_days: 730,
        is_hazardous: false,
        standard_cost: 150000,
        currency: 'VND',
        base_uom_id: '1',
        base_uom: mockUoMs[0],
        min_stock_quantity: 50,
        max_stock_quantity: 500,
        reorder_point: 100,
        is_active: true,
        created_by: 'admin',
        created_at: '2025-01-05T08:00:00Z',
        updated_at: '2025-01-12T16:00:00Z'
    },
    {
        id: '4',
        material_code: 'PKG-001',
        name: 'Glass Bottle 30ml',
        description: 'Frosted glass bottle with dropper',
        material_type: 'PACKAGING',
        category_id: '2',
        category: mockCategories[1],
        is_organic: false,
        is_natural: false,
        requires_cold_storage: false,
        is_hazardous: false,
        standard_cost: 25000,
        currency: 'VND',
        base_uom_id: '5',
        base_uom: mockUoMs[4],
        min_stock_quantity: 500,
        max_stock_quantity: 5000,
        reorder_point: 1000,
        is_active: true,
        created_by: 'admin',
        created_at: '2025-01-08T10:00:00Z',
        updated_at: '2025-01-15T12:00:00Z'
    },
    {
        id: '5',
        material_code: 'PKG-002',
        name: 'Pump Dispenser 50ml',
        description: 'Airless pump dispenser bottle',
        material_type: 'PACKAGING',
        category_id: '2',
        category: mockCategories[1],
        is_organic: false,
        is_natural: false,
        requires_cold_storage: false,
        is_hazardous: false,
        standard_cost: 35000,
        currency: 'VND',
        base_uom_id: '5',
        base_uom: mockUoMs[4],
        min_stock_quantity: 300,
        max_stock_quantity: 3000,
        reorder_point: 600,
        is_active: true,
        created_by: 'admin',
        created_at: '2025-01-09T11:00:00Z',
        updated_at: '2025-01-16T14:00:00Z'
    },
]

// ============================================
// SUPPLIERS
// ============================================

export const mockSuppliers: Supplier[] = [
    {
        id: '1',
        supplier_code: 'SUP-001',
        name: 'Korea Cosmetic Ingredients Co.',
        legal_name: 'Korea Cosmetic Ingredients Corporation',
        supplier_type: 'MANUFACTURER',
        tax_id: '123-45-67890',
        address: '123 Gangnam-daero, Gangnam-gu',
        city: 'Seoul',
        country: 'South Korea',
        phone: '+82-2-1234-5678',
        email: 'sales@kci.co.kr',
        website: 'https://www.kci.co.kr',
        payment_terms: 30,
        credit_limit: 500000000,
        currency: 'VND',
        status: 'ACTIVE',
        is_approved: true,
        approved_date: '2024-06-15',
        quality_rating: 5,
        delivery_rating: 4,
        overall_rating: 5,
        created_at: '2024-06-01T10:00:00Z',
        updated_at: '2025-01-20T09:00:00Z'
    },
    {
        id: '2',
        supplier_code: 'SUP-002',
        name: 'VietPack Solutions',
        legal_name: 'VietPack Solutions JSC',
        supplier_type: 'MANUFACTURER',
        tax_id: '0312345678',
        address: '456 Nguyen Van Linh, District 7',
        city: 'Ho Chi Minh City',
        country: 'Vietnam',
        phone: '+84-28-1234-5678',
        email: 'info@vietpack.vn',
        website: 'https://www.vietpack.vn',
        payment_terms: 15,
        credit_limit: 200000000,
        currency: 'VND',
        status: 'ACTIVE',
        is_approved: true,
        approved_date: '2024-08-20',
        quality_rating: 4,
        delivery_rating: 5,
        overall_rating: 4,
        created_at: '2024-08-10T09:00:00Z',
        updated_at: '2025-01-18T11:00:00Z'
    },
    {
        id: '3',
        supplier_code: 'SUP-003',
        name: 'ChemTrade China',
        legal_name: 'ChemTrade International Ltd',
        supplier_type: 'DISTRIBUTOR',
        tax_id: '91320000MA1MTXXX',
        address: '789 Pudong New Area',
        city: 'Shanghai',
        country: 'China',
        phone: '+86-21-1234-5678',
        email: 'export@chemtrade.cn',
        payment_terms: 45,
        credit_limit: 300000000,
        currency: 'VND',
        status: 'ACTIVE',
        is_approved: true,
        approved_date: '2024-09-10',
        quality_rating: 4,
        delivery_rating: 3,
        overall_rating: 4,
        created_at: '2024-09-01T08:00:00Z',
        updated_at: '2025-01-15T10:00:00Z'
    },
]

// ============================================
// PURCHASE REQUISITIONS
// ============================================

export const mockPRs: PurchaseRequisition[] = [
    {
        id: '1',
        pr_number: 'PR-2025-0001',
        title: 'Q1 Raw Materials Request',
        description: 'Quarterly restocking of active ingredients',
        requester_id: '1',
        requester_name: 'Nguyen Van A',
        department_id: '3',
        department_name: 'Production',
        status: 'APPROVED',
        priority: 'HIGH',
        required_date: '2025-02-15',
        items: [],
        total_amount: 25000000,
        currency: 'VND',
        submitted_at: '2025-01-20T09:00:00Z',
        approved_at: '2025-01-21T14:00:00Z',
        approved_by: '2',
        created_at: '2025-01-19T10:00:00Z',
        updated_at: '2025-01-21T14:00:00Z'
    },
    {
        id: '2',
        pr_number: 'PR-2025-0002',
        title: 'Packaging Materials - February',
        description: 'Bottles and dispensers for new serum line',
        requester_id: '1',
        requester_name: 'Nguyen Van A',
        department_id: '3',
        department_name: 'Production',
        status: 'SUBMITTED',
        priority: 'NORMAL',
        required_date: '2025-02-28',
        items: [],
        total_amount: 15000000,
        currency: 'VND',
        submitted_at: '2025-01-22T11:00:00Z',
        created_at: '2025-01-22T10:00:00Z',
        updated_at: '2025-01-22T11:00:00Z'
    },
    {
        id: '3',
        pr_number: 'PR-2025-0003',
        title: 'Emergency - Hyaluronic Acid',
        description: 'Urgent restock due to high demand',
        requester_id: '3',
        requester_name: 'Tran Thi B',
        department_id: '3',
        department_name: 'Production',
        status: 'DRAFT',
        priority: 'URGENT',
        required_date: '2025-01-30',
        items: [],
        total_amount: 50000000,
        currency: 'VND',
        created_at: '2025-01-24T08:00:00Z',
        updated_at: '2025-01-24T08:00:00Z'
    },
]

// ============================================
// PURCHASE ORDERS
// ============================================

export const mockPOs: PurchaseOrder[] = [
    {
        id: '1',
        po_number: 'PO-2025-0001',
        pr_id: '1',
        supplier_id: '1',
        supplier: mockSuppliers[0],
        status: 'CONFIRMED',
        order_date: '2025-01-22',
        expected_delivery_date: '2025-02-10',
        items: [],
        subtotal: 25000000,
        tax_amount: 2500000,
        total_amount: 27500000,
        currency: 'VND',
        payment_terms: 30,
        approved_by: '2',
        approved_at: '2025-01-22T15:00:00Z',
        confirmed_at: '2025-01-23T09:00:00Z',
        created_at: '2025-01-22T14:00:00Z',
        updated_at: '2025-01-23T09:00:00Z'
    },
    {
        id: '2',
        po_number: 'PO-2025-0002',
        supplier_id: '2',
        supplier: mockSuppliers[1],
        status: 'PENDING_APPROVAL',
        order_date: '2025-01-23',
        expected_delivery_date: '2025-02-05',
        items: [],
        subtotal: 15000000,
        tax_amount: 1500000,
        total_amount: 16500000,
        currency: 'VND',
        payment_terms: 15,
        created_at: '2025-01-23T10:00:00Z',
        updated_at: '2025-01-23T10:00:00Z'
    },
    {
        id: '3',
        po_number: 'PO-2025-0003',
        supplier_id: '3',
        supplier: mockSuppliers[2],
        status: 'PARTIALLY_RECEIVED',
        order_date: '2025-01-10',
        expected_delivery_date: '2025-01-25',
        items: [],
        subtotal: 40000000,
        tax_amount: 4000000,
        total_amount: 44000000,
        currency: 'VND',
        payment_terms: 45,
        approved_by: '2',
        approved_at: '2025-01-11T09:00:00Z',
        confirmed_at: '2025-01-12T10:00:00Z',
        created_at: '2025-01-10T14:00:00Z',
        updated_at: '2025-01-20T16:00:00Z'
    },
]

// ============================================
// STOCK
// ============================================

export const mockStock: Stock[] = [
    {
        id: '1',
        material_id: '1',
        material: mockMaterials[0],
        warehouse_id: '1',
        quantity: 15,
        reserved_quantity: 3,
        available_quantity: 12,
        uom_id: '1',
        uom: mockUoMs[0],
        last_movement_date: '2025-01-20T14:00:00Z'
    },
    {
        id: '2',
        material_id: '2',
        material: mockMaterials[1],
        warehouse_id: '1',
        quantity: 45,
        reserved_quantity: 10,
        available_quantity: 35,
        uom_id: '1',
        uom: mockUoMs[0],
        last_movement_date: '2025-01-22T10:00:00Z'
    },
    {
        id: '3',
        material_id: '3',
        material: mockMaterials[2],
        warehouse_id: '1',
        quantity: 80,
        reserved_quantity: 20,
        available_quantity: 60,
        uom_id: '1',
        uom: mockUoMs[0],
        last_movement_date: '2025-01-21T09:00:00Z'
    },
    {
        id: '4',
        material_id: '4',
        material: mockMaterials[3],
        warehouse_id: '1',
        quantity: 1500,
        reserved_quantity: 500,
        available_quantity: 1000,
        uom_id: '5',
        uom: mockUoMs[4],
        last_movement_date: '2025-01-19T16:00:00Z'
    },
    {
        id: '5',
        material_id: '5',
        material: mockMaterials[4],
        warehouse_id: '1',
        quantity: 200,
        reserved_quantity: 0,
        available_quantity: 200,
        uom_id: '5',
        uom: mockUoMs[4],
        last_movement_date: '2025-01-18T11:00:00Z'
    },
]

// ============================================
// GRN
// ============================================

export const mockGRNs: GoodsReceiptNote[] = [
    {
        id: '1',
        grn_number: 'GRN-2025-0001',
        po_id: '3',
        po_number: 'PO-2025-0003',
        supplier_id: '3',
        supplier: mockSuppliers[2],
        status: 'COMPLETED',
        receipt_date: '2025-01-20',
        warehouse_id: '1',
        items: [],
        received_by: '4',
        created_at: '2025-01-20T09:00:00Z',
        updated_at: '2025-01-20T14:00:00Z'
    },
    {
        id: '2',
        grn_number: 'GRN-2025-0002',
        po_id: '1',
        po_number: 'PO-2025-0001',
        supplier_id: '1',
        supplier: mockSuppliers[0],
        status: 'PENDING_QC',
        receipt_date: '2025-01-24',
        warehouse_id: '1',
        items: [],
        received_by: '4',
        created_at: '2025-01-24T10:00:00Z',
        updated_at: '2025-01-24T10:00:00Z'
    },
]

// ============================================
// BOM
// ============================================

export const mockBOMs: BillOfMaterials[] = [
    {
        id: '1',
        bom_code: 'BOM-SRM-001',
        product_id: '1',
        product: { id: '1', product_code: 'PRD-001', name: 'Hyaluronic Acid Serum 30ml', is_active: true, created_at: '', updated_at: '' },
        version: '1.0',
        status: 'ACTIVE',
        yield_quantity: 1000,
        yield_uom_id: '4',
        yield_uom: mockUoMs[3],
        items: [],
        effective_date: '2025-01-01',
        created_by: '1',
        approved_by: '2',
        approved_at: '2025-01-05T10:00:00Z',
        created_at: '2025-01-02T09:00:00Z',
        updated_at: '2025-01-05T10:00:00Z'
    },
    {
        id: '2',
        bom_code: 'BOM-CRM-001',
        product_id: '2',
        product: { id: '2', product_code: 'PRD-002', name: 'Niacinamide Cream 50ml', is_active: true, created_at: '', updated_at: '' },
        version: '2.1',
        status: 'ACTIVE',
        yield_quantity: 500,
        yield_uom_id: '4',
        yield_uom: mockUoMs[3],
        items: [],
        effective_date: '2025-01-10',
        created_by: '1',
        approved_by: '2',
        approved_at: '2025-01-12T14:00:00Z',
        created_at: '2025-01-08T10:00:00Z',
        updated_at: '2025-01-12T14:00:00Z'
    },
    {
        id: '3',
        bom_code: 'BOM-TNR-001',
        product_id: '3',
        product: { id: '3', product_code: 'PRD-003', name: 'Hydrating Toner 150ml', is_active: true, created_at: '', updated_at: '' },
        version: '1.0',
        status: 'DRAFT',
        yield_quantity: 2000,
        yield_uom_id: '4',
        yield_uom: mockUoMs[3],
        items: [],
        created_by: '1',
        created_at: '2025-01-20T11:00:00Z',
        updated_at: '2025-01-20T11:00:00Z'
    },
]

// ============================================
// WORK ORDERS
// ============================================

export const mockWorkOrders: WorkOrder[] = [
    {
        id: '1',
        wo_number: 'WO-2025-0001',
        bom_id: '1',
        bom: mockBOMs[0],
        product_id: '1',
        product: { id: '1', product_code: 'PRD-001', name: 'Hyaluronic Acid Serum 30ml', is_active: true, created_at: '', updated_at: '' },
        status: 'COMPLETED',
        priority: 'NORMAL',
        planned_quantity: 5000,
        completed_quantity: 5000,
        rejected_quantity: 50,
        uom_id: '4',
        uom: mockUoMs[3],
        planned_start_date: '2025-01-10',
        planned_end_date: '2025-01-12',
        actual_start_date: '2025-01-10',
        actual_end_date: '2025-01-12',
        material_issues: [],
        created_at: '2025-01-08T09:00:00Z',
        updated_at: '2025-01-12T17:00:00Z'
    },
    {
        id: '2',
        wo_number: 'WO-2025-0002',
        bom_id: '2',
        bom: mockBOMs[1],
        product_id: '2',
        product: { id: '2', product_code: 'PRD-002', name: 'Niacinamide Cream 50ml', is_active: true, created_at: '', updated_at: '' },
        status: 'IN_PROGRESS',
        priority: 'HIGH',
        planned_quantity: 3000,
        completed_quantity: 1800,
        rejected_quantity: 20,
        uom_id: '4',
        uom: mockUoMs[3],
        planned_start_date: '2025-01-22',
        planned_end_date: '2025-01-25',
        actual_start_date: '2025-01-22',
        material_issues: [],
        created_at: '2025-01-20T10:00:00Z',
        updated_at: '2025-01-24T16:00:00Z'
    },
    {
        id: '3',
        wo_number: 'WO-2025-0003',
        bom_id: '1',
        bom: mockBOMs[0],
        product_id: '1',
        product: { id: '1', product_code: 'PRD-001', name: 'Hyaluronic Acid Serum 30ml', is_active: true, created_at: '', updated_at: '' },
        status: 'PLANNED',
        priority: 'URGENT',
        planned_quantity: 10000,
        completed_quantity: 0,
        rejected_quantity: 0,
        uom_id: '4',
        uom: mockUoMs[3],
        planned_start_date: '2025-01-28',
        planned_end_date: '2025-01-31',
        material_issues: [],
        created_at: '2025-01-24T08:00:00Z',
        updated_at: '2025-01-24T08:00:00Z'
    },
    {
        id: '4',
        wo_number: 'WO-2025-0004',
        bom_id: '2',
        bom: mockBOMs[1],
        product_id: '2',
        product: { id: '2', product_code: 'PRD-002', name: 'Niacinamide Cream 50ml', is_active: true, created_at: '', updated_at: '' },
        status: 'RELEASED',
        priority: 'NORMAL',
        planned_quantity: 2000,
        completed_quantity: 0,
        rejected_quantity: 0,
        uom_id: '4',
        uom: mockUoMs[3],
        planned_start_date: '2025-02-01',
        planned_end_date: '2025-02-03',
        material_issues: [],
        created_at: '2025-01-23T14:00:00Z',
        updated_at: '2025-01-25T09:00:00Z'
    },
]
