-- +goose Up
CREATE TABLE budgets (
                         id UUID PRIMARY KEY,
                         user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                         category_id UUID NOT NULL REFERENCES categories(id),
                         budget_limit NUMERIC(18,2) NOT NULL,
                         period TEXT NOT NULL,
                         created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE budgets;