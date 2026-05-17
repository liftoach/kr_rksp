-- +goose Up
CREATE TABLE categories (
                            id UUID PRIMARY KEY,
                            user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                            name TEXT NOT NULL,
                            transaction_type TEXT NOT NULL
);

-- +goose Down
DROP TABLE categories;