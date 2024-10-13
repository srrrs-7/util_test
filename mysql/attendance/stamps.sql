CREATE TABLE IF NOT EXISTS stamps (
    id          BIGINT AUTO_INCREMENT PRIMARY KEY
    company_id  BIGINT NOT NULL COMMENT 'company id',
    staff_id    BIGINT NOT NULL COMMENT 'staff id',
    stamp_date  DATE PRIMARY KEY NOT NULL COMMENT '勤務日',
    stamp_type  SMALLINT PRIMARY KEY NOT NULL COMMENT '打刻タイプ 1:work_start, 2:rest_start, 3:rest_end, 4:work_end',
    stamp_time  DATETIME PRIMARY KEY NOT NULL COMMENT '打刻日時 && レコード作成日',
    stamp_memo  VARCHAR(255) DEFAULT "" COMMENT '打刻備考',
    updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES companies(id),
    FOREIGN KEY (staff_id) REFERENCES staffs(id)
);

CREATE INDEX idx_cid_sid_date ON stamps (company_id, staff_id, stamp_date);