DELETE FROM widgets WHERE dashboard_id IN (SELECT id FROM dashboards WHERE code = 'MAIN_DASHBOARD');
DELETE FROM dashboards WHERE code = 'MAIN_DASHBOARD';
DELETE FROM report_definitions WHERE is_system = true;
