CREATE TABLE users
(
    id      SERIAL PRIMARY KEY,
    user_id UUID UNIQUE NOT NULL,
    IP      VARCHAR(15) NOT NULL,
    email   VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE users_tokens
(
    id            SERIAL PRIMARY KEY,
    user_id       UUID NOT NULL UNIQUE ,
    refresh_token_hashed TEXT NOT NULL,
    issued_at     BIGINT NOT NULL,
    expires_at BIGINT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE,
    CONSTRAINT unique_refresh_token UNIQUE (user_id, refresh_token_hashed)
);

INSERT INTO users (user_id, IP, email)
VALUES
    ('123e4567-e89b-12d3-a456-426614174001', '192.168.1.1', 'user1@example.com'),
    ('123e4567-e89b-12d3-a456-426614174002', '192.168.1.2', 'user2@example.com'),
    ('123e4567-e89b-12d3-a456-426614174003', '192.168.1.3', 'user3@example.com'),
    ('123e4567-e89b-12d3-a456-426614174004', '192.168.1.4', 'user4@example.com'),
    ('123e4567-e89b-12d3-a456-426614174005', '192.168.1.5', 'user5@example.com'),
    ('123e4567-e89b-12d3-a456-426614174006', '192.168.1.6', 'user6@example.com'),
    ('123e4567-e89b-12d3-a456-426614174007', '192.168.1.7', 'user7@example.com'),
    ('123e4567-e89b-12d3-a456-426614174008', '192.168.1.8', 'user8@example.com'),
    ('123e4567-e89b-12d3-a456-426614174009', '192.168.1.9', 'user9@example.com'),
    ('123e4567-e89b-12d3-a456-426614174010', '192.168.1.10', 'user10@example.com');
