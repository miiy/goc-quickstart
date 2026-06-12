CREATE TABLE IF NOT EXISTS `files`
(
    `id`          bigint unsigned NOT NULL AUTO_INCREMENT,
    `owner_id`    bigint          NOT NULL DEFAULT 0,
    `owner_type`  varchar(64)     NOT NULL DEFAULT '' COMMENT 'owner type: currently user, future post/comment/system',
    `scene`       tinyint(1)      NOT NULL DEFAULT 0 COMMENT 'scene:0 unspecified, 1 avatar',
    `object_key`  varchar(512)    NOT NULL DEFAULT '',
    `url`         varchar(512)    NOT NULL DEFAULT '',
    `mime_type`   varchar(128)    NOT NULL DEFAULT '',
    `size`        bigint          NOT NULL DEFAULT 0 COMMENT 'file size in bytes',
    `checksum`    varchar(128)    NOT NULL DEFAULT '',
    `status`      tinyint(1)      NOT NULL DEFAULT 1 COMMENT 'status:0 unspecified, 1 active, 2 deleted',
    `created_by`  bigint          NOT NULL DEFAULT 0 COMMENT 'creator user id, 0 means system',
    `created_at`  timestamp       NULL     DEFAULT NULL,
    `updated_at`  timestamp       NULL     DEFAULT NULL,
    `deleted_at`  timestamp       NULL     DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_files_owner` (`owner_type`, `owner_id`),
    KEY `idx_files_scene` (`scene`),
    KEY `idx_files_checksum` (`checksum`),
    KEY `idx_files_deleted_at` (`deleted_at`)
) ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci;
