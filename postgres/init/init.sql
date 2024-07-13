CREATE EXTENSION vector;

CREATE TABLE vec_main_items (id bigserial PRIMARY KEY, embedding vector(3));
CREATE TABLE vec_sub_items (id bigserial PRIMARY KEY, embedding vector(3));

INSERT INTO vec_main_items (embedding)
VALUES
('[1,2,3]'),
('[4,5,6]')
;

INSERT INTO vec_sub_items (embedding)
VALUES
('[6,5,4]'),
('[3,2,1]')
;

SELECT *
FROM vec_main_items
ORDER BY embedding <-> '[3,1,2]' LIMIT 5;

SELECT *
FROM vec_sub_items
ORDER BY embedding <-> '[3,1,2]' LIMIT 5;

SELECT
  main.id AS vec_main_items,
  sub.id AS vec_sub_items,
  ROUND((main.embedding <-> sub.embedding)::numeric, 2) AS cosine_similarity
FROM vec_main_items AS main, vec_sub_items AS sub
WHERE main.id != sub.id;
