-- =============================================================================
-- PostgreSQL initialisation script for cmn-core local development
-- Executed automatically by the postgres container on first startup.
-- =============================================================================

-- ---- keycloak ---------------------------------------------------------------
CREATE DATABASE keycloak;
CREATE USER keycloak WITH PASSWORD 'keycloak';
GRANT ALL PRIVILEGES ON DATABASE keycloak TO keycloak;

-- ---- roundcube ---------------------------------------------------------------
CREATE DATABASE roundcube;
CREATE USER roundcube WITH PASSWORD 'roundcube';
GRANT ALL PRIVILEGES ON DATABASE roundcube TO roundcube;

-- ---- authentik ---------------------------------------------------------------
CREATE DATABASE authentik;
CREATE USER authentik WITH PASSWORD 'authentik';
GRANT ALL PRIVILEGES ON DATABASE authentik TO authentik;

-- ---- cmn_core (application) --------------------------------------------------
-- The default POSTGRES_DB=cmn_core is already created by the container.
-- Grant the default user (user) full access.
GRANT ALL PRIVILEGES ON DATABASE cmn_core TO "user";
