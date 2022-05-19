-- Create the alert rule table.
CREATE TABLE IF NOT EXISTS t_alert_rule
(
    rule_id          bigint not null auto_increment comment '规则id',
    alert_name       varchar(100) default NULL comment '告警名称',
    expression       varchar(500) default NULL comment '告警表达式',
    duration         varchar(20) default NULL comment '告警持续时间',
    alert_level      tinyint default NULL comment '告警严重级别,严重1，高2，中3，低4',
    alert_type       tinyint default NULL comment '告警类型,告警1，预警2，故障3',
    noitce           varchar(2000) default NULL comment '告警通知信息',
    description      varchar(2000) default NULL comment '告警通知描述',
    create_uid       bigint default NULL  comment '创建人',
    state            VARCHAR(2) default NULL COMMENT '数据有效状态 \'U\'：正常 \'E\'：失效',
    create_time      datetime default NULL  comment '创建时间',
    update_uid       bigint default NULL  comment '更新人',
    update_time      datetime default NULL  comment '更新时间',
    tenant_code      VARCHAR(50) default NULL comment '租户code',
    project_id       bigint default NULL comment '项目id',
    system_id        bigint default NULL comment '系统id',
    primary key (rule_id)
);

-- Create the rest of the tables
CREATE TABLE `AlertGroup` (
                              `ID` INT NOT NULL AUTO_INCREMENT,
                              `time` TIMESTAMP NOT NULL,
                              `receiver` VARCHAR(100) NOT NULL,
                              `status` VARCHAR(50) NOT NULL,
                              `externalURL` TEXT NOT NULL,
                              `groupKey` VARCHAR(255) NOT NULL,
                              KEY `idx_time` (`time`) USING BTREE,
                              KEY `idx_status_ts` (`status`, `time`) USING BTREE,
                              PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `GroupLabel` (
                              `ID` INT NOT NULL AUTO_INCREMENT,
                              `AlertGroupID` INT NOT NULL,
                              `GroupLabel` VARCHAR(100) NOT NULL,
                              `Value` VARCHAR(1000) NOT NULL,
                              FOREIGN KEY (AlertGroupID) REFERENCES AlertGroup (ID) ON DELETE CASCADE,
                              PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `CommonLabel` (
                               `ID` INT NOT NULL AUTO_INCREMENT,
                               `AlertGroupID` INT NOT NULL,
                               `Label` VARCHAR(100) NOT NULL,
                               `Value` VARCHAR(1000) NOT NULL,
                               FOREIGN KEY (AlertGroupID) REFERENCES AlertGroup (ID) ON DELETE CASCADE,
                               PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `CommonAnnotation` (
                                    `ID` INT NOT NULL AUTO_INCREMENT,
                                    `AlertGroupID` INT NOT NULL,
                                    `Annotation` VARCHAR(100) NOT NULL,
                                    `Value` VARCHAR(1000) NOT NULL,
                                    FOREIGN KEY (AlertGroupID) REFERENCES AlertGroup (ID) ON DELETE CASCADE,
                                    PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `Alert` (
                         `ID` INT NOT NULL AUTO_INCREMENT,
                         `alertGroupID` INT NOT NULL,
                         `status` VARCHAR(50) NOT NULL,
                         `startsAt` DATETIME NOT NULL,
                         `endsAt` DATETIME DEFAULT NULL,
                         `generatorURL` TEXT NOT NULL,
                         FOREIGN KEY (alertGroupID) REFERENCES AlertGroup (ID) ON DELETE CASCADE,
                         PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `AlertLabel` (
                              `ID` INT NOT NULL AUTO_INCREMENT,
                              `AlertID` INT NOT NULL,
                              `Label` VARCHAR(100) NOT NULL,
                              `Value` VARCHAR(1000) NOT NULL,
                              FOREIGN KEY (AlertID) REFERENCES Alert (ID) ON DELETE CASCADE,
                              PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `AlertAnnotation` (
                                   `ID` INT NOT NULL AUTO_INCREMENT,
                                   `AlertID` INT NOT NULL,
                                   `Annotation` VARCHAR(100) NOT NULL,
                                   `Value` VARCHAR(1000) NOT NULL,
                                   FOREIGN KEY (AlertID) REFERENCES Alert (ID) ON DELETE CASCADE,
                                   PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE Alert
    ADD `fingerprint` TEXT NOT NULL
;
