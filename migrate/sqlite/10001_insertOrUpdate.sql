-- 初始化系统默认任务
INSERT INTO `gm_tasks` (`id`, `uuid`, `sha256`, `user_id`, `title`, `desc`, `type`, `expression`, `method`, `method_params`, `status`, `created_at`, `updated_at`)
VALUES (1, 'c077b134-2458-4868-9371-eb7f58708756', '5cc777305bc35e5048f8d8b06f2097b50e0304ccef4107d3015800d113f7e7db', 0, '刷新用户配置的任务列表', '', 1, '16 */30 * * * *', 1, '', 2, 1755655377, 1755655377);
