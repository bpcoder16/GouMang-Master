-- 创建 gm_nodes 节点配置表
CREATE TABLE `gm_nodes` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `title` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '节点标题',
    `ip` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'ip',
    `port` int unsigned NOT NULL DEFAULT '0' COMMENT 'port',
    `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '备注',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态',
    `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间 (时间戳:秒)',
    `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '更新时间 (时间戳:秒)',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='节点配置表';

-- 创建 gm_tasks 任务配置表
CREATE TABLE `gm_tasks` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `uuid` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'uuid',
    `sha256` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '任务内容 sha256',
    `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户id',
    `title` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '任务标题',
    `tag` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '标签',
    `desc` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '任务描述',
    `type` tinyint NOT NULL DEFAULT '1' COMMENT '任务执行类型',
    `expression` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '任务表达式',
    `method` tinyint NOT NULL DEFAULT '0' COMMENT '任务执行方式',
    `method_params` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '任务执行方式参数',
    `next_run_time` bigint NOT NULL DEFAULT '0' COMMENT '下一次执行时间 (时间戳:毫秒)',
    `editable` tinyint NOT NULL DEFAULT '1' COMMENT '是否可编辑，1 可编辑 2 不可编辑',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态',
    `error_message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '错误信息',
    `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间 (时间戳:秒)',
    `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '更新时间 (时间戳:秒)',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_uuid` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='任务配置表';

-- 创建 gm_nodes_tasks 节点&任务关联表
CREATE TABLE `gm_nodes_tasks` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `node_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '节点 ID',
    `task_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '任务 ID',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_node_id_task_id` (`node_id`, `task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='节点&任务关联表';

-- 创建 gm_task_logs 任务日志表
CREATE TABLE `gm_task_logs` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `task_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '任务 ID',
    `task_title` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '任务标题',
    `task_type` tinyint NOT NULL DEFAULT '1' COMMENT '任务执行类型',
    `task_expression` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '任务表达式',
    `node_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '节点 ID',
    `started_at` bigint NOT NULL DEFAULT '0' COMMENT '开始时间 (时间戳:毫秒)',
    `ended_at` bigint NOT NULL DEFAULT '0' COMMENT '结束时间 (时间戳:毫秒)',
    `run_status` tinyint NOT NULL DEFAULT '1' COMMENT '运行状态',
    `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间 (时间戳:秒)',
    `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '更新时间 (时间戳:秒)',
    PRIMARY KEY (`id`),
    KEY `idx_task_id` (`task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='任务日志表';

-- 创建 gm_task_log_detail 任务日志详情表
CREATE TABLE `gm_task_log_detail` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `task_log_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '任务日志 ID',
    `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '内容',
    `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间 (时间戳:秒)',
    PRIMARY KEY (`id`),
    KEY `idx_task_log_id` (`task_log_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='任务日志详情表';