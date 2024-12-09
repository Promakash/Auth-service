CREATE TABLE users
(
    id      SERIAL PRIMARY KEY,
    user_id UUID UNIQUE NOT NULL,
    IP      VARCHAR(15) UNIQUE NOT NULL,
    email   VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE users_tokens
(
    id            SERIAL PRIMARY KEY,
    user_id       UUID NOT NULL,
    refresh_token_hashed TEXT NOT NULL,
    issued_at     TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE,
    CONSTRAINT unique_refresh_token UNIQUE (user_id, refresh_token_hashed)
);
