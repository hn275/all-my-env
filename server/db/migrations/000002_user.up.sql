CREATE TYPE vendors AS ENUM('github');

CREATE TABLE users (
    id INT PRIMARY KEY NOT NULL,
    created_at timestamp NOT NULL,
    vendor vendors,
    username VARCHAR(255) NOT NULL
);
