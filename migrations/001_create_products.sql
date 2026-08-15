CREATE TABLE products (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    category TEXT NOT NULL DEFAULT '',
    price NUMERIC(12, 2) NOT NULL
);

CREATE INDEX idx_products_title ON products(title);
CREATE INDEX idx_products_category ON products(category);
