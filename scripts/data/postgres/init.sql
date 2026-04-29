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

-- ---- pg_groups (Casdoor group UUID <-> display name mapping) -----------------
-- Created here so seed data can be inserted on first startup.
-- GORM AutoMigrate will reconcile the schema when the server starts.
CREATE TABLE IF NOT EXISTS pg_groups (
    id         BIGSERIAL    PRIMARY KEY,
    uuid       TEXT         NOT NULL,
    name       TEXT         NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_pg_groups_uuid ON pg_groups (uuid);

INSERT INTO pg_groups (uuid, name, created_at, updated_at) VALUES
  ('e4e83e99-0e8a-5d7b-9b0b-1fd02db4dad2', 'group001', NOW(), NOW()),
  ('08990183-8480-529a-b485-7b2260e1e7dd', 'group002', NOW(), NOW()),
  ('154d6397-ed14-51c8-b855-b87cf70844b9', 'group003', NOW(), NOW()),
  ('a3f1a256-5f34-5abe-a097-c7c4f3321450', 'group004', NOW(), NOW()),
  ('1297d847-237e-5b75-a10d-6e565b0dd102', 'group005', NOW(), NOW()),
  ('36d83b6c-81b6-52b4-bc45-10ec7af262b0', 'group006', NOW(), NOW()),
  ('5cb36dd8-f386-5695-918f-8384ae705d9f', 'group007', NOW(), NOW()),
  ('23616691-df42-5462-91f0-4da37959aac5', 'group008', NOW(), NOW()),
  ('506ae84b-58ba-5a81-9958-77ebe03eca7f', 'group009', NOW(), NOW()),
  ('c650aba7-5b91-5ef8-af1e-296bfb7edeb9', 'group010', NOW(), NOW()),
  ('bcd2b776-dfd5-54c7-a840-e001a90a6463', 'group011', NOW(), NOW()),
  ('3c281162-93a0-58de-b104-4b78d68131c2', 'group012', NOW(), NOW()),
  ('8c70d79c-6896-55a1-a6d0-3d0fd7200d36', 'group013', NOW(), NOW()),
  ('ccc4fd93-5422-5c96-b07c-b6a06f905e34', 'group014', NOW(), NOW()),
  ('3840d638-ea66-5240-a3a0-85a6b7b9a891', 'group015', NOW(), NOW()),
  ('9bad7a44-ca2c-5868-940e-8f9402d5dafb', 'group016', NOW(), NOW()),
  ('ac3f0a79-1924-5815-8ba6-d9e30394babd', 'group017', NOW(), NOW()),
  ('0790a6a9-f7ee-5403-9e3b-aa2b89a56f09', 'group018', NOW(), NOW()),
  ('97cac5f0-0ce1-59e0-ade1-4b8416d279ea', 'group019', NOW(), NOW()),
  ('98b97945-169e-50b1-b343-787b10d3bf25', 'group020', NOW(), NOW()),
  ('5f6a687e-a186-5178-a89f-e7c72f7eb58c', 'group021', NOW(), NOW()),
  ('bfd8f040-a307-5020-ab8e-78975e607e5d', 'group022', NOW(), NOW()),
  ('27446636-760d-5a29-98f2-26996137db8d', 'group023', NOW(), NOW()),
  ('3dc7e70b-7913-50ed-bd73-c85574c7d677', 'group024', NOW(), NOW()),
  ('4e9097d7-3d2e-5242-a94b-6309d4aa3160', 'group025', NOW(), NOW()),
  ('c95446f4-6760-59ca-9074-efc405312b93', 'group026', NOW(), NOW()),
  ('e90a2155-63b5-5fac-b520-9a7cc87b371a', 'group027', NOW(), NOW()),
  ('0e8ae533-d62c-582b-8f7b-b5d3df10ed29', 'group028', NOW(), NOW()),
  ('b6676809-6205-5d43-83ae-5ec44090feaa', 'group029', NOW(), NOW()),
  ('603f929c-ab38-557b-abe9-b5d41e88c4b3', 'group030', NOW(), NOW()),
  ('957d639f-c79b-5558-8d54-a0dab37ea299', 'group031', NOW(), NOW()),
  ('7e4f6ac1-5746-5fec-866b-e69b70d2e51d', 'group032', NOW(), NOW()),
  ('fed3e4e7-a88b-5455-a9e3-e2d82284349b', 'group033', NOW(), NOW()),
  ('a4fcc95f-273f-54db-900a-0dda527ef733', 'group034', NOW(), NOW()),
  ('052d50c4-327a-5d6b-8d66-7e8e49778823', 'group035', NOW(), NOW()),
  ('7f5a639a-fbf5-5608-b554-ef7b541cdc83', 'group036', NOW(), NOW()),
  ('1b11336c-c0d8-55a2-91dc-d9876e079495', 'group037', NOW(), NOW()),
  ('81e7f254-772f-552a-90bb-861a530f6eb7', 'group038', NOW(), NOW()),
  ('3d1527c5-873e-59b4-b2f9-090212abf1ef', 'group039', NOW(), NOW()),
  ('9bbfc498-8952-50ae-8f94-a5cdc861d9fa', 'group040', NOW(), NOW()),
  ('d60f3e38-e640-5d1c-a6d2-0208c88ccc9d', 'group041', NOW(), NOW()),
  ('18a038fa-4e7e-5968-9d5f-fc908b55ab47', 'group042', NOW(), NOW()),
  ('ad590257-83bd-5e4a-bd6b-43ca854bacf3', 'group043', NOW(), NOW()),
  ('6bc6c74a-a2db-595f-ad99-1bf0a35d56aa', 'group044', NOW(), NOW()),
  ('a78eacd6-a4c3-5115-a774-2057fcab3332', 'group045', NOW(), NOW()),
  ('76b4ab65-1acb-5384-bc1d-479b5789ccef', 'group046', NOW(), NOW()),
  ('ea4e869f-d78a-5576-b8cd-41f88da8ad7c', 'group047', NOW(), NOW()),
  ('a4f75168-0a97-5838-bd46-7dc2b7fa8a07', 'group048', NOW(), NOW()),
  ('95f27b7c-103c-584c-8d90-3884ed3606bc', 'group049', NOW(), NOW()),
  ('8f477837-1a47-54cd-9866-9987999e0b84', 'group050', NOW(), NOW()),
  ('de40fb99-8d00-5b04-ab2c-d34264deba80', 'group051', NOW(), NOW()),
  ('bea9f714-a021-5391-9f45-8d1b436e5a46', 'group052', NOW(), NOW()),
  ('29eaecd8-389e-5d2a-97d4-d69bb3c36263', 'group053', NOW(), NOW()),
  ('4438f26a-b1e3-5c78-8092-6250201c711c', 'group054', NOW(), NOW()),
  ('e286e59e-52c6-5633-b7a4-489ad7a14cfa', 'group055', NOW(), NOW()),
  ('c1cdb5e5-c255-5cce-a6f8-e342a6520132', 'group056', NOW(), NOW()),
  ('69e0141c-b7bd-5113-b61c-db614f1f1de9', 'group057', NOW(), NOW()),
  ('008805f3-3fef-59de-9e93-96328741afbf', 'group058', NOW(), NOW()),
  ('297d9966-47bc-519f-834e-076ee3e7721c', 'group059', NOW(), NOW()),
  ('f1e2f9c0-f23e-5f33-8edb-9faf85c6bbfb', 'group060', NOW(), NOW()),
  ('0d448c39-3cc1-5ece-b884-57d11a751669', 'group061', NOW(), NOW()),
  ('be073412-89aa-5782-a7d0-d6b6aa7c58f1', 'group062', NOW(), NOW()),
  ('609fef8d-24ab-536e-815e-46eb8e203993', 'group063', NOW(), NOW()),
  ('63833830-854f-5b4e-9de8-b7e925bd6ad7', 'group064', NOW(), NOW()),
  ('d7956cee-6241-53f7-9221-0aa541bd86bd', 'group065', NOW(), NOW()),
  ('2a91149e-a013-54b2-979a-69e6aabd42a0', 'group066', NOW(), NOW()),
  ('d1261c43-30ae-539f-b05e-e4e1e0c7e264', 'group067', NOW(), NOW()),
  ('cd447537-963c-55a8-a953-a9d0e96e5f8e', 'group068', NOW(), NOW()),
  ('c62a7a6c-70cd-5c60-90b6-b677f6f6e20d', 'group069', NOW(), NOW()),
  ('62ba514e-0bfc-5f70-8882-39bd323a0923', 'group070', NOW(), NOW()),
  ('d560804c-905f-5bb5-b497-eb83580ddf58', 'group071', NOW(), NOW()),
  ('fb1f0d1b-861d-592a-b891-c7d480fb47b6', 'group072', NOW(), NOW()),
  ('74bf0a5c-88ac-5087-84a1-08adf3652405', 'group073', NOW(), NOW()),
  ('048afbf7-e592-50e0-b0b0-ba6768efab68', 'group074', NOW(), NOW()),
  ('c99b8d89-33b7-5b8c-a849-ad5df1180c0d', 'group075', NOW(), NOW()),
  ('8da14123-6a02-5bea-bbb6-98d550cae47b', 'group076', NOW(), NOW()),
  ('5e33a7e1-e95f-5acc-a4e4-e5da59596ebf', 'group077', NOW(), NOW()),
  ('f15732aa-8314-57d8-85aa-0140d5c3d3aa', 'group078', NOW(), NOW()),
  ('6b7e7c2d-22d6-582a-968a-091654fea4e8', 'group079', NOW(), NOW()),
  ('7080d780-0288-5199-8b7e-4cda8470c370', 'group080', NOW(), NOW()),
  ('87408f1a-e433-5f65-b445-7d2eb9473164', 'group081', NOW(), NOW()),
  ('25505b4e-3cdb-550e-b8dd-9e199f01d5b4', 'group082', NOW(), NOW()),
  ('8f59c9e4-2832-53f9-809c-72c9d5184a91', 'group083', NOW(), NOW()),
  ('f89edc59-95b9-500a-b3c5-e73084fa73f3', 'group084', NOW(), NOW()),
  ('08030e82-bc84-52f5-b47e-ead5ddfa8562', 'group085', NOW(), NOW()),
  ('6ee41e07-ca29-5a3f-9be1-ee99e047bcf8', 'group086', NOW(), NOW()),
  ('9a93834b-7fb8-5e80-a05f-dea32a1c5a45', 'group087', NOW(), NOW()),
  ('18f4905e-6aa1-59c4-94c0-14f97c3fcdda', 'group088', NOW(), NOW()),
  ('006c29fd-0389-5284-826c-434d8fc92103', 'group089', NOW(), NOW()),
  ('c7672c29-7ff1-5f56-af2d-0d1d74895f3d', 'group090', NOW(), NOW()),
  ('d97eaf94-4535-55f7-95f6-49c90423127a', 'group091', NOW(), NOW()),
  ('69893a4c-50e9-5f8f-af38-74289c4e6031', 'group092', NOW(), NOW()),
  ('0c637eb6-794f-551b-8286-d35c95bae698', 'group093', NOW(), NOW()),
  ('ae1c6b39-ed91-5021-af73-ac47bec3b6d9', 'group094', NOW(), NOW()),
  ('9ad477d3-0683-5f81-94c1-bac91e257c68', 'group095', NOW(), NOW()),
  ('747cba1a-90d7-5589-bd72-052bf7627e88', 'group096', NOW(), NOW()),
  ('6ca0f778-da7f-586d-80da-2443a2f6e3a0', 'group097', NOW(), NOW()),
  ('1b3d8003-3adf-5aed-bc1f-b87e5a11f187', 'group098', NOW(), NOW()),
  ('3be50426-6c64-5503-887a-495123703926', 'group099', NOW(), NOW()),
  ('a52bdb66-f374-5d5b-82df-64b1275ed424', 'group100', NOW(), NOW())
ON CONFLICT (uuid) DO UPDATE
  SET name = EXCLUDED.name, updated_at = NOW();
