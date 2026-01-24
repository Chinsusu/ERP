-- Seed Data for Warehouses, Zones, and Locations

-- Main Warehouse
INSERT INTO warehouses (id, code, name, warehouse_type, address, is_active) VALUES
    ('a1b2c3d4-1111-1111-1111-111111111111', 'WH-MAIN', 'Main Warehouse', 'MAIN', '123 Industrial Park, District 9, HCMC', true),
    ('a1b2c3d4-2222-2222-2222-222222222222', 'WH-COLD', 'Cold Storage Warehouse', 'COLD_STORAGE', '123 Industrial Park, District 9, HCMC', true),
    ('a1b2c3d4-3333-3333-3333-333333333333', 'WH-FG', 'Finished Goods Warehouse', 'FINISHED_GOODS', '125 Industrial Park, District 9, HCMC', true);

-- Zones for Main Warehouse
INSERT INTO zones (id, warehouse_id, code, name, zone_type, is_active) VALUES
    ('b1b2c3d4-1111-1111-1111-111111111111', 'a1b2c3d4-1111-1111-1111-111111111111', 'ZONE-RECV', 'Receiving Zone', 'RECEIVING', true),
    ('b1b2c3d4-2222-2222-2222-222222222222', 'a1b2c3d4-1111-1111-1111-111111111111', 'ZONE-QUAR', 'Quarantine Zone', 'QUARANTINE', true),
    ('b1b2c3d4-3333-3333-3333-333333333333', 'a1b2c3d4-1111-1111-1111-111111111111', 'ZONE-STOR', 'Storage Zone', 'STORAGE', true),
    ('b1b2c3d4-4444-4444-4444-444444444444', 'a1b2c3d4-1111-1111-1111-111111111111', 'ZONE-PICK', 'Picking Zone', 'PICKING', true),
    ('b1b2c3d4-5555-5555-5555-555555555555', 'a1b2c3d4-1111-1111-1111-111111111111', 'ZONE-SHIP', 'Shipping Zone', 'SHIPPING', true);

-- Cold Zone for Cold Storage Warehouse
INSERT INTO zones (id, warehouse_id, code, name, zone_type, temperature_min, temperature_max, is_active) VALUES
    ('b1b2c3d4-6666-6666-6666-666666666666', 'a1b2c3d4-2222-2222-2222-222222222222', 'ZONE-COLD-01', 'Cold Storage Zone 1', 'COLD', 2.00, 8.00, true),
    ('b1b2c3d4-7777-7777-7777-777777777777', 'a1b2c3d4-2222-2222-2222-222222222222', 'ZONE-COLD-02', 'Cold Storage Zone 2', 'COLD', 2.00, 8.00, true);

-- Zones for Finished Goods Warehouse
INSERT INTO zones (id, warehouse_id, code, name, zone_type, is_active) VALUES
    ('b1b2c3d4-8888-8888-8888-888888888888', 'a1b2c3d4-3333-3333-3333-333333333333', 'ZONE-FG-STOR', 'FG Storage Zone', 'STORAGE', true),
    ('b1b2c3d4-9999-9999-9999-999999999999', 'a1b2c3d4-3333-3333-3333-333333333333', 'ZONE-FG-SHIP', 'FG Shipping Zone', 'SHIPPING', true);

-- Locations for Quarantine Zone
INSERT INTO locations (id, zone_id, code, aisle, rack, shelf, bin, capacity, is_active) VALUES
    ('c1b2c3d4-1111-1111-1111-111111111111', 'b1b2c3d4-2222-2222-2222-222222222222', 'Q01-R01-S01', 'Q01', 'R01', 'S01', NULL, 1000.00, true),
    ('c1b2c3d4-2222-2222-2222-222222222222', 'b1b2c3d4-2222-2222-2222-222222222222', 'Q01-R01-S02', 'Q01', 'R01', 'S02', NULL, 1000.00, true),
    ('c1b2c3d4-3333-3333-3333-333333333333', 'b1b2c3d4-2222-2222-2222-222222222222', 'Q01-R02-S01', 'Q01', 'R02', 'S01', NULL, 1000.00, true);

-- Locations for Storage Zone
INSERT INTO locations (id, zone_id, code, aisle, rack, shelf, bin, capacity, is_active) VALUES
    ('c1b2c3d4-4444-4444-4444-444444444444', 'b1b2c3d4-3333-3333-3333-333333333333', 'A01-R01-S01', 'A01', 'R01', 'S01', NULL, 500.00, true),
    ('c1b2c3d4-5555-5555-5555-555555555555', 'b1b2c3d4-3333-3333-3333-333333333333', 'A01-R01-S02', 'A01', 'R01', 'S02', NULL, 500.00, true),
    ('c1b2c3d4-6666-6666-6666-666666666666', 'b1b2c3d4-3333-3333-3333-333333333333', 'A01-R01-S03', 'A01', 'R01', 'S03', NULL, 500.00, true),
    ('c1b2c3d4-7777-7777-7777-777777777777', 'b1b2c3d4-3333-3333-3333-333333333333', 'A01-R02-S01', 'A01', 'R02', 'S01', NULL, 500.00, true),
    ('c1b2c3d4-8888-8888-8888-888888888888', 'b1b2c3d4-3333-3333-3333-333333333333', 'A02-R01-S01', 'A02', 'R01', 'S01', NULL, 500.00, true),
    ('c1b2c3d4-9999-9999-9999-999999999999', 'b1b2c3d4-3333-3333-3333-333333333333', 'A02-R01-S02', 'A02', 'R01', 'S02', NULL, 500.00, true);

-- Locations for Cold Storage
INSERT INTO locations (id, zone_id, code, aisle, rack, shelf, bin, capacity, is_active) VALUES
    ('c1b2c3d4-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'b1b2c3d4-6666-6666-6666-666666666666', 'C01-R01-S01', 'C01', 'R01', 'S01', NULL, 200.00, true),
    ('c1b2c3d4-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'b1b2c3d4-6666-6666-6666-666666666666', 'C01-R01-S02', 'C01', 'R01', 'S02', NULL, 200.00, true),
    ('c1b2c3d4-cccc-cccc-cccc-cccccccccccc', 'b1b2c3d4-7777-7777-7777-777777777777', 'C02-R01-S01', 'C02', 'R01', 'S01', NULL, 200.00, true);

-- Locations for Finished Goods Storage
INSERT INTO locations (id, zone_id, code, aisle, rack, shelf, bin, capacity, is_active) VALUES
    ('c1b2c3d4-dddd-dddd-dddd-dddddddddddd', 'b1b2c3d4-8888-8888-8888-888888888888', 'FG01-R01-S01', 'FG01', 'R01', 'S01', NULL, 800.00, true),
    ('c1b2c3d4-eeee-eeee-eeee-eeeeeeeeeeee', 'b1b2c3d4-8888-8888-8888-888888888888', 'FG01-R01-S02', 'FG01', 'R01', 'S02', NULL, 800.00, true),
    ('c1b2c3d4-ffff-ffff-ffff-ffffffffffff', 'b1b2c3d4-8888-8888-8888-888888888888', 'FG02-R01-S01', 'FG02', 'R01', 'S01', NULL, 800.00, true);
