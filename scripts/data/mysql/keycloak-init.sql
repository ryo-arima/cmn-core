-- Keycloak database initialization
-- This script is executed on first container initialization.

CREATE DATABASE IF NOT EXISTS `keycloak` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
GRANT ALL PRIVILEGES ON `keycloak`.* TO 'user'@'%';
FLUSH PRIVILEGES;
