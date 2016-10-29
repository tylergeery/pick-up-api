CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(50),
    password VARCHAR(60),
    name VARCHAR(50),
    facebook_id INT,
    is_active INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'PST'),
    updated_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'PST')
);

INSERT INTO users (email, password, name, created_at)
VALUES ('tyler.geery@yahoo.com', 'test', 'Tyler Geery', NOW());
