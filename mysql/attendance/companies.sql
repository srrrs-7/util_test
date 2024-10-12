CREATE TABLE IF NOT EXISTS companies (
    id                BIGINT AUTO_INCREMENT PRIMARY KEY,
    name              VARCHAR(50) NOT NULL COMMENT '会社名',
    age               INT NOT NULL COMMENT '年齢',
    gender            ENUM('1','2','3') NOT NULL COMMENT '1:woman, 2:man, 3:other',
    representative    VARCHAR(50) NOT NULL COMMENT '代表者指名'
    created_at        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted           TINYINT NOT NULL DEFAULT 0
);