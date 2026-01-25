-- User in-app notifications table
CREATE TABLE IF NOT EXISTS user_notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    
    -- Content
    title VARCHAR(300) NOT NULL,
    message TEXT NOT NULL,
    
    -- Categorization
    notification_type VARCHAR(50) DEFAULT 'INFO', -- INFO, SUCCESS, WARNING, ERROR
    category VARCHAR(50) DEFAULT 'SYSTEM', -- APPROVAL, ALERT, MESSAGE, SYSTEM
    
    -- Deep linking
    link_url TEXT, -- Link to relevant page
    link_text VARCHAR(100),
    
    -- Related entity (for context)
    entity_type VARCHAR(50), -- PO, PR, WO, LOT, etc.
    entity_id UUID,
    
    -- Read status
    is_read BOOLEAN DEFAULT false,
    read_at TIMESTAMP,
    
    -- Dismiss status
    is_dismissed BOOLEAN DEFAULT false,
    dismissed_at TIMESTAMP,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_notifications_user ON user_notifications(user_id);
CREATE INDEX idx_user_notifications_read ON user_notifications(user_id, is_read);
CREATE INDEX idx_user_notifications_category ON user_notifications(category);
CREATE INDEX idx_user_notifications_created ON user_notifications(created_at DESC);
CREATE INDEX idx_user_notifications_entity ON user_notifications(entity_type, entity_id);

COMMENT ON TABLE user_notifications IS 'In-app notifications displayed to users in the UI';
COMMENT ON COLUMN user_notifications.notification_type IS 'Visual type: INFO, SUCCESS, WARNING, ERROR';
COMMENT ON COLUMN user_notifications.category IS 'Functional category: APPROVAL, ALERT, MESSAGE, SYSTEM';
