CREATE TABLE IF NOT EXISTS `posts`
(
    `id`           bigint unsigned NOT NULL AUTO_INCREMENT,
    `user_id`      bigint unsigned NOT NULL DEFAULT 0,
    `title`        varchar(255)    NOT NULL DEFAULT '',
    `summary`      varchar(512)    NOT NULL DEFAULT '',
    `cover_url`    varchar(1024)   NOT NULL DEFAULT '',
    `content`      text            NOT NULL,
    `status`       tinyint(1)      NOT NULL DEFAULT 0 COMMENT 'status:0 unspecified, 1 draft, 2 published, 3 pending_review',
    `tags`         varchar(255)    NOT NULL DEFAULT '' COMMENT 'json array',
    `category_id`   bigint unsigned NOT NULL DEFAULT 0,
    `published_at` timestamp       NULL     DEFAULT NULL,
    `created_at`  timestamp       NULL     DEFAULT NULL,
    `updated_at`  timestamp       NULL     DEFAULT NULL,
    `deleted_at`  timestamp       NULL     DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_posts_user_id` (`user_id`),
    KEY `idx_posts_category_id` (`category_id`)
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
