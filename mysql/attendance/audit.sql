CREATE TABLE IF NOT EXISTS audits (
    conpaniny_id    BIGINT NOT NULL COMMENT 'company id',
    manager_id      BIGINT NOT NULL COMMENT 'manager id',
    staff_id        BIGINT NOT NULL COMMENT 'staff id',
    audit_date      DATE NOT NULL   COMMENT 'audit date',
    domain_id       BIGINT NOT NULL COMMENT 'service domain id',
    action_id       BIGINT NOT NULL COMMENT 'action id',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (domain_id) REFERENCES domains(id),
    FOREIGN KEY (action_id) REFERENCES actions(id)
);

CREATE TABLE IF NOT EXISTS domains (
    id          BIGINT AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(50) NOT NULL COMMENT 'service domain name',
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS actions (
    id          BIGINT AUTO_INCREMENT PRIMARY KEY,
    action      VARCHAR(50) NOT NULL COMMENT 'action',
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);