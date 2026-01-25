-- Notification templates table
CREATE TABLE IF NOT EXISTS notification_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(200) NOT NULL,
    notification_type VARCHAR(50) NOT NULL, -- EMAIL, IN_APP, SMS, PUSH
    subject_template TEXT,
    body_template TEXT NOT NULL, -- HTML for email, plain text for others
    variables JSONB DEFAULT '[]', -- Available variables for this template
    is_active BOOLEAN DEFAULT true,
    created_by UUID,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notification_templates_code ON notification_templates(template_code);
CREATE INDEX idx_notification_templates_type ON notification_templates(notification_type);
CREATE INDEX idx_notification_templates_active ON notification_templates(is_active);

COMMENT ON TABLE notification_templates IS 'Notification templates for emails, in-app, SMS, and push notifications';
COMMENT ON COLUMN notification_templates.template_code IS 'Unique code for template lookup, e.g., PASSWORD_RESET, LOW_STOCK_ALERT';
COMMENT ON COLUMN notification_templates.notification_type IS 'Type of notification: EMAIL, IN_APP, SMS, PUSH';
COMMENT ON COLUMN notification_templates.variables IS 'JSON array of available template variables';
