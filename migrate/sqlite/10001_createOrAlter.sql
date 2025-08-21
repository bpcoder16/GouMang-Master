-- 创建 gm_nodes 节点配置表
CREATE TABLE `gm_nodes` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `title` TEXT NOT NULL DEFAULT '',
    `ip` TEXT NOT NULL DEFAULT '',
    `port` INTEGER NOT NULL DEFAULT 0,
    `remark` TEXT NOT NULL DEFAULT '',
    `status` INTEGER NOT NULL DEFAULT 1,
    `created_at` INTEGER NOT NULL DEFAULT 0,
    `updated_at` INTEGER NOT NULL DEFAULT 0,
    CONSTRAINT `uniq_title` UNIQUE (`title`),
    CONSTRAINT `uniq_ip_port` UNIQUE (`ip`, `port`)
);

-- 创建 gm_tasks 任务配置表
CREATE TABLE `gm_tasks` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `uuid` TEXT NOT NULL DEFAULT '',
    `sha256` TEXT NOT NULL DEFAULT '',
    `user_id` INTEGER NOT NULL DEFAULT 0,
    `title` TEXT NOT NULL DEFAULT '',
    `tag` TEXT NOT NULL DEFAULT '',
    `desc` TEXT NOT NULL DEFAULT '',
    `type` INTEGER NOT NULL DEFAULT 1,
    `expression` TEXT NOT NULL DEFAULT '',
    `method` INTEGER NOT NULL DEFAULT 1,
    `method_params` TEXT NOT NULL DEFAULT '',
    `next_run_time` INTEGER NOT NULL DEFAULT 0,
    `editable` INTEGER NOT NULL DEFAULT 0,
    `status` INTEGER NOT NULL DEFAULT 1,
    `error_message` TEXT NOT NULL DEFAULT '',
    `created_at` INTEGER NOT NULL DEFAULT 0,
    `updated_at` INTEGER NOT NULL DEFAULT 0,
    CONSTRAINT `uniq_title` UNIQUE (`title`)
);

-- 创建 gm_nodes_tasks 节点配置表
CREATE TABLE `gm_nodes_tasks` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `node_id` INTEGER NOT NULL DEFAULT 0,
    `task_id` INTEGER NOT NULL DEFAULT 0,
    CONSTRAINT `uniq_node_id_task_id` UNIQUE (`node_id`, `task_id`)
);

-- 创建 gm_task_logs 任务日志表
CREATE TABLE `gm_task_logs` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `task_id` INTEGER NOT NULL DEFAULT 0,
    `task_title` TEXT NOT NULL DEFAULT '',
    `task_type` INTEGER NOT NULL DEFAULT 1,
    `task_expression` TEXT NOT NULL DEFAULT '',
    `node_id` INTEGER NOT NULL DEFAULT 0,
    `started_at` INTEGER NOT NULL DEFAULT 0,
    `ended_at` INTEGER NOT NULL DEFAULT 0,
    `run_status` INTEGER NOT NULL DEFAULT 0,
    `created_at` INTEGER NOT NULL DEFAULT 0,
    `updated_at` INTEGER NOT NULL DEFAULT 0,
);

-- 创建 gm_task_logs 索引
CREATE INDEX `idx_task_id` ON `gm_task_logs` (`task_id`);

-- 创建 gm_task_log_detail 任务日志详情表
CREATE TABLE `gm_task_log_detail` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `task_log_id` INTEGER NOT NULL DEFAULT 0,
    `content` TEXT NOT NULL DEFAULT '',
    `created_at` INTEGER NOT NULL DEFAULT 0,
);

-- 创建 gm_task_log_detail 索引
CREATE INDEX `idx_task_log_id` ON `gm_task_logs` (`task_log_id`);