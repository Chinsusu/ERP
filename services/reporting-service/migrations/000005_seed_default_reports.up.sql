-- Seed default report definitions
INSERT INTO report_definitions (code, name, description, report_type, query_template, parameters, columns, is_system) VALUES

-- Inventory Reports
('STOCK_SUMMARY', 'Stock Summary', 'Current stock levels by warehouse and material', 'INVENTORY',
'SELECT 
    m.code as material_code,
    m.name as material_name,
    mc.name as category,
    w.name as warehouse,
    COALESCE(SUM(s.quantity), 0) as quantity,
    m.uom,
    m.reorder_point,
    CASE WHEN COALESCE(SUM(s.quantity), 0) < m.reorder_point THEN ''LOW'' ELSE ''OK'' END as status
FROM materials m
LEFT JOIN material_categories mc ON m.category_id = mc.id
LEFT JOIN stock_on_hand s ON m.id = s.material_id
LEFT JOIN warehouses w ON s.warehouse_id = w.id
WHERE ({{.warehouse_id}} IS NULL OR w.id = {{.warehouse_id}})
  AND ({{.category_id}} IS NULL OR m.category_id = {{.category_id}})
GROUP BY m.id, m.code, m.name, mc.name, w.name, m.uom, m.reorder_point
ORDER BY m.code',
'[{"name": "warehouse_id", "type": "uuid", "required": false, "label": "Warehouse"},
  {"name": "category_id", "type": "uuid", "required": false, "label": "Category"}]',
'[{"field": "material_code", "header": "Material Code", "width": 120},
  {"field": "material_name", "header": "Material Name", "width": 200},
  {"field": "category", "header": "Category", "width": 150},
  {"field": "warehouse", "header": "Warehouse", "width": 150},
  {"field": "quantity", "header": "Quantity", "type": "number", "width": 100},
  {"field": "uom", "header": "UoM", "width": 80},
  {"field": "reorder_point", "header": "Reorder Point", "type": "number", "width": 120},
  {"field": "status", "header": "Status", "width": 80}]',
true),

('EXPIRY_REPORT', 'Lot Expiry Report', 'Lots expiring within specified days', 'INVENTORY',
'SELECT 
    l.lot_number,
    m.code as material_code,
    m.name as material_name,
    l.expiry_date,
    (l.expiry_date - CURRENT_DATE) as days_remaining,
    COALESCE(SUM(s.quantity), 0) as quantity,
    m.uom,
    w.name as warehouse
FROM lots l
JOIN materials m ON l.material_id = m.id
LEFT JOIN stock_on_hand s ON l.id = s.lot_id
LEFT JOIN warehouses w ON s.warehouse_id = w.id
WHERE l.expiry_date <= CURRENT_DATE + INTERVAL ''{{.days}} days''
  AND l.status = ''ACTIVE''
  AND COALESCE(s.quantity, 0) > 0
GROUP BY l.id, l.lot_number, m.code, m.name, l.expiry_date, m.uom, w.name
ORDER BY l.expiry_date ASC',
'[{"name": "days", "type": "integer", "required": true, "default": 30, "label": "Days Until Expiry"}]',
'[{"field": "lot_number", "header": "Lot Number", "width": 150},
  {"field": "material_code", "header": "Material Code", "width": 120},
  {"field": "material_name", "header": "Material Name", "width": 200},
  {"field": "expiry_date", "header": "Expiry Date", "type": "date", "width": 120},
  {"field": "days_remaining", "header": "Days Left", "type": "number", "width": 100},
  {"field": "quantity", "header": "Quantity", "type": "number", "width": 100},
  {"field": "uom", "header": "UoM", "width": 80},
  {"field": "warehouse", "header": "Warehouse", "width": 150}]',
true),

('STOCK_MOVEMENT', 'Stock Movement Report', 'Stock ins and outs by period', 'INVENTORY',
'SELECT 
    m.code as material_code,
    m.name as material_name,
    sm.movement_type,
    sm.reference_type,
    sm.reference_number,
    sm.quantity,
    m.uom,
    sm.created_at as movement_date
FROM stock_movements sm
JOIN materials m ON sm.material_id = m.id
WHERE sm.created_at BETWEEN {{.from_date}} AND {{.to_date}}
ORDER BY sm.created_at DESC',
'[{"name": "from_date", "type": "date", "required": true, "label": "From Date"},
  {"name": "to_date", "type": "date", "required": true, "label": "To Date"}]',
'[{"field": "material_code", "header": "Material Code", "width": 120},
  {"field": "material_name", "header": "Material Name", "width": 200},
  {"field": "movement_type", "header": "Type", "width": 100},
  {"field": "reference_type", "header": "Reference Type", "width": 120},
  {"field": "reference_number", "header": "Reference", "width": 150},
  {"field": "quantity", "header": "Quantity", "type": "number", "width": 100},
  {"field": "uom", "header": "UoM", "width": 80},
  {"field": "movement_date", "header": "Date", "type": "datetime", "width": 150}]',
true),

-- Procurement Reports
('PO_SUMMARY', 'Purchase Order Summary', 'PO status summary by period', 'PROCUREMENT',
'SELECT 
    po.po_number,
    s.name as supplier,
    po.status,
    po.total_amount,
    po.currency,
    po.order_date,
    po.expected_delivery_date,
    u.full_name as created_by
FROM purchase_orders po
JOIN suppliers s ON po.supplier_id = s.id
LEFT JOIN users u ON po.created_by = u.id
WHERE po.order_date BETWEEN {{.from_date}} AND {{.to_date}}
ORDER BY po.order_date DESC',
'[{"name": "from_date", "type": "date", "required": true, "label": "From Date"},
  {"name": "to_date", "type": "date", "required": true, "label": "To Date"}]',
'[{"field": "po_number", "header": "PO Number", "width": 150},
  {"field": "supplier", "header": "Supplier", "width": 200},
  {"field": "status", "header": "Status", "width": 120},
  {"field": "total_amount", "header": "Amount", "type": "currency", "width": 120},
  {"field": "currency", "header": "Currency", "width": 80},
  {"field": "order_date", "header": "Order Date", "type": "date", "width": 120},
  {"field": "expected_delivery_date", "header": "Expected Delivery", "type": "date", "width": 140}]',
true),

('SUPPLIER_PERFORMANCE', 'Supplier Performance', 'Supplier ratings and delivery performance', 'PROCUREMENT',
'SELECT 
    s.code as supplier_code,
    s.name as supplier_name,
    s.supplier_type,
    COALESCE(se.overall_score, 0) as rating,
    COUNT(po.id) as total_orders,
    SUM(CASE WHEN po.status = ''COMPLETED'' THEN 1 ELSE 0 END) as completed_orders,
    ROUND(AVG(CASE WHEN grn.received_date <= po.expected_delivery_date THEN 100 ELSE 0 END), 1) as on_time_rate
FROM suppliers s
LEFT JOIN supplier_evaluations se ON s.id = se.supplier_id
LEFT JOIN purchase_orders po ON s.id = po.supplier_id
LEFT JOIN goods_receipt_notes grn ON po.id = grn.po_id
WHERE s.is_active = true
GROUP BY s.id, s.code, s.name, s.supplier_type, se.overall_score
ORDER BY se.overall_score DESC NULLS LAST',
'[]',
'[{"field": "supplier_code", "header": "Code", "width": 100},
  {"field": "supplier_name", "header": "Supplier Name", "width": 200},
  {"field": "supplier_type", "header": "Type", "width": 120},
  {"field": "rating", "header": "Rating", "type": "number", "width": 80},
  {"field": "total_orders", "header": "Total Orders", "type": "number", "width": 100},
  {"field": "completed_orders", "header": "Completed", "type": "number", "width": 100},
  {"field": "on_time_rate", "header": "On-Time %", "type": "percentage", "width": 100}]',
true),

-- Production Reports
('PRODUCTION_OUTPUT', 'Production Output', 'Work order completion by period', 'PRODUCTION',
'SELECT 
    wo.wo_number,
    p.name as product_name,
    wo.planned_quantity,
    wo.completed_quantity,
    ROUND(wo.completed_quantity * 100.0 / NULLIF(wo.planned_quantity, 0), 1) as completion_rate,
    wo.status,
    wo.planned_start_date,
    wo.actual_end_date
FROM work_orders wo
JOIN products p ON wo.product_id = p.id
WHERE wo.created_at BETWEEN {{.from_date}} AND {{.to_date}}
ORDER BY wo.created_at DESC',
'[{"name": "from_date", "type": "date", "required": true, "label": "From Date"},
  {"name": "to_date", "type": "date", "required": true, "label": "To Date"}]',
'[{"field": "wo_number", "header": "WO Number", "width": 150},
  {"field": "product_name", "header": "Product", "width": 200},
  {"field": "planned_quantity", "header": "Planned", "type": "number", "width": 100},
  {"field": "completed_quantity", "header": "Completed", "type": "number", "width": 100},
  {"field": "completion_rate", "header": "Completion %", "type": "percentage", "width": 100},
  {"field": "status", "header": "Status", "width": 120},
  {"field": "planned_start_date", "header": "Planned Start", "type": "date", "width": 120},
  {"field": "actual_end_date", "header": "Actual End", "type": "date", "width": 120}]',
true),

('QC_SUMMARY', 'QC Summary Report', 'Quality check pass/fail rates', 'PRODUCTION',
'SELECT 
    qc.qc_type,
    COUNT(*) as total_checks,
    SUM(CASE WHEN qc.overall_result = ''PASS'' THEN 1 ELSE 0 END) as passed,
    SUM(CASE WHEN qc.overall_result = ''FAIL'' THEN 1 ELSE 0 END) as failed,
    ROUND(SUM(CASE WHEN qc.overall_result = ''PASS'' THEN 1 ELSE 0 END) * 100.0 / COUNT(*), 1) as pass_rate
FROM qc_records qc
WHERE qc.qc_date BETWEEN {{.from_date}} AND {{.to_date}}
GROUP BY qc.qc_type
ORDER BY qc.qc_type',
'[{"name": "from_date", "type": "date", "required": true, "label": "From Date"},
  {"name": "to_date", "type": "date", "required": true, "label": "To Date"}]',
'[{"field": "qc_type", "header": "QC Type", "width": 150},
  {"field": "total_checks", "header": "Total Checks", "type": "number", "width": 120},
  {"field": "passed", "header": "Passed", "type": "number", "width": 100},
  {"field": "failed", "header": "Failed", "type": "number", "width": 100},
  {"field": "pass_rate", "header": "Pass Rate %", "type": "percentage", "width": 100}]',
true),

-- Sales Reports
('SALES_SUMMARY', 'Sales Summary', 'Sales by period and customer', 'SALES',
'SELECT 
    so.order_number,
    c.name as customer_name,
    so.order_date,
    so.total_amount,
    so.currency,
    so.status,
    so.payment_status
FROM sales_orders so
JOIN customers c ON so.customer_id = c.id
WHERE so.order_date BETWEEN {{.from_date}} AND {{.to_date}}
ORDER BY so.order_date DESC',
'[{"name": "from_date", "type": "date", "required": true, "label": "From Date"},
  {"name": "to_date", "type": "date", "required": true, "label": "To Date"}]',
'[{"field": "order_number", "header": "Order Number", "width": 150},
  {"field": "customer_name", "header": "Customer", "width": 200},
  {"field": "order_date", "header": "Order Date", "type": "date", "width": 120},
  {"field": "total_amount", "header": "Amount", "type": "currency", "width": 120},
  {"field": "currency", "header": "Currency", "width": 80},
  {"field": "status", "header": "Status", "width": 120},
  {"field": "payment_status", "header": "Payment", "width": 120}]',
true),

('TOP_PRODUCTS', 'Top Selling Products', 'Best selling products by quantity or value', 'SALES',
'SELECT 
    p.code as product_code,
    p.name as product_name,
    pc.name as category,
    SUM(soi.quantity) as total_quantity,
    SUM(soi.line_total) as total_value,
    COUNT(DISTINCT so.id) as order_count
FROM sales_order_items soi
JOIN sales_orders so ON soi.sales_order_id = so.id
JOIN products p ON soi.product_id = p.id
LEFT JOIN product_categories pc ON p.category_id = pc.id
WHERE so.order_date BETWEEN {{.from_date}} AND {{.to_date}}
  AND so.status != ''CANCELLED''
GROUP BY p.id, p.code, p.name, pc.name
ORDER BY total_value DESC
LIMIT {{.limit}}',
'[{"name": "from_date", "type": "date", "required": true, "label": "From Date"},
  {"name": "to_date", "type": "date", "required": true, "label": "To Date"},
  {"name": "limit", "type": "integer", "required": false, "default": 20, "label": "Top N"}]',
'[{"field": "product_code", "header": "Product Code", "width": 120},
  {"field": "product_name", "header": "Product Name", "width": 200},
  {"field": "category", "header": "Category", "width": 150},
  {"field": "total_quantity", "header": "Quantity", "type": "number", "width": 100},
  {"field": "total_value", "header": "Total Value", "type": "currency", "width": 120},
  {"field": "order_count", "header": "Orders", "type": "number", "width": 80}]',
true),

('LOW_STOCK_ITEMS', 'Low Stock Items', 'Items below reorder point', 'INVENTORY',
'SELECT 
    m.code as material_code,
    m.name as material_name,
    mc.name as category,
    COALESCE(SUM(s.quantity), 0) as current_stock,
    m.reorder_point,
    m.reorder_point - COALESCE(SUM(s.quantity), 0) as shortage,
    m.uom
FROM materials m
LEFT JOIN material_categories mc ON m.category_id = mc.id
LEFT JOIN stock_on_hand s ON m.id = s.material_id
WHERE m.is_active = true
  AND m.reorder_point > 0
GROUP BY m.id, m.code, m.name, mc.name, m.reorder_point, m.uom
HAVING COALESCE(SUM(s.quantity), 0) < m.reorder_point
ORDER BY (m.reorder_point - COALESCE(SUM(s.quantity), 0)) DESC',
'[]',
'[{"field": "material_code", "header": "Material Code", "width": 120},
  {"field": "material_name", "header": "Material Name", "width": 200},
  {"field": "category", "header": "Category", "width": 150},
  {"field": "current_stock", "header": "Current Stock", "type": "number", "width": 120},
  {"field": "reorder_point", "header": "Reorder Point", "type": "number", "width": 120},
  {"field": "shortage", "header": "Shortage", "type": "number", "width": 100},
  {"field": "uom", "header": "UoM", "width": 80}]',
true);

-- Seed default dashboard
INSERT INTO dashboards (code, name, description, is_default, is_system, visibility) VALUES
('MAIN_DASHBOARD', 'Main Dashboard', 'Overview of key metrics', true, true, 'PUBLIC');

-- Get dashboard ID for widgets
DO $$
DECLARE
    dashboard_uuid UUID;
BEGIN
    SELECT id INTO dashboard_uuid FROM dashboards WHERE code = 'MAIN_DASHBOARD';
    
    -- Insert widgets
    INSERT INTO widgets (dashboard_id, widget_type, title, data_source, config, position_x, position_y, width, height) VALUES
    (dashboard_uuid, 'KPI', 'Total Stock Value', '/api/v1/stats/inventory', 
     '{"icon": "pi-box", "color": "blue", "format": "currency"}', 0, 0, 3, 1),
    (dashboard_uuid, 'KPI', 'Pending Orders', '/api/v1/stats/procurement', 
     '{"icon": "pi-shopping-cart", "color": "orange", "field": "pending_orders"}', 3, 0, 3, 1),
    (dashboard_uuid, 'KPI', 'Expiring Soon', '/api/v1/stats/inventory', 
     '{"icon": "pi-clock", "color": "red", "field": "expiring_30_days"}', 6, 0, 3, 1),
    (dashboard_uuid, 'KPI', 'Active Work Orders', '/api/v1/stats/production', 
     '{"icon": "pi-cog", "color": "green", "field": "active_work_orders"}', 9, 0, 3, 1),
    (dashboard_uuid, 'PIE_CHART', 'Stock by Category', '/api/v1/stats/inventory/by-category', 
     '{"showLegend": true}', 0, 1, 6, 3),
    (dashboard_uuid, 'BAR_CHART', 'Monthly Sales', '/api/v1/stats/sales/monthly', 
     '{"xField": "month", "yField": "total", "color": "#4CAF50"}', 6, 1, 6, 3),
    (dashboard_uuid, 'TABLE', 'Low Stock Items', '/api/v1/reports/LOW_STOCK_ITEMS/preview', 
     '{"limit": 10}', 0, 4, 6, 3),
    (dashboard_uuid, 'TABLE', 'Expiring Lots', '/api/v1/reports/EXPIRY_REPORT/preview', 
     '{"limit": 10, "days": 30}', 6, 4, 6, 3);
END $$;
