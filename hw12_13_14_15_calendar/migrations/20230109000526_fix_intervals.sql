-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.events ALTER COLUMN duration TYPE interval;
ALTER TABLE public.events ALTER COLUMN notify_term TYPE interval;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.events ALTER COLUMN duration TYPE interval minute;
ALTER TABLE public.events ALTER COLUMN notify_term TYPE interval day;
-- +goose StatementEnd
