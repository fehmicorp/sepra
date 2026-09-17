-- Extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS citext;

-- Create auth schema
CREATE SCHEMA IF NOT EXISTS auth;

-- Function to handle updated_at timestamps automatically
CREATE OR REPLACE FUNCTION auth.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

---

-- 1. Credentials Table (Scoped to an Organization)
CREATE TABLE IF NOT EXISTS auth.credentials (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES org.organizations(id) ON DELETE CASCADE,
    username CITEXT UNIQUE,
    password_hash TEXT,
    status_code SMALLINT DEFAULT 1,
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Identities Table (Email, Phone, Username, OAuth linked to user)
CREATE TABLE IF NOT EXISTS auth.identities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.credentials (user_id) ON DELETE CASCADE,
    identity_type VARCHAR(10) CHECK (identity_type IN ('email', 'phone', 'username', 'oauth')),
    identifier CITEXT NOT NULL,
    provider VARCHAR(20) DEFAULT 'local',
    provider_id TEXT,
    is_verified BOOLEAN DEFAULT false,
    is_primary BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (identity_type, identifier, provider)
);

-- 3. Multi-Factor Authentication (MFA)
CREATE TABLE IF NOT EXISTS auth.mfa (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.credentials (user_id) ON DELETE CASCADE,
    mfa_enabled BOOLEAN DEFAULT false,
    mfa_type VARCHAR(20) CHECK (mfa_type IN ('totp', 'sms', 'email')),
    secret TEXT NOT NULL,
    recovery_codes TEXT[],
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 4. Sessions
CREATE TABLE IF NOT EXISTS auth.sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.credentials (user_id) ON DELETE CASCADE,
    refresh_token TEXT UNIQUE NOT NULL,
    user_agent TEXT,
    ip_address INET,
    is_revoked BOOLEAN DEFAULT false,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 5. Login History
CREATE TABLE IF NOT EXISTS auth.login_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.credentials (user_id) ON DELETE CASCADE,
    login_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    logout_at TIMESTAMP WITH TIME ZONE,
    ip_address INET,
    user_agent TEXT,
    success BOOLEAN DEFAULT true,
    failure_reason TEXT,
    metadata JSONB DEFAULT '{}'::jsonb
);

-- 6. Mobile Verifications (OTP Tracking)
CREATE TABLE IF NOT EXISTS auth.mobile_number_verifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.credentials (user_id) ON DELETE CASCADE,
    verification_code TEXT,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 7. Mobile Device Registrations
CREATE TABLE IF NOT EXISTS auth.mobile_device_registrations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.credentials (user_id) ON DELETE CASCADE,
    device_token TEXT UNIQUE NOT NULL,
    device_type VARCHAR(20) CHECK (device_type IN ('ios', 'android')),
    registered_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_location JSONB,
    last_used_at TIMESTAMP WITH TIME ZONE
);

-- 8. Security Questions
CREATE TABLE IF NOT EXISTS auth.security_questions (
    user_id UUID PRIMARY KEY REFERENCES auth.credentials (user_id) ON DELETE CASCADE,
    question TEXT[] NOT NULL,
    answer_hash TEXT[] NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

---

-- Performance Indexes
CREATE INDEX idx_credentials_org ON auth.credentials(organization_id);
CREATE INDEX idx_identities_identifier ON auth.identities(identifier);
CREATE INDEX idx_sessions_user ON auth.sessions(user_id);

---

-- Triggers to auto-update the 'updated_at' column
CREATE TRIGGER update_credentials_modtime BEFORE UPDATE ON auth.credentials FOR EACH ROW EXECUTE FUNCTION auth.update_updated_at_column();
CREATE TRIGGER update_identities_modtime BEFORE UPDATE ON auth.identities FOR EACH ROW EXECUTE FUNCTION auth.update_updated_at_column();
CREATE TRIGGER update_mfa_modtime BEFORE UPDATE ON auth.mfa FOR EACH ROW EXECUTE FUNCTION auth.update_updated_at_column();
CREATE TRIGGER update_sessions_modtime BEFORE UPDATE ON auth.sessions FOR EACH ROW EXECUTE FUNCTION auth.update_updated_at_column();
CREATE TRIGGER update_security_questions_modtime BEFORE UPDATE ON auth.security_questions FOR EACH ROW EXECUTE FUNCTION auth.update_updated_at_column();