DROP TABLE IF EXISTS users;
CREATE TABLE users (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    cc VARCHAR(255) NOT NULL,
    age VARCHAR(255) NOT NULL,
    birth_date VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    address VARCHAR(255) NOT NULL,
    suburb VARCHAR(255) NOT NULL,
    voting_place VARCHAR(255) NOT NULL,
    civil_status VARCHAR(255) NOT NULL,
    phone VARCHAR(255) NOT NULL,
    ecan Bool NOT NULL DEFAULT false
);
CREATE TABLE children (
    id VARCHAR(32) PRIMARY KEY,
    user_id VARCHAR(32) NOT NULL,
    name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    age VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE services (
    id VARCHAR(32) PRIMARY KEY,
    user_id VARCHAR(32) NOT NULL,
    service VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
