CREATE TABLE IF NOT EXISTS stamps (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY
    company_id      BIGINT NOT NULL                 COMMENT 'company id',
    staff_id        BIGINT NOT NULL                 COMMENT 'staff id',
    stamp_date      DATE PRIMARY KEY NOT NULL       COMMENT '勤務日',
    stamp_type      SMALLINT PRIMARY KEY NOT NULL   COMMENT '打刻タイプ 1:work_start, 2:rest_start, 3:rest_end, 4:work_end',
    stamp_time      TIME PRIMARY KEY NOT NULL       COMMENT '打刻日時 && レコード作成日',
    stamp_memo      VARCHAR(255) DEFAULT ""         COMMENT '打刻備考',
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
CREATE INDEX idx_cid_sid_date ON stamps (company_id, staff_id, stamp_date);

CREATE TABLE IF NOT EXISTS stamp_requests (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY
    company_id      BIGINT PRIMARY KEY          COMMENT 'company id',
    staff_id        BIGINT PRIMARY KEY          COMMENT 'staff id',
    stamp_date      DATE PRIMARY KEY            COMMENT '勤務日',
    stamp_type      SMALLINT PRIMARY KEY        COMMENT '打刻タイプ 1:work_start, 2:rest_start, 3:rest_end, 4:work_end',
    stamp_time      TIME PRIMARY KEY NOT NULL   COMMENT '打刻日時 && レコード作成日',
    request_memo    VARCHAR(255) DEFAULT ""     COMMENT '打刻申請備考',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS stamp_request_approves (
    id                  BIGINT AUTO_INCREMENT PRIMARY KEY
    manager_id          BIGINT NOT NULL             COMMENT 'manager id',
    stamp_request_id    BIGINT NOT NULL             COMMENT 'stamp request id',
    approve_flag        TINYINT NOT NULL DEFAULT 0  COMMENT '承認フラグ 0:not approved, 1:approved'
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    FOREIGN KEY (manager_id) REFERENCES staffs(id),
    FOREIGN KEY (stamp_request_id) REFERENCES stamp_requests(id)
);

-- get update record
SELECT 
    id,
    company_id,
    staff_id,
    stamp_date,
    stamp_type,
    stamp_time,
    stamp_memo
FROM stamps
WHERE id = 1;

-- get all stamp per daily staff records
SELECT 
    id,
    company_id,
    staff_id,
    stamp_date,
    stamp_type,
    stamp_time,
    stamp_memo
FROM stamps
WHERE company_id = 1
    AND staff_id = 1
    AND stamp_date = '2024-01-01';