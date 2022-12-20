-- +goose Up
-- +goose StatementBegin
INSERT INTO public.users
VALUES ('ab8e3706-7ad8-11ed-95f7-d00d1b9e4cfe', 'Ivan', 'ivan@otus.ru'),
       ('90bdce82-7ad8-11ed-99c1-d00d1b9e4cfe', 'Vitaly', 'vitaly@mail.ru'),
       ('973454b8-7ae0-11ed-97ae-d00d1b9e4cfe', 'Sergey', 'gray@yandex.ru');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM public.users
WHERE id IN
      ('ab8e3706-7ad8-11ed-95f7-d00d1b9e4cfe',
       '90bdce82-7ad8-11ed-99c1-d00d1b9e4cfe',
       '973454b8-7ae0-11ed-97ae-d00d1b9e4cfe');
-- +goose StatementEnd
