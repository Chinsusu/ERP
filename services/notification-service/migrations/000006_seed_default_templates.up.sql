-- Seed default notification templates
INSERT INTO notification_templates (template_code, name, notification_type, subject_template, body_template, variables) VALUES

-- Password Reset Email
('PASSWORD_RESET', 'Password Reset Request', 'EMAIL',
'Reset Your Password - ERP Cosmetics',
'<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
<h2 style="color: #2563eb;">Password Reset Request</h2>
<p>Hi {{.UserName}},</p>
<p>You requested to reset your password for ERP Cosmetics System.</p>
<p>Click the button below to reset your password:</p>
<p style="text-align: center; margin: 30px 0;">
<a href="{{.ResetLink}}" style="background-color: #2563eb; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px;">Reset Password</a>
</p>
<p style="color: #666; font-size: 14px;">This link will expire in 1 hour.</p>
<p style="color: #666; font-size: 14px;">If you didn''t request this, please ignore this email.</p>
<hr style="border: 1px solid #eee; margin: 20px 0;">
<p style="color: #999; font-size: 12px;">ERP Cosmetics Team</p>
</div>
</body>
</html>',
'["UserName", "ResetLink"]'),

-- Account Locked Alert
('ACCOUNT_LOCKED', 'Account Locked Alert', 'EMAIL',
'Account Locked - ERP Cosmetics',
'<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
<h2 style="color: #dc2626;">Account Locked</h2>
<p>Hi {{.UserName}},</p>
<p>Your account has been locked due to multiple failed login attempts.</p>
<p><strong>Lock Duration:</strong> {{.LockDuration}}</p>
<p><strong>Locked At:</strong> {{.LockedAt}}</p>
<p>If this wasn''t you, please contact the administrator immediately.</p>
<hr style="border: 1px solid #eee; margin: 20px 0;">
<p style="color: #999; font-size: 12px;">ERP Cosmetics Security Team</p>
</div>
</body>
</html>',
'["UserName", "LockDuration", "LockedAt"]'),

-- Low Stock Alert
('LOW_STOCK_ALERT', 'Low Stock Alert', 'EMAIL',
'Low Stock Alert: {{.MaterialName}} - ERP Cosmetics',
'<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
<h2 style="color: #f59e0b;">⚠️ Low Stock Alert</h2>
<table style="width: 100%; border-collapse: collapse; margin: 20px 0;">
<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Material:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee;">{{.MaterialCode}} - {{.MaterialName}}</td></tr>
<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Current Stock:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee;">{{.CurrentQuantity}} {{.UoM}}</td></tr>
<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Reorder Point:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee;">{{.ReorderPoint}} {{.UoM}}</td></tr>
<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Shortage:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee; color: #dc2626;">{{.Shortage}} {{.UoM}}</td></tr>
</table>
<p><strong>Action Required:</strong> Please create a Purchase Requisition.</p>
<p style="text-align: center; margin: 30px 0;">
<a href="{{.StockLink}}" style="background-color: #2563eb; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px;">View Stock</a>
</p>
</div>
</body>
</html>',
'["MaterialCode", "MaterialName", "CurrentQuantity", "ReorderPoint", "Shortage", "UoM", "StockLink"]'),

-- Lot Expiring Alert
('LOT_EXPIRING_ALERT', 'Lot Expiring Soon Alert', 'EMAIL',
'Lot Expiring Soon: {{.LotNumber}} - ERP Cosmetics',
'<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
<h2 style="color: #f59e0b;">⏰ Lot Expiring Soon</h2>
<table style="width: 100%; border-collapse: collapse; margin: 20px 0;">
<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Lot Number:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee;">{{.LotNumber}}</td></tr>
<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Material:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee;">{{.MaterialCode}} - {{.MaterialName}}</td></tr>
<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Expiry Date:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee; color: #dc2626;">{{.ExpiryDate}}</td></tr>
<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Days Until Expiry:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee; color: #f59e0b;">{{.DaysRemaining}} days</td></tr>
<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Quantity:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee;">{{.Quantity}} {{.UoM}}</td></tr>
<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Location:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee;">{{.Location}}</td></tr>
</table>
<p><strong>Action Required:</strong> Please prioritize using this lot (FEFO) or plan disposal.</p>
<p style="text-align: center; margin: 30px 0;">
<a href="{{.LotLink}}" style="background-color: #2563eb; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px;">View Lot Details</a>
</p>
</div>
</body>
</html>',
'["LotNumber", "MaterialCode", "MaterialName", "ExpiryDate", "DaysRemaining", "Quantity", "UoM", "Location", "LotLink"]'),

-- Certificate Expiring Alert
('CERT_EXPIRING_ALERT', 'Certificate Expiring Soon', 'EMAIL',
'Certificate Expiring: {{.SupplierName}} - ERP Cosmetics',
'<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
<h2 style="color: #f59e0b;">📜 Certificate Expiring Soon</h2>
<table style="width: 100%; border-collapse: collapse; margin: 20px 0;">
<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Supplier:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee;">{{.SupplierName}}</td></tr>
<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Certificate Type:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee;">{{.CertificateType}}</td></tr>
<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Certificate Number:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee;">{{.CertificateNumber}}</td></tr>
<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Expiry Date:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee; color: #dc2626;">{{.ExpiryDate}}</td></tr>
<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Days Until Expiry:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee; color: #f59e0b;">{{.DaysRemaining}} days</td></tr>
</table>
<p><strong>Action Required:</strong> Request updated certificate from supplier.</p>
<p style="text-align: center; margin: 30px 0;">
<a href="{{.CertificateLink}}" style="background-color: #2563eb; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px;">View Certificate</a>
</p>
</div>
</body>
</html>',
'["SupplierName", "CertificateType", "CertificateNumber", "ExpiryDate", "DaysRemaining", "CertificateLink"]'),

-- PR Pending Approval (In-App)
('PR_PENDING_APPROVAL', 'PR Pending Approval', 'IN_APP',
'Purchase Requisition Pending Approval',
'PR {{.PRNumber}} from {{.RequesterName}} ({{.DepartmentName}}) requires your approval. Total: {{.TotalAmount}}',
'["PRNumber", "RequesterName", "DepartmentName", "TotalAmount", "PRLink"]'),

-- PO Created Notification
('PO_CREATED', 'Purchase Order Created', 'IN_APP',
'Purchase Order Created',
'PO {{.PONumber}} has been created for {{.SupplierName}}. Total: {{.TotalAmount}}',
'["PONumber", "SupplierName", "TotalAmount", "POLink"]'),

-- QC Failed Alert
('QC_FAILED', 'Quality Check Failed', 'BOTH',
'QC Failed: {{.WorkOrderNumber}} - ERP Cosmetics',
'<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
<h2 style="color: #dc2626;">❌ Quality Check Failed</h2>
<table style="width: 100%; border-collapse: collapse; margin: 20px 0;">
<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Work Order:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee;">{{.WorkOrderNumber}}</td></tr>
<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Product:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee;">{{.ProductName}}</td></tr>
<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Batch Number:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee;">{{.BatchNumber}}</td></tr>
<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>QC Type:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee;">{{.QCType}}</td></tr>
<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Failed Items:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee; color: #dc2626;">{{.FailedItems}}</td></tr>
<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Inspector:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee;">{{.InspectorName}}</td></tr>
</table>
<p><strong>Action Required:</strong> Review QC results and create NCR if necessary.</p>
<p style="text-align: center; margin: 30px 0;">
<a href="{{.QCLink}}" style="background-color: #dc2626; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px;">View QC Report</a>
</p>
</div>
</body>
</html>',
'["WorkOrderNumber", "ProductName", "BatchNumber", "QCType", "FailedItems", "InspectorName", "QCLink"]'),

-- Order Confirmation
('ORDER_CONFIRMATION', 'Order Confirmation', 'EMAIL',
'Order Confirmation: {{.OrderNumber}} - ERP Cosmetics',
'<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
<h2 style="color: #16a34a;">✅ Order Confirmed</h2>
<p>Dear {{.CustomerName}},</p>
<p>Thank you for your order. Your order has been confirmed.</p>
<table style="width: 100%; border-collapse: collapse; margin: 20px 0;">
<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Order Number:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee;">{{.OrderNumber}}</td></tr>
<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Order Date:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee;">{{.OrderDate}}</td></tr>
<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Total Amount:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee;">{{.TotalAmount}}</td></tr>
<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Estimated Delivery:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee;">{{.DeliveryDate}}</td></tr>
</table>
<p>We will notify you when your order ships.</p>
<hr style="border: 1px solid #eee; margin: 20px 0;">
<p style="color: #999; font-size: 12px;">ERP Cosmetics Team</p>
</div>
</body>
</html>',
'["CustomerName", "OrderNumber", "OrderDate", "TotalAmount", "DeliveryDate"]');

-- Seed default alert rules
INSERT INTO alert_rules (rule_code, name, description, rule_type, conditions, notification_type, recipients) VALUES

('STOCK_LOW_ALERT', 'Low Stock Alert', 'Notify when stock falls below reorder point', 'STOCK_LOW',
'{"threshold_type": "REORDER_POINT", "check_interval": "1h"}',
'BOTH',
'[{"role": "PROCUREMENT_MANAGER"}, {"role": "WAREHOUSE_MANAGER"}]'),

('LOT_EXPIRY_30_DAYS', 'Lot Expiring in 30 Days', 'Notify when lot expires within 30 days', 'LOT_EXPIRY',
'{"days_before_expiry": 30, "check_interval": "daily"}',
'EMAIL',
'[{"role": "WAREHOUSE_MANAGER"}, {"role": "QC_MANAGER"}]'),

('LOT_EXPIRY_7_DAYS', 'Lot Expiring in 7 Days', 'Urgent: Lot expires within 7 days', 'LOT_EXPIRY',
'{"days_before_expiry": 7, "check_interval": "daily"}',
'BOTH',
'[{"role": "WAREHOUSE_MANAGER"}, {"role": "PRODUCTION_MANAGER"}]'),

('CERT_EXPIRY_90_DAYS', 'Certificate Expiring in 90 Days', 'Supplier certificate expires within 90 days', 'CERT_EXPIRY',
'{"days_before_expiry": 90, "check_interval": "daily"}',
'EMAIL',
'[{"role": "PROCUREMENT_MANAGER"}]'),

('CERT_EXPIRY_30_DAYS', 'Certificate Expiring in 30 Days', 'Urgent: Supplier certificate expires within 30 days', 'CERT_EXPIRY',
'{"days_before_expiry": 30, "check_interval": "daily"}',
'BOTH',
'[{"role": "PROCUREMENT_MANAGER"}, {"role": "COMPLIANCE_OFFICER"}]'),

('APPROVAL_PENDING_24H', 'Approval Pending > 24h', 'Notify when approval pending for more than 24 hours', 'APPROVAL_PENDING',
'{"hours_pending": 24, "check_interval": "1h"}',
'IN_APP',
'[{"role": "ADMIN"}]');
