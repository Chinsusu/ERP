-- Alert rules configuration table
CREATE TABLE IF NOT EXISTS alert_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rule_code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    
    -- Rule type and conditions
    rule_type VARCHAR(50) NOT NULL, -- STOCK_LOW, LOT_EXPIRY, CERT_EXPIRY, APPROVAL_PENDING, TEMP_OUT_OF_RANGE
    conditions JSONB NOT NULL DEFAULT '{}', -- Rule conditions
    
    -- Template and recipients
    notification_template_id UUID REFERENCES notification_templates(id),
    notification_type VARCHAR(50) DEFAULT 'IN_APP', -- EMAIL, IN_APP, BOTH
    recipients JSONB DEFAULT '[]', -- Array of {role: string} or {user_id: uuid}
    
    -- Schedule
    check_interval VARCHAR(20) DEFAULT '1h', -- Cron-like interval: 1h, 30m, daily, etc.
    last_checked_at TIMESTAMP,
    next_check_at TIMESTAMP,
    
    -- Status
    is_active BOOLEAN DEFAULT true,
    
    created_by UUID,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_alert_rules_code ON alert_rules(rule_code);
CREATE INDEX idx_alert_rules_type ON alert_rules(rule_type);
CREATE INDEX idx_alert_rules_active ON alert_rules(is_active);
CREATE INDEX idx_alert_rules_next_check ON alert_rules(next_check_at) WHERE is_active = true;

COMMENT ON TABLE alert_rules IS 'Configurable alert rules for automated notifications';
COMMENT ON COLUMN alert_rules.rule_type IS 'Type: STOCK_LOW, LOT_EXPIRY, CERT_EXPIRY, APPROVAL_PENDING, TEMP_OUT_OF_RANGE';
COMMENT ON COLUMN alert_rules.conditions IS 'JSON object with rule-specific conditions';
COMMENT ON COLUMN alert_rules.recipients IS 'JSON array of recipients: [{role: "WAREHOUSE_MANAGER"}, {user_id: "uuid"}]';
