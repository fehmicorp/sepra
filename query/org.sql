-- Extensions (safe to run globally)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS citext;

-- Create the dedicated schema for organization and licensing management
CREATE SCHEMA IF NOT EXISTS org;

-- Function to handle updated_at timestamps automatically within the org schema
CREATE OR REPLACE FUNCTION org.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

---

-- 1. Organizations Table
CREATE TABLE IF NOT EXISTS org.organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Licenses Table (Organization-wide limits for DSMs and Agents)
CREATE TABLE IF NOT EXISTS org.licenses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES org.organizations(id) ON DELETE CASCADE,
    max_dsm_count INT NOT NULL DEFAULT 1,
    max_agent_count INT NOT NULL DEFAULT 5,
    license_key VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 3. Deployed DSM Instances Table
CREATE TABLE IF NOT EXISTS org.dsms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES org.organizations(id) ON DELETE CASCADE,
    license_id UUID NOT NULL REFERENCES org.licenses(id) ON DELETE RESTRICT,
    dsm_identifier VARCHAR(255) UNIQUE NOT NULL,
    hostname VARCHAR(255),
    ip_address VARCHAR(45),
    status VARCHAR(50) DEFAULT 'active',
    registered_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_heartbeat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

---

-- Performance Indexes
CREATE INDEX idx_licenses_org ON org.licenses(organization_id);
CREATE INDEX idx_licenses_key ON org.licenses(license_key);
CREATE INDEX idx_dsms_org ON org.dsms(organization_id);
CREATE INDEX idx_dsms_license ON org.dsms(license_id);
CREATE INDEX idx_dsms_identifier ON org.dsms(dsm_identifier);

---

-- Triggers to auto-update the 'updated_at' columns
CREATE TRIGGER update_organizations_modtime BEFORE UPDATE ON org.organizations FOR EACH ROW EXECUTE FUNCTION org.update_updated_at_column();
CREATE TRIGGER update_licenses_modtime BEFORE UPDATE ON org.licenses FOR EACH ROW EXECUTE FUNCTION org.update_updated_at_column();