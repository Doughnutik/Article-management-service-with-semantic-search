CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS embeddings (
    id INT PRIMARY KEY,
    embedding VECTOR(1024) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_article_embeddings_vector
ON article_embeddings
USING ivfflat (embedding vector_cosine_ops)
WITH (lists = 10);