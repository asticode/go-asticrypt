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

-- create table email
CREATE TABLE IF NOT EXISTS email (
    id int(10) unsigned NOT NULL AUTO_INCREMENT,
    user_id int(10) unsigned DEFAULT NULL,
    addr VARCHAR(255),
    validation_token CHAR(100) NOT NULL,
    validated_at datetime DEFAULT NULL,
    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_email_user FOREIGN KEY (user_id) REFERENCES user(id),
    UNIQUE KEY addr (addr),
    KEY validation_token(validation_token)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;