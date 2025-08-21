CREATE TABLE `gm_tasks` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `uuid` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'uuid',
    `sha256` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '任务内容 sha256',
    `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户id',
    `title` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '任务标题',
    `desc` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '任务描述',
    `type` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '任务执行类型',
    `expression` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '任务表达式',
    `method` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '任务执行方式',
    `method_params` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '任务执行方式参数',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态',
    `error_message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '错误信息',
    `created_at` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建时间（时间戳）',
    `updated_at` bigint unsigned NOT NULL DEFAULT '0' COMMENT '更新时间（时间戳）',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_title` (`title`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='任务配置表';