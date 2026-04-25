-- =============================================================================
-- PostgreSQL initialisation script for cmn-core local development
-- Executed automatically by the postgres container on first startup.
-- =============================================================================

-- ---- keycloak ---------------------------------------------------------------
CREATE DATABASE keycloak;
CREATE USER keycloak WITH PASSWORD 'keycloak';
GRANT ALL PRIVILEGES ON DATABASE keycloak TO keycloak;
\c keycloak
GRANT ALL ON SCHEMA public TO keycloak;

-- ---- roundcube ---------------------------------------------------------------
CREATE DATABASE roundcube;
CREATE USER roundcube WITH PASSWORD 'roundcube';
GRANT ALL PRIVILEGES ON DATABASE roundcube TO roundcube;
\c roundcube
GRANT ALL ON SCHEMA public TO roundcube;

-- ---- casdoor ---------------------------------------------------------------
CREATE DATABASE casdoor;
CREATE USER casdoor WITH PASSWORD 'casdoor';
GRANT ALL PRIVILEGES ON DATABASE casdoor TO casdoor;
\c casdoor
GRANT ALL ON SCHEMA public TO casdoor;

-- ---- cmn_core (application) --------------------------------------------------
-- The default POSTGRES_DB=cmn_core is already created by the container.
-- Grant the default user (user) full access.
GRANT ALL PRIVILEGES ON DATABASE cmn_core TO "user";
\c cmn_core
GRANT ALL ON SCHEMA public TO "user";
