-- Email delivery logs table
CREATE TABLE IF NOT EXISTS email_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    notification_id UUID REFERENCES notifications(id),
    
    -- Email details
    from_email VARCHAR(255) NOT NULL,
    from_name VARCHAR(200),
    to_email VARCHAR(255) NOT NULL,
    cc_emails TEXT[], -- Array of CC addresses
    bcc_emails TEXT[], -- Array of BCC addresses
    
    -- Content
    subject TEXT NOT NULL,
    body_html TEXT,
    body_text TEXT,
    
    -- Delivery status
    status VARCHAR(50) DEFAULT 'PENDING', -- PENDING, SENT, DELIVERED, BOUNCED, FAILED
    smtp_response TEXT,
    message_id VARCHAR(255), -- SMTP message ID
    
    -- Tracking
    opened_at TIMESTAMP,
    clicked_at TIMESTAMP,
    bounced_at TIMESTAMP,
    
    -- Error handling
    error_code VARCHAR(50),
    error_message TEXT,
    
    sent_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_email_logs_notification ON email_logs(notification_id);
CREATE INDEX idx_email_logs_to_email ON email_logs(to_email);
CREATE INDEX idx_email_logs_status ON email_logs(status);
CREATE INDEX idx_email_logs_sent ON email_logs(sent_at);

COMMENT ON TABLE email_logs IS 'Email delivery tracking and logs';
COMMENT ON COLUMN email_logs.status IS 'Delivery status: PENDING, SENT, DELIVERED, BOUNCED, FAILED';
