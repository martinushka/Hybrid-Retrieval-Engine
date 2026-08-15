ALTER TABLE products
ADD COLUMN embedding vector(384);

CREATE INDEX idx_products_embedding
ON products
USING hnsw (embedding vector_cosine_ops);
