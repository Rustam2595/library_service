CREATE TABLE IF NOT EXISTS Users (
    uid VARCHAR(36) PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    pass TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON Users (email);

CREATE TABLE IF NOT EXISTS Books(
    bid VARCHAR(36) PRIMARY KEY,
    label TEXT NOT NULL,
    author TEXT NOT NULL,
    deleted BOOLEAN NOT NULL DEFAULT false,
    user_uid VARCHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    CONSTRAINT fk_books_user FOREIGN KEY (user_uid) REFERENCES Users(uid) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_books_user_uid ON Books (user_uid); --user_uid будет уникальным