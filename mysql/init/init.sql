USE test;

-- create table
CREATE TABLE products (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255),
  embedding VECTOR(3)
);

-- insert data
INSERT INTO products (name, embedding) VALUES
  ('Product A', CAST('[0.1, 0.2, 0.3]' AS JSON)),
  ('Product B', CAST('[0.4, 0.5, 0.6]' AS JSON)),
  ('Product C', CAST('[0.7, 0.8, 0.9]' AS JSON))
;

-- コサイン類似度
SELECT
  p1.name AS product1,
  p2.name AS product2,
  ROUND(DOT_PRODUCT(p1.embedding, p2.embedding) / (SQRT(DOT_PRODUCT(p1.embedding, p1.embedding)) * SQRT(DOT_PRODUCT(p2.embedding, p2.embedding))), 2) AS cosine_similarity
FROM products p1, products p2
WHERE p1.id = p2.id
;
