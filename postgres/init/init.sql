-- vector extension install
CREATE EXTENSION IF NOT EXISTS pgvector;

-- create table
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name TEXT,
    features vector(100)
);

-- insert data
INSERT INTO products (name, features) VALUES
  ('Product A', '{1.2, 3.4, 5.6, 7.8, 9.0, 0.1, 0.2, 0.3, 0.4, 0.5}'),
  ('Product B', '{0.1, 0.2, 0.3, 0.4, 0.5, 1.2, 3.4, 5.6, 7.8, 9.0}')
;

-- create index
CREATE INDEX products_features_idx ON products USING GIST (features gist_vector_ops); 