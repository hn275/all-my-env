CREATE TABLE repos (
    id BIGSERIAL PRIMARY KEY NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL,
    name VARCHAR(255) NOT NULL,
    repo_url VARCHAR(255) NOT NULL,
    owner_id BIGSERIAL NOT NULL
);

CREATE TABLE variables (
    id BIGSERIAL PRIMARY KEY NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    key VARCHAR(100) NOT NULL,
    value TEXT NOT NULL,
    nonce TEXT NOT NULL,
    repo_id BIGSERIAL NOT NULL,
    CONSTRAINT fk_repos
        FOREIGN KEY(repo_id)
        REFERENCES repos(id)
);
