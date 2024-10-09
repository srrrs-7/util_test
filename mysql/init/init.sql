-- create table
CREATE TABLE IF NOT EXISTS products (
  id          INT AUTO_INCREMENT PRIMARY KEY,
  name        VARCHAR(255) NOT NULL,
  embedding   VECTOR(3) NOT NULL,
  created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
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
FROM 
  products p1, products p2
WHERE 
  p1.id = p2.id
;


CREATE TABLE IF NOT EXISTS users (
  id          INT AUTO_INCREMENT PRIMARY KEY,
  name        VARCHAR(255) NOT NULL,
  age         INT NOT NULL,
  gender      ENUM('1','2','3') NOT NULL COMMENT '1:woman, 2:man, 3:other',
  created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_scores (
  user_id                 INT NOT NULL COMMENT 'user id',
  commitment_point        INT NOT NULL COMMENT '達成力',
  expertise_skill_point   INT NOT NULL COMMENT '専門スキル',
  resolving_point         INT NOT NULL COMMENT '問題解決力',
  communication_point     INT NOT NULL COMMENT 'コミュニケーション',
  leadership_skill_point  INT NOT NULL COMMENT 'リーダーシップ',
  proactiveness_point     INT NOT NULL COMMENT '積極性',
  responsibility_point    INT NOT NULL COMMENT '責任性',
  cooperativeness_point   INT NOT NULL COMMENT '協調性',
  willingness_point       INT NOT NULL COMMENT '意志力',
  stress_tolerance_point  INT NOT NULL COMMENT 'ストレス耐性',
  integrity_point         INT NOT NULL COMMENT '一貫性',
  accountability_point    INT NOT NULL COMMENT '説明力',
  adaptability_point      INT NOT NULL COMMENT '適応力',
  drive_for_improvement   INT NOT NULL COMMENT '向上',
  created_at              DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at              DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

INSERT INTO users (name, age, gender, created_at, updated_at) VALUES
('佐藤 太郎', 25, 2, NOW(), NOW()),
('山田 花子', 30, 1, NOW(), NOW()),
('田中 一郎', 28, 2, NOW(), NOW()),
('鈴木 愛', 22, 1, NOW(), NOW()),
('高橋 大輔', 35, 2, NOW(), NOW()),
('伊藤 美咲', 27, 1, NOW(), NOW()),
('渡辺 健太', 32, 2, NOW(), NOW()),
('中村 由美', 29, 1, NOW(), NOW()),
('小林 正義', 26, 2, NOW(), NOW()),
('加藤 あゆみ', 31, 1, NOW(), NOW());

-- user_scoresテーブルにデータを挿入
INSERT INTO user_scores (user_id, commitment_point, expertise_skill_point, resolving_point, communication_point, leadership_skill_point, proactiveness_point, responsibility_point, cooperativeness_point, willingness_point, stress_tolerance_point, integrity_point, accountability_point, adaptability_point, drive_for_improvement, created_at, updated_at) VALUES
(1, 8, 7, 6, 9, 5, 7, 8, 9, 6, 7, 8, 9, 6, 7, NOW(), NOW()),
(2, 6, 9, 8, 7, 6, 8, 7, 6, 9, 8, 7, 6, 9, 8, NOW(), NOW()),
(3, 7, 8, 9, 6, 7, 9, 8, 7, 6, 9, 8, 7, 6, 9, NOW(), NOW()),
(4, 9, 7, 6, 8, 5, 7, 8, 9, 6, 7, 8, 9, 6, 7, NOW(), NOW()),
(5, 5, 8, 7, 9, 6, 8, 7, 6, 9, 8, 7, 6, 9, 8, NOW(), NOW()),
(6, 8, 9, 6, 7, 8, 9, 6, 7, 8, 9, 6, 7, 8, 9, NOW(), NOW()),
(7, 7, 6, 8, 9, 7, 6, 8, 9, 7, 6, 8, 9, 7, 6, NOW(), NOW()),
(8, 9, 8, 7, 6, 9, 8, 7, 6, 8, 9, 7, 6, 8, 9, NOW(), NOW()),
(9, 6, 7, 9, 8, 6, 7, 9, 8, 6, 7, 9, 8, 6, 7, NOW(), NOW()),
(10, 8, 6, 7, 9, 8, 6, 7, 9, 8, 6, 7, 9, 8, 6, NOW(), NOW());

-- NTILE query
SELECT 
  *,
  NTILE(3) OVER (PARTITION BY gender ORDER BY age) AS nt
FROM 
  users
;