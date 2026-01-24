-- Rollback: Delete seed data

-- Delete categories (children first due to FK)
DELETE FROM categories WHERE level > 0;
DELETE FROM categories WHERE level = 0;

-- Delete unit conversions
DELETE FROM unit_conversions;

-- Delete units of measure
DELETE FROM units_of_measure;
