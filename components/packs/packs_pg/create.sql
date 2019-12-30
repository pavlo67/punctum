CREATE TABLE packs (
    id                BIGSERIAL                PRIMARY KEY,
    identity_key      VARCHAR(255)             NOT NULL,
    address_from      VARCHAR(255)             NOT NULL,
    address_to        TEXT                     ,
    options           TEXT                     ,
    type_key          VARCHAR(255)             NOT NULL,
    content           TEXT                     ,
    history           TEXT                     ,
    created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_packs_identity_key ON packs(identity_key);

CREATE INDEX idx_packs_type_key     ON packs(type_key);

CREATE INDEX idx_packs_created_at   ON packs(created_at);


