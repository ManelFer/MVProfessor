CREATE TABLE users (
    id  SERIAL  PRIMARY KEY,
    name    VARCHAR(100)    NOT NULL,
    password_hash   VARCHAR(255)    NOT NULL,
    created_at  TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP   DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE users
ADD email VARCHAR(150) UNIQUE NOT NULL

CREATE INDEX idx_users_email ON users(email);

