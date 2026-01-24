-- Migration: Seed default data
-- Default categories, units of measure, and conversions

-- ============================
-- Units of Measure
-- ============================

-- Base units (weight)
INSERT INTO units_of_measure (id, code, name, name_en, symbol, uom_type, is_base_unit, conversion_factor) VALUES
('a0000000-0000-0000-0000-000000000001', 'KG', 'Kilogram', 'Kilogram', 'kg', 'WEIGHT', true, 1),
('a0000000-0000-0000-0000-000000000002', 'G', 'Gram', 'Gram', 'g', 'WEIGHT', false, 0.001),
('a0000000-0000-0000-0000-000000000003', 'MG', 'Miligram', 'Milligram', 'mg', 'WEIGHT', false, 0.000001);

-- Update base_unit_id for weight units
UPDATE units_of_measure SET base_unit_id = 'a0000000-0000-0000-0000-000000000001' WHERE code IN ('G', 'MG');

-- Base units (volume)
INSERT INTO units_of_measure (id, code, name, name_en, symbol, uom_type, is_base_unit, conversion_factor) VALUES
('a0000000-0000-0000-0000-000000000011', 'L', 'Lít', 'Liter', 'L', 'VOLUME', true, 1),
('a0000000-0000-0000-0000-000000000012', 'ML', 'Mililit', 'Milliliter', 'mL', 'VOLUME', false, 0.001);

-- Update base_unit_id for volume units
UPDATE units_of_measure SET base_unit_id = 'a0000000-0000-0000-0000-000000000011' WHERE code = 'ML';

-- Quantity units
INSERT INTO units_of_measure (id, code, name, name_en, symbol, uom_type, is_base_unit, conversion_factor) VALUES
('a0000000-0000-0000-0000-000000000021', 'PCS', 'Cái', 'Piece', 'pcs', 'QUANTITY', true, 1),
('a0000000-0000-0000-0000-000000000022', 'BOX', 'Hộp', 'Box', 'box', 'QUANTITY', true, 1),
('a0000000-0000-0000-0000-000000000023', 'SET', 'Bộ', 'Set', 'set', 'QUANTITY', true, 1),
('a0000000-0000-0000-0000-000000000024', 'BTL', 'Chai', 'Bottle', 'btl', 'QUANTITY', true, 1),
('a0000000-0000-0000-0000-000000000025', 'TUBE', 'Tuýp', 'Tube', 'tube', 'QUANTITY', true, 1),
('a0000000-0000-0000-0000-000000000026', 'JAR', 'Hũ', 'Jar', 'jar', 'QUANTITY', true, 1);

-- ============================
-- Unit Conversions
-- ============================

-- Weight conversions (bidirectional)
INSERT INTO unit_conversions (from_unit_id, to_unit_id, conversion_factor) VALUES
('a0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000002', 1000),    -- KG to G
('a0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001', 0.001),   -- G to KG
('a0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000003', 1000),    -- G to MG
('a0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000002', 0.001),   -- MG to G
('a0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000003', 1000000), -- KG to MG
('a0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001', 0.000001); -- MG to KG

-- Volume conversions (bidirectional)
INSERT INTO unit_conversions (from_unit_id, to_unit_id, conversion_factor) VALUES
('a0000000-0000-0000-0000-000000000011', 'a0000000-0000-0000-0000-000000000012', 1000),   -- L to ML
('a0000000-0000-0000-0000-000000000012', 'a0000000-0000-0000-0000-000000000011', 0.001);  -- ML to L

-- ============================
-- Material Categories
-- ============================

-- Root categories for materials
INSERT INTO categories (id, code, name, name_en, category_type, path, level, sort_order) VALUES
('b0000000-0000-0000-0000-000000000001', 'RAW-MATERIALS', 'Nguyên liệu thô', 'Raw Materials', 'MATERIAL', '/RAW-MATERIALS/', 0, 1),
('b0000000-0000-0000-0000-000000000002', 'PACKAGING', 'Bao bì', 'Packaging', 'MATERIAL', '/PACKAGING/', 0, 2),
('b0000000-0000-0000-0000-000000000003', 'CONSUMABLES', 'Vật tư tiêu hao', 'Consumables', 'MATERIAL', '/CONSUMABLES/', 0, 3);

-- Sub-categories for Raw Materials
INSERT INTO categories (id, code, name, name_en, category_type, parent_id, path, level, sort_order) VALUES
('b0000000-0000-0000-0000-000000000011', 'ACTIVE-ING', 'Hoạt chất', 'Active Ingredients', 'MATERIAL', 'b0000000-0000-0000-0000-000000000001', '/RAW-MATERIALS/ACTIVE-ING/', 1, 1),
('b0000000-0000-0000-0000-000000000012', 'EMOLLIENTS', 'Chất làm mềm', 'Emollients', 'MATERIAL', 'b0000000-0000-0000-0000-000000000001', '/RAW-MATERIALS/EMOLLIENTS/', 1, 2),
('b0000000-0000-0000-0000-000000000013', 'EMULSIFIERS', 'Chất nhũ hóa', 'Emulsifiers', 'MATERIAL', 'b0000000-0000-0000-0000-000000000001', '/RAW-MATERIALS/EMULSIFIERS/', 1, 3),
('b0000000-0000-0000-0000-000000000014', 'PRESERVATIVES', 'Chất bảo quản', 'Preservatives', 'MATERIAL', 'b0000000-0000-0000-0000-000000000001', '/RAW-MATERIALS/PRESERVATIVES/', 1, 4),
('b0000000-0000-0000-0000-000000000015', 'FRAGRANCES', 'Hương liệu', 'Fragrances', 'MATERIAL', 'b0000000-0000-0000-0000-000000000001', '/RAW-MATERIALS/FRAGRANCES/', 1, 5),
('b0000000-0000-0000-0000-000000000016', 'COLORANTS', 'Chất tạo màu', 'Colorants', 'MATERIAL', 'b0000000-0000-0000-0000-000000000001', '/RAW-MATERIALS/COLORANTS/', 1, 6),
('b0000000-0000-0000-0000-000000000017', 'THICKENERS', 'Chất làm đặc', 'Thickeners', 'MATERIAL', 'b0000000-0000-0000-0000-000000000001', '/RAW-MATERIALS/THICKENERS/', 1, 7),
('b0000000-0000-0000-0000-000000000018', 'VITAMINS', 'Vitamin', 'Vitamins', 'MATERIAL', 'b0000000-0000-0000-0000-000000000011', '/RAW-MATERIALS/ACTIVE-ING/VITAMINS/', 2, 1);

-- Sub-categories for Packaging
INSERT INTO categories (id, code, name, name_en, category_type, parent_id, path, level, sort_order) VALUES
('b0000000-0000-0000-0000-000000000021', 'PRIMARY-PKG', 'Bao bì sơ cấp', 'Primary Packaging', 'MATERIAL', 'b0000000-0000-0000-0000-000000000002', '/PACKAGING/PRIMARY-PKG/', 1, 1),
('b0000000-0000-0000-0000-000000000022', 'SECONDARY-PKG', 'Bao bì thứ cấp', 'Secondary Packaging', 'MATERIAL', 'b0000000-0000-0000-0000-000000000002', '/PACKAGING/SECONDARY-PKG/', 1, 2),
('b0000000-0000-0000-0000-000000000023', 'LABELS', 'Nhãn mác', 'Labels', 'MATERIAL', 'b0000000-0000-0000-0000-000000000002', '/PACKAGING/LABELS/', 1, 3);

-- ============================
-- Product Categories
-- ============================

INSERT INTO categories (id, code, name, name_en, category_type, path, level, sort_order) VALUES
('c0000000-0000-0000-0000-000000000001', 'SKINCARE', 'Chăm sóc da', 'Skincare', 'PRODUCT', '/SKINCARE/', 0, 1),
('c0000000-0000-0000-0000-000000000002', 'HAIRCARE', 'Chăm sóc tóc', 'Haircare', 'PRODUCT', '/HAIRCARE/', 0, 2),
('c0000000-0000-0000-0000-000000000003', 'BODYCARE', 'Chăm sóc cơ thể', 'Bodycare', 'PRODUCT', '/BODYCARE/', 0, 3);

-- Sub-categories for Skincare
INSERT INTO categories (id, code, name, name_en, category_type, parent_id, path, level, sort_order) VALUES
('c0000000-0000-0000-0000-000000000011', 'SERUM', 'Tinh chất', 'Serum', 'PRODUCT', 'c0000000-0000-0000-0000-000000000001', '/SKINCARE/SERUM/', 1, 1),
('c0000000-0000-0000-0000-000000000012', 'CREAM', 'Kem dưỡng', 'Cream', 'PRODUCT', 'c0000000-0000-0000-0000-000000000001', '/SKINCARE/CREAM/', 1, 2),
('c0000000-0000-0000-0000-000000000013', 'CLEANSER', 'Sữa rửa mặt', 'Cleanser', 'PRODUCT', 'c0000000-0000-0000-0000-000000000001', '/SKINCARE/CLEANSER/', 1, 3),
('c0000000-0000-0000-0000-000000000014', 'TONER', 'Nước hoa hồng', 'Toner', 'PRODUCT', 'c0000000-0000-0000-0000-000000000001', '/SKINCARE/TONER/', 1, 4),
('c0000000-0000-0000-0000-000000000015', 'MASK', 'Mặt nạ', 'Mask', 'PRODUCT', 'c0000000-0000-0000-0000-000000000001', '/SKINCARE/MASK/', 1, 5),
('c0000000-0000-0000-0000-000000000016', 'SUNSCREEN', 'Kem chống nắng', 'Sunscreen', 'PRODUCT', 'c0000000-0000-0000-0000-000000000001', '/SKINCARE/SUNSCREEN/', 1, 6);

-- Sub-categories for Haircare
INSERT INTO categories (id, code, name, name_en, category_type, parent_id, path, level, sort_order) VALUES
('c0000000-0000-0000-0000-000000000021', 'SHAMPOO', 'Dầu gội', 'Shampoo', 'PRODUCT', 'c0000000-0000-0000-0000-000000000002', '/HAIRCARE/SHAMPOO/', 1, 1),
('c0000000-0000-0000-0000-000000000022', 'CONDITIONER', 'Dầu xả', 'Conditioner', 'PRODUCT', 'c0000000-0000-0000-0000-000000000002', '/HAIRCARE/CONDITIONER/', 1, 2),
('c0000000-0000-0000-0000-000000000023', 'HAIR-TREATMENT', 'Dưỡng tóc', 'Hair Treatment', 'PRODUCT', 'c0000000-0000-0000-0000-000000000002', '/HAIRCARE/HAIR-TREATMENT/', 1, 3);

-- Sub-categories for Bodycare
INSERT INTO categories (id, code, name, name_en, category_type, parent_id, path, level, sort_order) VALUES
('c0000000-0000-0000-0000-000000000031', 'BODY-LOTION', 'Dưỡng thể', 'Body Lotion', 'PRODUCT', 'c0000000-0000-0000-0000-000000000003', '/BODYCARE/BODY-LOTION/', 1, 1),
('c0000000-0000-0000-0000-000000000032', 'SHOWER-GEL', 'Sữa tắm', 'Shower Gel', 'PRODUCT', 'c0000000-0000-0000-0000-000000000003', '/BODYCARE/SHOWER-GEL/', 1, 2),
('c0000000-0000-0000-0000-000000000033', 'HAND-CREAM', 'Kem tay', 'Hand Cream', 'PRODUCT', 'c0000000-0000-0000-0000-000000000003', '/BODYCARE/HAND-CREAM/', 1, 3);
