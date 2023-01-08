-- +goose Up
-- +goose StatementBegin
DO $$ BEGIN
    CREATE TYPE public.events_notify_status as ENUM ('none', 'blocked', 'notified');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;
ALTER TABLE IF EXISTS public.events ADD COLUMN notify_status public.events_notify_status DEFAULT 'none'::events_notify_status;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.events DROP COLUMN IF EXISTS notify_status;
DROP TYPE IF EXISTS public.events_notify_status;
-- +goose StatementEnd
