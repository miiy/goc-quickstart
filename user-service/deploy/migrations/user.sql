CREATE TABLE IF NOT EXISTS `users`
(
    `id`                bigint unsigned NOT NULL AUTO_INCREMENT,
    `username`          varchar(255)    NOT NULL DEFAULT '',
    `password`          varchar(255)    NOT NULL DEFAULT '',
    `email`             varchar(255)    NOT NULL DEFAULT '',
    `email_verified_time` timestamp        NULL     DEFAULT NULL,
    `phone`             varchar(255)    NOT NULL DEFAULT '',
    `unionid`           varchar(255)    NOT NULL DEFAULT 'wechat unionid',
    `mp_openid`         varchar(255)    NOT NULL DEFAULT 'mini program openid',
    `mp_session_key`    varchar(255)    NOT NULL DEFAULT 'mini program session key',
    `status`            tinyint(1)      NOT NULL DEFAULT 0 COMMENT 'status:0 default, 1 active, 2 disable',
    `create_time`        timestamp       NULL     DEFAULT NULL,
    `update_time`        timestamp       NULL     DEFAULT NULL,
    `delete_time`        timestamp       NULL     DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `users_username_unique` (`username`)
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `user_profiles`
(
    `id`                bigint unsigned NOT NULL AUTO_INCREMENT,
    `user_id`           bigint unsigned NOT NULL DEFAULT 0,
    `status`            tinyint(1)      NOT NULL DEFAULT 0 COMMENT 'status:0 default, 1 active, 2 disable',
    `create_time`        timestamp       NULL     DEFAULT NULL,
    `update_time`        timestamp       NULL     DEFAULT NULL,
    `delete_time`        timestamp       NULL     DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `users_user_id_unique` (`user_id`)
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    COLLATE = utf8mb4_unicode_ci;
