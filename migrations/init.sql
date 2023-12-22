BEGIN;
CREATE SCHEMA url_shortener;
CREATE TABLE url_shortener.urls(
                           short_url VARCHAR(10) PRIMARY KEY,
                           original_url VARCHAR(255) UNIQUE NOT NULL
);
COMMIT;
