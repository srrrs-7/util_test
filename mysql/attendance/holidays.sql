CREATE TABLE IF NOT EXISTS holidays (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    company_id      BIGINT PRIMARY KEY          COMMENT 'company id',
    name            VARCHAR(50) NOT NULL        COMMENT 'holiday name',
    start_time      TIME NOT NULL               COMMENT 'holiday start time',
    end_time        TIME NOT NULL               COMMENT 'holiday end time',
    deemed_flag     TINYINT NOT NULL DEFAULT 0  COMMENT 'みなし休日 0:normal holiday, 1:deemed holiday'
    holiday_memo    VARCHAR(255) DEFAULT ""     COMMENT '休日備考',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);


-- 全休、半休、時間休対応
CREATE TABLE IF NOT EXISTS holiday_assignments (
    company_id      BIGINT PRIMARY KEY          COMMENT 'company id',
    staff_id        BIGINT PRIMARY KEY          COMMENT 'staff id',
    holiday_date    DATE PRIMARY KEY            COMMENT '休日日',
    holiday_id      BIGINT PRIMARY KEY          COMMENT '休日タイプ 有給、振休、代休、祝日 etc',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (holiday_id) REFERENCES holidays(id)
);


CREATE TABLE IF NOT EXISTS staff_paid_holidays (
    company_id      BIGINT PRIMARY KEY          COMMENT 'company id',
    staff_id        BIGINT PRIMARY KEY          COMMENT 'staff id',
    paid_holiday_count INT NOT NULL DEFAULT 0   COMMENT 'paid holiday count',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS holiday_requests (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY
    company_id      BIGINT PRIMARY KEY          COMMENT 'company id',
    staff_id        BIGINT PRIMARY KEY          COMMENT 'staff id',
    holiday_date    DATE PRIMARY KEY            COMMENT '勤務日',
    holiday_id      BIGINT PRIMARY KEY          COMMENT '休日タイプ 有給、振休、代休、祝日 etc',
    request_memo    VARCHAR(255) DEFAULT ""     COMMENT '休暇申請備考',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    FOREIGN KEY (holiday_id) REFERENCES holidays(id)
);

CREATE TABLE IF NOT EXISTS holiday_request_approves (
    id                  BIGINT AUTO_INCREMENT PRIMARY KEY
    manager_id          BIGINT NOT NULL             COMMENT 'manager id',
    holiday_request_id    BIGINT NOT NULL           COMMENT 'holiday request id',
    approve_flag        TINYINT NOT NULL DEFAULT 0  COMMENT '承認フラグ 0:not approved, 1:approved'
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    FOREIGN KEY (manager_id) REFERENCES staffs(id),
    FOREIGN KEY (holiday_request_id) REFERENCES holiday_requests(id)
);

-- get holiday master
SELECT 
    id,
    name,
    start_time,
    end_time,
    holiday_memo,
    deemed_flag
FROM holidays
WHERE company_id = 1;

-- ger holiday assignment daily staff holiday records
SELECT
    id,
    staff_id,
    holiday_date,
    holiday_id
FROM holiday_assignments
WHERE company_id = 1
    AND holiday_date BETWEEN '2024-01-01' AND '2024-01-31';

-- get holiday assign detail
SELECT
    staff.id                AS staff_id,
    staff.name              AS staff_name,
    holiday.id              AS holiday_id,
    holiday.name            AS holiday_name,
    assign.holiday_date     AS holiday_date,
    holiday.start_time      AS holiday_start_time,
    holiday.end_time        AS holiday_end_time,
    holiday.holiday_memo    AS holiday_memo,
    holiday.deemed_flag     AS holiday_deemed_flag
FROM holiday_assignments AS assign
INNER JOIN holidays as holiday
    ON assign.holiday_id = holiday.id
    AND assign.company_id = holiday.company_id
    AND assign.company_id = 1
    AND assign.holiday_date BETWEEN '2024-01-01' AND '2024-01-31';
INNER JOIN staffs as staff
    ON assign.company_id = staff.company_id
    AND assign.staff_id = staff.id
;

-- check exist duplicate holiday
SELECT EXISTS (
    SELECT 1
    FROM holiday_assignments AS assign
    WHERE company_id = 1
        AND staff_id = 1
        AND holiday_date = '2024-01-01'
        AND holiday_id = 1
);

-- check exist duplicate part time holiday
SELECT EXISTS (
    SELECT 1
    FROM holiday_assignments AS assign
    INNER JOIN holidays as holiday
        ON assign.holiday_id = holiday.id
        AND assign.company_id = holiday.company_id
        AND assign.company_id = 1
        AND assign.staff_id = 1
        AND assign.holiday_date = '2024-01-01'
    WHERE holiday.start_time >= '09:00:00' AND holiday.end_time <= '18:00:00';
);
