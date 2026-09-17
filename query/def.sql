ALTER TABLE public.users 
ADD COLUMN IF NOT EXISTS is_protected BOOLEAN DEFAULT FALSE;

CREATE OR REPLACE FUNCTION public.prevent_admin_removal_or_disable()
RETURNS TRIGGER AS $$
BEGIN
    -- Check if the user being modified is a protected admin
    IF OLD.is_protected = TRUE THEN
        -- If it's a DELETE operation, block it
        IF (TG_OP = 'DELETE') THEN
            RAISE EXCEPTION 'Security Error: Protected system admin accounts cannot be deleted.';
        END IF;

        -- If it's an UPDATE operation checking status/active state, block disabling
        IF (TG_OP = 'UPDATE') AND (NEW.is_active = FALSE) THEN
            RAISE EXCEPTION 'Security Error: Protected system admin accounts cannot be disabled.';
        END IF;
    END IF;

    IF (TG_OP = 'DELETE') THEN
        RETURN OLD;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Attach the trigger to block deletions and deactivations on public.users
DROP TRIGGER IF EXISTS trg_protect_system_admin ON public.users;
CREATE TRIGGER trg_protect_system_admin
    BEFORE DELETE OR UPDATE ON public.users
    FOR EACH ROW
    EXECUTE FUNCTION public.prevent_admin_removal_or_disable();

INSERT INTO org.organizations (id, name, is_active)
VALUES ('00000000-0000-0000-0000-000000000001', 'System Administration', TRUE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO auth.credentials (user_id, organization_id, username, password_hash, status_code)
VALUES (
    '11111111-1111-1111-1111-111111111111',
    '00000000-0000-0000-0000-000000000001',
    'sysadmin',
    crypt('$uper@dmin', gen_salt('bf')),
    1
)
ON CONFLICT (user_id) DO NOTHING;

INSERT INTO public.users (user_id, organization_id, first_name, last_name, job_title, is_protected)
VALUES (
    '11111111-1111-1111-1111-111111111111',
    '00000000-0000-0000-0000-000000000001',
    'System',
    'Administrator',
    'Super Admin',
    TRUE
)
ON CONFLICT (user_id) DO UPDATE SET is_protected = TRUE;