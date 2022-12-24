-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.users (
   id uuid NOT NULL,
   name character varying(255) NOT NULL,
   email character varying(50) NOT NULL,
   CONSTRAINT users_email_key UNIQUE (email),
   PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.users;
-- +goose StatementEnd
