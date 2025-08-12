-- mysql
CREATE TABLE `gm_tasks` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `uuid` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'uuid',
    `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户id',
    `title` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '任务标题',
    `desc` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '任务描述',
    `type` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '任务执行类型',
    `expression` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '任务表达式',
    `method` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '任务执行方式',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '-1 删除 1 待启用 2 已启用',
    `created_at` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
    `updated_at` bigint unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='任务配置表';

-- sqlite
CREATE TABLE gm_tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uuid TEXT NOT NULL DEFAULT '',
    user_id INTEGER NOT NULL DEFAULT 0,
    title TEXT NOT NULL DEFAULT '',
    desc TEXT NOT NULL DEFAULT '',
    type INTEGER NOT NULL DEFAULT 1,
    expression TEXT NOT NULL DEFAULT '',
    method INTEGER NOT NULL DEFAULT 1,
    status INTEGER NOT NULL DEFAULT 1,
    created_at INTEGER NOT NULL DEFAULT 0,
    updated_at INTEGER NOT NULL DEFAULT 0
);