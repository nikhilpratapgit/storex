BEGIN;

CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TYPE asset_type AS ENUM (
    'laptop',
    'keyboard',
    'mouse',
    'mobile'
);

CREATE TYPE asset_status AS ENUM (
    'available',
    'assigned',
    'in_service',
    'for_repair',
    'damaged'
);

CREATE TYPE user_role AS ENUM (
    'admin',
    'employee',
    'project-manager',
    'asset-manager',
    'employee-manager'
);

CREATE TYPE user_type AS ENUM (
    'full-time',
    'intern',
    'freelancer'
);

CREATE TYPE owner_type AS ENUM (
    'client',
    'remotestate'
);

CREATE TYPE connection_type AS ENUM (
    'wired',
    'wireless'
);

CREATE TABLE IF NOT EXISTS users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name          TEXT        NOT NULL,
    email         TEXT        NOT NULL,
    role          user_role   DEFAULT 'employee',
    type          user_type   NOT NULL,
    phone_no      TEXT        NOT NULL,
    password      TEXT        NOT NULL,
    created_at    TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    archived_at   TIMESTAMPTZ
    );

CREATE UNIQUE INDEX idx_unique_email
    ON users (email)
    WHERE archived_at IS NULL;

CREATE TABLE assets (
                        id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        brand           TEXT          NOT NULL,
                        model           TEXT          NOT NULL,
                        serial_number       TEXT UNIQUE   NOT NULL,
                        type            asset_type    NOT NULL,
                        status          asset_status  DEFAULT 'available',
                        owner           owner_type    NOT NULL,

                        assigned_by_id  UUID REFERENCES users(id),
                        assigned_to     UUID REFERENCES users(id),
                        assigned_on     TIMESTAMPTZ,

                        warranty_start  TIMESTAMPTZ   NOT NULL,
                        warranty_end    TIMESTAMPTZ   NOT NULL,

                        service_start   TIMESTAMPTZ,
                        service_end     TIMESTAMPTZ,
                        returned_on     TIMESTAMPTZ,

                        created_at      TIMESTAMPTZ   DEFAULT now(),
                        updated_at      TIMESTAMPTZ,
                        archived_at     TIMESTAMPTZ,
                        archived_by     UUID REFERENCES users(id)
);

CREATE INDEX idx_asset_id
    ON assets(id);

CREATE INDEX idx_assigned_to
    ON assets(assigned_to);

CREATE TABLE IF NOT EXISTS user_session (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    archived_at TIMESTAMPTZ
    );

CREATE TABLE laptops (
                        id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        asset_id         UUID UNIQUE REFERENCES assets(id),
                        processor        TEXT,
                        ram              TEXT,
                        storage          TEXT,
                        operating_system TEXT,
                        charger          TEXT,
                        device_password  TEXT NOT NULL
);

CREATE TABLE keyboards (
                          id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                          asset_id     UUID UNIQUE REFERENCES assets(id),
                          layout       TEXT,
                          connectivity connection_type
);

CREATE TABLE mouses (
                       id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       asset_id     UUID UNIQUE REFERENCES assets(id),
                       dpi          INT,
                       connectivity connection_type
);

CREATE TABLE mobiles (
                        id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        asset_id         UUID UNIQUE REFERENCES assets(id),
                        operating_system               TEXT NOT NULL,
                        ram              TEXT NOT NULL,
                        storage          TEXT NOT NULL,
                        charger          TEXT,
                        device_password  TEXT NOT NULL
);

COMMIT;