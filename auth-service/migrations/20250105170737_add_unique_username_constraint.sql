-- +goose Up
ALTER TABLE public.users ADD CONSTRAINT unique_username UNIQUE (username);

-- +goose Down
ALTER TABLE public.users DROP CONSTRAINT unique_username;
