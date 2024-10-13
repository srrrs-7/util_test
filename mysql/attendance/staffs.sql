CREATE TABLE IF NOT EXISTS staffs (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    company_id      BIGINT NOT NULL COMMENT 'company id',
    work_type_id    INT NOT NULL COMMENT 'スタッフ勤務形態',
    group_id        INT NOT NULL COMMENT '所属グループ',
    name            VARCHAR(50) NOT NULL COMMENT 'staff name',
    birth_of_date   DATE NOT NULL COMMENT 'staff birth od date',
    gender          SMALLINT NOT NULL COMMENT '1:woman, 2:man, 3:other',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted         TINYINT NOT NULL DEFAULT 0 COMMENT '削除フラグ 1:deleted, 2:not deleted',
    FOREIGN KEY (company_id) REFERENCES companies(id),
    FOREIGN KEY (work_type_id) REFERENCES work_types(id),
    FOREIGN KEY (group_id) REFERENCES groups(id)
);

CREATE INDEX idx_cid_sid_gid ON stamps (company_id, staff_id, group_id);
CREATE INDEX idx_cid_sid_wid ON stamps (company_id, staff_id, work_type_id);