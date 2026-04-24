-- Roundcube database initialization
-- This script is executed only on first container initialization (empty data dir).
-- For subsequent recreates with a fresh container (ephemeral storage) it will run again.

CREATE DATABASE IF NOT EXISTS `roundcube` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER IF NOT EXISTS 'user'@'%' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON `roundcube`.* TO 'user'@'%';
FLUSH PRIVILEGES;

-- Optional: ensure primary application DB also exists (already created via env but idempotent)
CREATE DATABASE IF NOT EXISTS `locky` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;