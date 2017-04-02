CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(50),
    password VARCHAR(60),
    name VARCHAR(50),
    facebook_id BIGINT,
    is_active SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'PST'),
    updated_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'PST')
);

INSERT INTO users (email, password, name, created_at)
VALUES ('tyler@geerydev.com', 'test', 'Tyler Geery', NOW());
