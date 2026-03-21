CREATE TABLE transactions (
                              id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                              user_id UUID NOT NULL,
                              amount NUMERIC(18, 8) NOT NULL,
                              currency VARCHAR(10) NOT NULL,
                              transaction_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                              categories TEXT[],
                              account VARCHAR(100),
                              receipt_url TEXT,
                              created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);