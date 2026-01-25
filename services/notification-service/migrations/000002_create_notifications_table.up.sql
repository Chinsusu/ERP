-- Notifications queue table (email/SMS outbound tracking)
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    notification_type VARCHAR(50) NOT NULL, -- EMAIL, SMS, PUSH
    template_id UUID REFERENCES notification_templates(id),
    
    -- Recipient information
    recipient_user_id UUID,
    recipient_email VARCHAR(255),
    recipient_phone VARCHAR(50),
    
    -- Content
    subject TEXT,
    body TEXT NOT NULL,
    
    -- Additional data
    metadata JSONB DEFAULT '{}', -- Event data, links, etc.
    
    -- Status tracking
    status VARCHAR(50) DEFAULT 'PENDING', -- PENDING, SENT, FAILED, RETRYING
    priority VARCHAR(20) DEFAULT 'NORMAL', -- HIGH, NORMAL, LOW
    
    -- Timestamps and retry info
    sent_at TIMESTAMP,
    failed_at TIMESTAMP,
    retry_count INT DEFAULT 0,
    max_retries INT DEFAULT 3,
    next_retry_at TIMESTAMP,
    error_message TEXT,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_user ON notifications(recipient_user_id);
CREATE INDEX idx_notifications_status ON notifications(status);
CREATE INDEX idx_notifications_type ON notifications(notification_type);
CREATE INDEX idx_notifications_priority ON notifications(priority);
CREATE INDEX idx_notifications_next_retry ON notifications(next_retry_at) WHERE status = 'RETRYING';
CREATE INDEX idx_notifications_pending ON notifications(created_at) WHERE status = 'PENDING';

COMMENT ON TABLE notifications IS 'Outbound notification queue for emails, SMS, and push notifications';
COMMENT ON COLUMN notifications.status IS 'Delivery status: PENDING, SENT, FAILED, RETRYING';
COMMENT ON COLUMN notifications.priority IS 'Processing priority: HIGH, NORMAL, LOW';
