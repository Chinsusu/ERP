-- BOM Versions - version history tracking
CREATE TABLE IF NOT EXISTS bom_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bom_id UUID NOT NULL REFERENCES boms(id) ON DELETE CASCADE,
    version INT NOT NULL,
    change_reason TEXT,
    changed_by UUID,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    snapshot JSONB NOT NULL, -- Full BOM snapshot at this version
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bom_versions_bom_id ON bom_versions(bom_id);
CREATE UNIQUE INDEX idx_bom_versions_bom_version ON bom_versions(bom_id, version);
