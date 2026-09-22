CREATE TABLE IF NOT EXISTS product (
                                       id            BIGSERIAL    NOT NULL UNIQUE,
                                       guid          UUID         NOT NULL PRIMARY KEY,
                                       name          VARCHAR(255) NOT NULL UNIQUE,
    description   TEXT,
    price         BIGINT       NOT NULL CHECK (price > 0),
    category_guid UUID         NOT NULL REFERENCES category(guid) ON DELETE RESTRICT,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS idx_product_category_guid ON product (category_guid);