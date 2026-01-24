-- Remove seed data
DELETE FROM locations WHERE zone_id IN (SELECT id FROM zones WHERE warehouse_id IN (SELECT id FROM warehouses WHERE code IN ('WH-MAIN', 'WH-COLD', 'WH-FG')));
DELETE FROM zones WHERE warehouse_id IN (SELECT id FROM warehouses WHERE code IN ('WH-MAIN', 'WH-COLD', 'WH-FG'));
DELETE FROM warehouses WHERE code IN ('WH-MAIN', 'WH-COLD', 'WH-FG');
