-- Remove seed data
DELETE FROM campaigns WHERE campaign_code = 'CAMP-SERUM-2024';
DELETE FROM kols WHERE kol_code IN ('KOL-0001', 'KOL-0002', 'KOL-0003', 'KOL-0004');
DELETE FROM kol_tiers WHERE code IN ('MEGA', 'MACRO', 'MICRO', 'NANO');
