-- Temperature Logs table (for cold storage monitoring)
CREATE TABLE IF NOT EXISTS temperature_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    zone_id UUID NOT NULL REFERENCES zones(id),
    temperature DECIMAL(5,2) NOT NULL,
    humidity DECIMAL(5,2),
    recorded_at TIMESTAMP NOT NULL,
    recorded_by UUID,
    is_alert BOOLEAN DEFAULT false,
    alert_type VARCHAR(20), -- HIGH_TEMP, LOW_TEMP
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_temp_logs_zone ON temperature_logs(zone_id);
CREATE INDEX idx_temp_logs_date ON temperature_logs(recorded_at);
CREATE INDEX idx_temp_logs_alert ON temperature_logs(is_alert);
