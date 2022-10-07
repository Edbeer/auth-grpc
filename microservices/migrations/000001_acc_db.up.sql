DROP TABLE IF EXISTS users CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE users
(
    user_id      UUID PRIMARY KEY            DEFAULT uuid_generate_v4(),
    name         VARCHAR(32)                 NOT NULL check ( name <> '' ),
    email        VARCHAR(64)                 NOT NULL check ( email <> '' ),
    password     VARCHAR(250)                NOT NULL CHECK ( octet_length(password) <> 0 ),
    role         VARCHAR(10)                 NOT NULL DEFAULT 'user',
    created_at   TIMESTAMP                   NOT NULL DEFAULT now()
);