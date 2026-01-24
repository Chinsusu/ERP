-- Seed data for KOL tiers
INSERT INTO kol_tiers (id, code, name, description, min_followers, max_followers, auto_approve_samples, priority) VALUES
('c0000001-0001-0001-0001-000000000001', 'MEGA', 'Mega Influencer', 'More than 1 million followers', 1000000, NULL, true, 1),
('c0000001-0001-0001-0001-000000000002', 'MACRO', 'Macro Influencer', '100K to 1M followers', 100000, 1000000, false, 2),
('c0000001-0001-0001-0001-000000000003', 'MICRO', 'Micro Influencer', '10K to 100K followers', 10000, 100000, false, 3),
('c0000001-0001-0001-0001-000000000004', 'NANO', 'Nano Influencer', 'Less than 10K followers', 0, 10000, false, 4);

-- Seed sample KOLs
INSERT INTO kols (id, kol_code, name, email, phone, tier_id, category, instagram_handle, instagram_followers, youtube_channel, youtube_subscribers, tiktok_handle, tiktok_followers, avg_engagement_rate, niche, collaboration_rate, status) VALUES
('d0000001-0001-0001-0001-000000000001', 'KOL-0001', 'Beauty By Linh', 'contact@beautybylinh.com', '+84 901 234 567', 'c0000001-0001-0001-0001-000000000003', 'BEAUTY_BLOGGER', '@beautybylinh', 50000, 'Beauty By Linh', 25000, '@beautybylinh', 30000, 5.50, 'Skincare', 5000000, 'ACTIVE'),
('d0000001-0001-0001-0001-000000000002', 'KOL-0002', 'Makeup Queen VN', 'collab@makeupqueenvn.com', '+84 902 345 678', 'c0000001-0001-0001-0001-000000000002', 'INFLUENCER', '@makeupqueenvn', 500000, 'Makeup Queen VN', 200000, '@makeupqueenvn', 800000, 4.20, 'Makeup', 50000000, 'ACTIVE'),
('d0000001-0001-0001-0001-000000000003', 'KOL-0003', 'Dr. Skin Expert', 'drskinexpert@gmail.com', '+84 903 456 789', 'c0000001-0001-0001-0001-000000000002', 'EXPERT', '@drskinexpert', 150000, 'Dr Skin Expert', 80000, '@drskinexpert', 50000, 6.80, 'Skincare', 30000000, 'ACTIVE'),
('d0000001-0001-0001-0001-000000000004', 'KOL-0004', 'Natural Beauty Tips', 'hello@naturalbeautytips.vn', '+84 904 567 890', 'c0000001-0001-0001-0001-000000000003', 'BEAUTY_BLOGGER', '@naturalbeautytips', 75000, NULL, 0, '@naturalbeautytips', 100000, 7.20, 'Natural Cosmetics', 8000000, 'ACTIVE');

-- Seed sample campaign
INSERT INTO campaigns (id, campaign_code, name, description, campaign_type, start_date, end_date, budget, status, channels) VALUES
('e0000001-0001-0001-0001-000000000001', 'CAMP-SERUM-2024', 'Vitamin C Serum Launch', 'Launch campaign for new Vitamin C serum', 'PRODUCT_LAUNCH', '2024-03-01', '2024-03-31', 100000000, 'PLANNED', '["INSTAGRAM", "TIKTOK", "YOUTUBE"]'::jsonb);
