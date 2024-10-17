CREATE TABLE IF NOT EXISTS daily_attendances (
    company_id              BIGINT NOT NULL             COMMENT 'company id',
    staff_id                BIGINT NOT NULL             COMMENT 'staff id',
    attendance_date         DATE PRIMARY KEY            COMMENT '勤務日',
    work_start              DATETIME NOT NULL           COMMENT '出勤時刻',
    work_end                DATETIME NOT NULL           COMMENT '退勤時刻',
    work_time               INT NOT NULL                COMMENT '労働時間unixtime',
    rest_start              DATETIME NOT NULL           COMMENT '休憩開始時刻',
    rest_end                DATETIME NOT NULL           COMMENT '休憩終了時刻',
    rest_time               INT NOT NULL                COMMENT '休憩時間unixtime',
    shift_start_time        DATETIME NOT NULL           COMMENT 'シフト開始時刻',
    shift_end_time          DATETIME NOT NULL           COMMENT 'シフト終了時刻',
    shift_time              INT NOT NULL                COMMENT 'シフト時間unixtime',
    regular_work_time       INT NOT NULL                COMMENT '所定内労働時間unixtime',
    regular_over_work_time  INT NOT NULL                COMMENT '所定外労働時間unixtime',
    legal_work_time         INT NOT NULL                COMMENT '法定内労働時間unixtime',
    legal_over_work_time    INT NOT NULL                COMMENT '法定外労働時間unixtime',
    night_work_time         INT NOT NULL                COMMENT '夜間労働時間unixtime',
    early_work_time         INT NOT NULL                COMMENT '早出労働時間unixtime',
    health_time             INT NOT NULL                COMMENT '健康管理時間unixtime',
    holiday_work_time       INT NOT NULL                COMMENT '休日労働時間unixtime',
    holiday_type            SMALLINT NOT NULL           COMMENT '休日タイプ 0:無休, 1:全休, 2:半休, 3:時間休 4:legal holiday 5:public holiday 6: paid holiday',
    absence_flag            TINYINT NOT NULL DEFAULT 0  COMMENT '欠勤フラグ 0:not absence, 1:absence',
    late_flag               TINYINT NOT NULL DEFAULT 0  COMMENT '遅刻フラグ 0:not late, 1:late',
    early_flag              TINYINT NOT NULL DEFAULT 0  COMMENT '早退フラグ 0:not early, 1:early',
    lack_rest_flag          TINYINT NOT NULL DEFAULT 0  COMMENT '休憩不足フラグ 0:not lack rest, 1:lack rest',
    created_at              DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS closed_groups (
    company_id      BIGINT PRIMARY KEY              COMMENT 'company id',
    group_id        BIGINT PRIMARY KEY              COMMENT 'closed group id',
    closed_month    SMALLINT PRIMARY KEY NOT NULL   COMMENT '締め月',
    closed_flag     TINYINT NOT NULL DEFAULT 0      COMMENT '締めフラグ 0:not closed(締め解除), 1:closed',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- get closed group
SELECT EXISTS (
    SELECT 1
    FROM closed_groups
    WHERE company_id = 1
        AND group_id = 1
        AND closed_month = 1
        AND closed_flag = 1
);

-- get closed staff
SELECT EXISTS (
    SELECT 1
    FROM staffs AS staff
    INNER JOIN closed_groups AS closed
        ON staff.company_id = closed.company_id
        AND staff.group_id = closed.group_id
    WHERE
        AND staff.company_id = 1
        AND staff.id = 1
        AND closed.closed_month = 1
        AND closed.closed_flag = 1
);
