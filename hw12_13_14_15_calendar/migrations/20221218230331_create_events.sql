-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.events (
    id uuid NOT NULL,
    title character varying(255) NOT NULL,
    date timestamp with time zone NOT NULL,
    duration interval minute NOT NULL,
    owner_id uuid NOT NULL,
    description text,
    notify_term interval day,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    PRIMARY KEY (id),
    CONSTRAINT owner_id_fkey FOREIGN KEY (owner_id)
        REFERENCES public.users(id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.events;
-- +goose StatementEnd
