CREATE TABLE IF NOT EXISTS `posts`
(
    `id`           bigint unsigned NOT NULL AUTO_INCREMENT,
    `author_id`    bigint unsigned NOT NULL DEFAULT 0,
    `title`        varchar(255)    NOT NULL DEFAULT '',
    `content`      text            NOT NULL,
    `status`       tinyint(1)      NOT NULL DEFAULT 0 COMMENT 'status:0 unspecified, 1 draft, 2 published',
    `tags`         varchar(255)    NOT NULL DEFAULT '' COMMENT 'json array',
    `category_id`   bigint unsigned NOT NULL DEFAULT 0,
    `created_at`  timestamp       NULL     DEFAULT NULL,
    `updated_at`  timestamp       NULL     DEFAULT NULL,
    `deleted_at`  timestamp       NULL     DEFAULT NULL,
    PRIMARY KEY (`id`)
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `categories`
(
    `id`           bigint unsigned NOT NULL AUTO_INCREMENT,
    `name`         varchar(255)    NOT NULL DEFAULT '',
    `parent_id`    bigint unsigned NOT NULL DEFAULT 0,
    `path`         varchar(255)    NOT NULL DEFAULT '',
    `created_at`  timestamp       NULL     DEFAULT NULL,
    `updated_at`  timestamp       NULL     DEFAULT NULL,
    `deleted_at`  timestamp       NULL     DEFAULT NULL,
    PRIMARY KEY (`id`)
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    COLLATE = utf8mb4_unicode_ci;
