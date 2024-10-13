CREATE TABLE IF NOT EXISTS companies (
    id                BIGINT AUTO_INCREMENT PRIMARY KEY,
    name              VARCHAR(50) NOT NULL COMMENT '会社名',
    representative    VARCHAR(50) NOT NULL COMMENT '代表者指名'
    created_at        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted           TINYINT NOT NULL DEFAULT 0
);