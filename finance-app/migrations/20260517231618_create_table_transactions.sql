-- +goose Up
CREATE TABLE transactions (
                              id UUID PRIMARY KEY,
                              user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                              amount NUMERIC(18,2) NOT NULL,
                              type TEXT NOT NULL,
                              category TEXT NOT NULL,
                              description TEXT,
                              created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE transactions;