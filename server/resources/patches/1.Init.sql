-- create table user
CREATE TABLE IF NOT EXISTS user (
    id int(10) unsigned NOT NULL AUTO_INCREMENT,
    client_public_key_hash BINARY(20) NOT NULL,
    client_public_key TEXT,
    server_private_key TEXT,
    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY client_public_key_hash (client_public_key_hash)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- create table account
CREATE TABLE IF NOT EXISTS account (
    id int(10) unsigned NOT NULL AUTO_INCREMENT,
    user_id int(10) unsigned DEFAULT NULL,
    ext_id VARCHAR(255) NOT NULL,
    provider int(10) unsigned NOT NULL,
    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_account_user FOREIGN KEY (user_id) REFERENCES user(id),
    UNIQUE KEY ext_id (ext_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;