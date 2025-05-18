-- +goose Up
CREATE TABLE IF NOT EXISTS public.users (
                                     id SERIAL PRIMARY KEY,
                                     username VARCHAR(255) NOT NULL,
                                     password VARCHAR(255) NOT NULL,
                                     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS public.users;
