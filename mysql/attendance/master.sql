CREATE TABLE IF NOT EXISTS companies (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    name            VARCHAR(50) NOT NULL        COMMENT '会社名',
    representative  VARCHAR(50) NOT NULL        COMMENT '代表者指名',
    deleted         TINYINT NOT NULL DEFAULT 0  COMMENT '削除フラグ 1:deleted, 2:not deleted',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
);


CREATE TABLE IF NOT EXISTS staffs (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    company_id      BIGINT NOT NULL             COMMENT 'company id',
    work_type_id    INT NOT NULL                COMMENT 'スタッフ勤務形態',
    group_id        INT NOT NULL                COMMENT '所属グループ',
    name            VARCHAR(50) NOT NULL        COMMENT 'staff name',
    birth_of_date   DATE NOT NULL               COMMENT 'staff birth od date',
    gender          SMALLINT NOT NULL           COMMENT '1:woman, 2:man, 3:other',
    deleted         TINYINT NOT NULL DEFAULT 0  COMMENT '削除フラグ 1:deleted, 2:not deleted',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES companies(id),
    FOREIGN KEY (work_type_id) REFERENCES work_types(id),
    FOREIGN KEY (group_id) REFERENCES groups(id)
);
CREATE INDEX idx_cid_sid_gid ON stamps (company_id, staff_id, group_id);
CREATE INDEX idx_cid_sid_wid ON stamps (company_id, staff_id, work_type_id);


CREATE TABLE IF NOT EXISTS groups (
    id          BIGINT AUTO_INCREMENT PRIMARY KEY,
    company_id  BIGINT PRIMARY KEY          COMMENT 'company id',
    manager_id  BIGINT NOT NULL             COMMENT 'manager id',
    name        VARCHAR(100) NOT NULL       COMMENT 'グループ名',
    parent_id   BIGINT DEFAULT 0            COMMENT '親グループのID(トップレベルの場合は0)',
    description VARCHAR(255) DEFAULT ''     COMMENT 'グループの説明',
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES companies(id),
    FOREIGN KEY (manager_id) REFERENCES staffs(id)
);


CREATE TABLE IF NOT EXISTS work_types (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    name            VARCHAR(50) NOT NULL        COMMENT 'work type name',
    company_id      BIGINT NOT NULL             COMMENT 'company id',
    staff_id        BIGINT NOT NULL             COMMENT 'staff id',
    deleted         TINYINT NOT NULL DEFAULT 0  COMMENT '削除フラグ 1:deleted, 2:not deleted',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES companies(id),
    FOREIGN KEY (staff_id) REFERENCES staffs(id)
);

CREATE TABLE IF NOT EXISTS wages (
    company_id              BIGINT PRIMARY KEY      COMMENT 'company id',
    work_wage               INT NOT NULL DEFAULT 0  COMMENT '通常賃金',
    night_work_wage_rate    INT NOT NULL DEFAULT 0  COMMENT '夜間労働割増賃金',
    early_work_wage_rate    INT NOT NULL DEFAULT 0  COMMENT '早出労働割増賃金',
    holiday_work_wage_rate  INT NOT NULL DEFAULT 0  COMMENT '休日労働割増賃金',
    created_at              DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES companies(id)
);


-- get group structure (recursive CTE query)
WITH RECURSIVE group_hierarchy AS (
    SELECT id, name, parent_id
    FROM groups
    WHERE company_id = 1
        AND id = 1
UNION ALL
    SELECT g.id, g.name, g.parent_id
    FROM groups g
    INNER JOIN group_hierarchy gh
        ON company_id = 1 
        AND g.parent_id = gh.id
)
SELECT * FROM group_hierarchy;