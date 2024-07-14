-- create table
CREATE TABLE IF NOT EXISTS products (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255),
  embedding VECTOR(3)
);

-- insert data
INSERT IGNORE INTO products (name, embedding) VALUES
  ('Product A', TO_VECTOR('[0.1, 0.2, 0.3]')),
  ('Product B', TO_VECTOR('[0.4, 0.5, 0.6]')),
  ('Product C', TO_VECTOR('[0.7, 0.8, 0.9]'))
;

SELECT
  p1.name AS product1,
  p2.name AS product2,
  VECTOR_TO_STRING(p1.embedding),
  VECTOR_TO_STRING(p2.embedding)
FROM products p1, products p2
WHERE p1.id = p2.id
;
