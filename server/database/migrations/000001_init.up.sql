CREATE TABLE users (
    id INT NOT NULL PRIMARY KEY,
    login CHAR(255) NOT NULL,
    email CHAR(255) NOT NULL,
    refresh_token CHAR(32) DEFAULT NULL
);

CREATE TABLE repositories (
    id INT NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    full_name  VARCHAR(255) NOT NULL,
    user_id INT NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE variables (
    id INT NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NULL,
    variable_key CHAR(255) NOT NULL,
    variable_value VARCHAR(255) NOT NULL,
    repository_id INT NOT NULL,
    CONSTRAINT fk_repo_var FOREIGN KEY (repository_id) REFERENCES repositories(id) ON DELETE CASCADE,
    CONSTRAINT unique_key_repo UNIQUE KEY (repository_id, variable_key)
);

CREATE TABLE permissions (
    id INT NOT NULL PRIMARY KEY,
    repository_id INT NOT NULL,
    user_id INT NOT NULL,
    CONSTRAINT fk_repo_perm FOREIGN KEY (repository_id) REFERENCES repositories(id) ON DELETE CASCADE,
    CONSTRAINT unique_permissions UNIQUE KEY (repository_id, user_id)
);
