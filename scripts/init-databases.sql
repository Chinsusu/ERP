-- Create all databases for ERP microservices
CREATE DATABASE auth_db;
CREATE DATABASE user_db;
CREATE DATABASE master_data_db;
CREATE DATABASE supplier_db;
CREATE DATABASE procurement_db;
CREATE DATABASE wms_db;
CREATE DATABASE manufacturing_db;
CREATE DATABASE sales_db;
CREATE DATABASE marketing_db;
CREATE DATABASE finance_db;
CREATE DATABASE reporting_db;
CREATE DATABASE notification_db;
CREATE DATABASE ai_db;

-- Create service-specific users (optional, for better security)
-- Uncomment and customize if you want separate users per service

-- CREATE USER auth_service WITH PASSWORD 'auth_password';
-- GRANT ALL PRIVILEGES ON DATABASE auth_db TO auth_service;

-- CREATE USER user_service WITH PASSWORD 'user_password';
-- GRANT ALL PRIVILEGES ON DATABASE user_db TO user_service;

-- Add more users as needed for other services
