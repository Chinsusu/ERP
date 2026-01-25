-- Seed default file categories
INSERT INTO file_categories (code, name, description, allowed_extensions, max_file_size, storage_bucket) VALUES

('DOCUMENT', 'Documents', 'General documents', 
 '["pdf", "doc", "docx", "xls", "xlsx", "txt"]', 
 10485760, -- 10 MB
 'documents'),

('IMAGE', 'Images', 'Image files', 
 '["jpg", "jpeg", "png", "gif", "webp", "svg"]', 
 5242880, -- 5 MB
 'images'),

('CERTIFICATE', 'Certificates', 'Supplier and product certificates', 
 '["pdf"]', 
 5242880, -- 5 MB
 'certificates'),

('CONTRACT', 'Contracts', 'Supplier contracts and agreements', 
 '["pdf", "doc", "docx"]', 
 20971520, -- 20 MB
 'contracts'),

('REPORT', 'Reports', 'Generated reports and exports', 
 '["pdf", "xlsx", "csv"]', 
 52428800, -- 50 MB
 'reports'),

('AVATAR', 'Avatars', 'User profile pictures', 
 '["jpg", "jpeg", "png"]', 
 2097152, -- 2 MB
 'avatars'),

('PRODUCT_IMAGE', 'Product Images', 'Product catalog images', 
 '["jpg", "jpeg", "png", "webp"]', 
 5242880, -- 5 MB
 'products'),

('QC_PHOTO', 'QC Photos', 'Quality control inspection photos', 
 '["jpg", "jpeg", "png"]', 
 10485760, -- 10 MB
 'qc-photos'),

('SIGNATURE', 'Signatures', 'Digital signatures', 
 '["png", "jpg"]', 
 1048576, -- 1 MB
 'signatures'),

('ATTACHMENT', 'Attachments', 'General attachments', 
 '["pdf", "doc", "docx", "xls", "xlsx", "jpg", "png", "zip"]', 
 20971520, -- 20 MB
 'attachments');
