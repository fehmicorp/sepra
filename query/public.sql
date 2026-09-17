-- 1. Public Users Profile Table
CREATE TABLE IF NOT EXISTS public.users (
    user_id UUID PRIMARY KEY REFERENCES auth.credentials(user_id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES org.organizations(id) ON DELETE CASCADE,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    avatar_url TEXT,
    job_title VARCHAR(100),
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Performance Index
CREATE INDEX idx_public_users_org ON public.users(organization_id);

-- 3. Trigger to auto-update the 'updated_at' column
-- (Reusing the function from the org schema or creating a local public one)
CREATE OR REPLACE FUNCTION public.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_public_users_modtime 
BEFORE UPDATE ON public.users 
FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();