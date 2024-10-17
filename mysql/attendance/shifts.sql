CREATE TABLE IF NOT EXISTS shifts (
    id          BIGINT AUTO_INCREMENT PRIMARY KEY,
    company_id  BIGINT PRIMARY KEY          COMMENT 'company id',
    name        VARCHAR(50) NOT NULL        COMMENT 'shift name',
    start_time  TIME NOT NULL               COMMENT 'shift start time',
    end_time    TIME NOT NULL               COMMENT 'shift end time',
    deemed_flag TINYINT NOT NULL DEFAULT 0  COMMENT 'みなし勤務シフト 0:normal shift, 1:deemed shift'
    shift_memo  VARCHAR(255) DEFAULT ""     COMMENT 'シフト備考',
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS shift_assignments (
    company_id  BIGINT PRIMARY KEY  COMMENT 'company id',
    staff_id    BIGINT PRIMARY KEY  COMMENT 'staff id',
    shift_date  DATE PRIMARY KEY    COMMENT 'シフト日',
    shift_id    BIGINT PRIMARY KEY  COMMENT 'シフトタイプ',
    shift_time  INT NOT NULL        COMMENT 'シフト時間unixtime',
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES companies(id),
    FOREIGN KEY (staff_id) REFERENCES staffs(id),
    FOREIGN KEY (shift_id) REFERENCES shifts(id)
);

CREATE TABLE IF NOT EXISTS shift_requests (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY
    company_id      BIGINT PRIMARY KEY      COMMENT 'company id',
    staff_id        BIGINT PRIMARY KEY      COMMENT 'staff id',
    shift_date      DATE PRIMARY KEY        COMMENT '勤務日',
    start_time      TIME NOT NULL           COMMENT 'shift start time',
    end_time        TIME NOT NULL           COMMENT 'shift end time',
    request_memo    VARCHAR(255) DEFAULT "" COMMENT 'シフト申請備考',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS shift_request_approves (
    id                  BIGINT AUTO_INCREMENT PRIMARY KEY
    manager_id          BIGINT NOT NULL             COMMENT 'manager id',
    shift_request_id    BIGINT NOT NULL             COMMENT 'shift request id',
    approve_flag        TINYINT NOT NULL DEFAULT 0  COMMENT '承認フラグ 0:not approved, 1:approved'
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    FOREIGN KEY (manager_id) REFERENCES staffs(id),
    FOREIGN KEY (shift_request_id) REFERENCES shift_requests(id)
);

-- get all shift master
SELECT 
    id,
    name,
    start_time,
    end_time,
    shift_memo,
    deemed_flag
FROM shifts
WHERE company_id = 1;

-- ger shift assignment per employee and date
SELECT
    id,
    staff_id,
    shift_date,
    shift_id
FROM shift_assignments
WHERE company_id = 1
    AND shift_date BETWEEN '2024-01-01' AND '2024-01-31';

-- get all shift assign daily staff detail records
SELECT
    staff.id            AS staff_id,
    staff.name          AS staff_name,
    shift.id            AS shift_id,
    shift.name          AS shift_name,
    assign.shift_date   AS shift_date,
    shift.start_time    AS shift_start_time,
    shift.end_time      AS shift_end_time,
    shift.shift_memo    AS shift_memo,
    shift.deemed_flag   AS shift_deemed_flag
FROM shift_assignments AS assign
INNER JOIN shifts as shift
    ON assign.company_id = shift.company_id
    AND assign.shift_id = shift.id
    AND assign.company_id = 1
    AND assign.shift_date BETWEEN '2024-01-01' AND '2024-01-31';
INNER JOIN staffs as staff
    ON assign.company_id = staff.company_id
    AND assign.staff_id = staff.id
;

-- check exist duplicate shift
SELECT EXISTS (
    SELECT 1
    FROM shift_assignments AS assign
    WHERE company_id = 1
        AND staff_id = 1
        AND shift_date = '2024-01-01'
        AND shift_id = 1
);

-- check exist duplicate part time shift
SELECT EXISTS (
    SELECT 1
    FROM shift_assignments AS assign
    INNER JOIN shifts as shift
        ON assign.shift_id = shift.id
        AND assign.company_id = shift.company_id
        AND assign.company_id = 1
        AND assign.staff_id = 1
        AND assign.shift_date = '2024-01-01'
    WHERE shift.start_time >= '09:00:00' AND shift.end_time <= '18:00:00';
);