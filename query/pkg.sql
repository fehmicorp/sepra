CREATE SCHEMA IF NOT EXISTS pkg;

-- Function to handle updated_at timestamps automatically within the pkg schema
CREATE OR REPLACE FUNCTION pkg.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TABLE IF NOT EXISTS pkg.packages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,                          -- e.g., "dsm-installer", "endpoint-agent"
    target_type VARCHAR(50) NOT NULL CHECK (target_type IN ('dsm', 'agent')), -- Target component type
    version VARCHAR(50) NOT NULL,                        -- Semantic version e.g., "1.2.4"
    file_path TEXT NOT NULL,                             -- Storage path or object storage URL
    checksum_sha256 VARCHAR(64) NOT NULL,                -- For integrity verification by DSM/Agent
    file_size BIGINT NOT NULL,                           -- Size in bytes
    is_active BOOLEAN DEFAULT TRUE,                      -- Whether this version is available for deployment
    release_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (target_type, version)
);

-- 2. Package Downloads & Update Audit Log Table
CREATE TABLE IF NOT EXISTS pkg.downloads_audit (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES org.organizations(id) ON DELETE SET NULL, -- Optional org tracking
    package_id UUID NOT NULL REFERENCES pkg.packages(id) ON DELETE RESTRICT,
    entity_type VARCHAR(50) NOT NULL CHECK (entity_type IN ('dsm', 'agent')), -- Who requested it
    entity_identifier VARCHAR(255) NOT NULL,             -- Unique DSM identifier or Agent UUID/ID
    ip_address INET,
    status VARCHAR(50) DEFAULT 'success',                -- 'success', 'failed', 'downloading'
    metadata JSONB DEFAULT '{}'::jsonb,                  -- Extra info (e.g., OS architecture, previous version)
    downloaded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

---

-- Performance Indexes
CREATE INDEX idx_packages_target_version ON pkg.packages(target_type, version);
CREATE INDEX idx_packages_active ON pkg.packages(is_active);
CREATE INDEX idx_audit_org ON pkg.downloads_audit(organization_id);
CREATE INDEX idx_audit_package ON pkg.downloads_audit(package_id);
CREATE INDEX idx_audit_entity ON pkg.downloads_audit(entity_type, entity_identifier);

---

-- Trigger to auto-update the 'updated_at' column on packages
CREATE TRIGGER update_packages_modtime BEFORE UPDATE ON pkg.packages FOR EACH ROW EXECUTE FUNCTION pkg.update_updated_at_column();