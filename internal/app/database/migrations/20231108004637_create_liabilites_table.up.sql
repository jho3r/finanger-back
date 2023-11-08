CREATE TABLE liabilities (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    amount NUMERIC(10, 2) NOT NULL,
    description TEXT NOT NULL,
    financial_asset_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    creditor_name VARCHAR(255) NOT NULL,
    category_id INTEGER NOT NULL,
    adquisition_date TIMESTAMP NOT NULL DEFAULT NOW(),
    due_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    FOREIGN KEY (financial_asset_id) REFERENCES financial_assets(id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id)
);
