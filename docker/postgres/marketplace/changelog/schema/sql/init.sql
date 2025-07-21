CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL PRIMARY KEY,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE TABLE ads
(
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(100) NOT NULL,
    description TEXT         NOT NULL,
    image_url   TEXT         NOT NULL,
    price       BIGINT       NOT NULL CHECK (price > 0),
    author_id   INT          NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_ads_created_at ON ads(created_at DESC);

CREATE INDEX idx_ads_price ON ads(price DESC);