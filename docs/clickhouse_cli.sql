CREATE TABLE `clickhouse_servers` (
  `id` varchar(36) NOT NULL,
  `org_id` varchar(36) DEFAULT NULL,
  `host` varchar(255) DEFAULT NULL,
  `port` varchar(16) DEFAULT NULL,
  `username` varchar(255) DEFAULT NULL,
  `cluster` varchar(255) DEFAULT NULL,
  `shards` json DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `clickhouse_servers` (`id`, `org_id`, `host`, `port`, `username`, `cluster`, `shards`) VALUES
('1', 'tiki_aff', 'localhost', '9001', 'default', 'tiki_aff', '[\"shard_1\", \"shard_3\"]'),
('2', 'tiki_aff', 'localhost', '9002', 'default', 'tiki_aff', '[\"shard_1\", \"shard_2\"]'),
('3', 'tiki_aff', 'localhost', '9003', 'default', 'tiki_aff', '[\"shard_2\", \"shard_3\"]');


CREATE TABLE `clickhouse_users` (
  `username` varchar(255) NOT NULL,
  `password` text,
  `allow_databases` text,
  `status` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `clickhouse_users` (`username`, `password`, `allow_databases`, `status`) VALUES
('default', NULL, '[]', 'active');