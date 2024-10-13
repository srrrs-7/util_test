CREATE TABLE IF NOT EXISTS shifts (
    id          BIGINT AUTO_INCREMENT PRIMARY KEY,
    company_id  BIGINT NOT NULL COMMENT 'company id',
    name        VARCHAR(50) NOT NULL COMMENT 'shift name',
    start_time  TIME NOT NULL COMMENT 'shift start time',
    end_time    TIME NOT NULL COMMENT 'shift end time',
    deemed_flag TINYINT NOT NULL DEFAULT 0 COMMENT 'みなし勤務シフト 0:normal shift, 1:deemed shift'
    shift_memo  VARCHAR(255) DEFAULT "" COMMENT 'シフト備考',
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES companies(id),
    FOREIGN KEY (staff_id) REFERENCES staffs(id),
);

CREATE TABLE IF NOT EXISTS shift_assignments (
    id          BIGINT AUTO_INCREMENT PRIMARY KEY,
    company_id  BIGINT NOT NULL COMMENT 'company id',
    staff_id    BIGINT NOT NULL COMMENT 'staff id',
    shift_id    BIGINT NOT NULL COMMENT 'シフトタイプ',
    shift_date  DATE NOT NULL COMMENT 'シフト日',
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES companies(id),
    FOREIGN KEY (staff_id) REFERENCES staffs(id),
    FOREIGN KEY (shift_id) REFERENCES shifts(id)
);