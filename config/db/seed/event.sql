CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    title VARCHAR(50),
    description text,
    owner_id INT,
    `date` VARCHAR(50),
    cost NUMERIC(2) DEFAULT 0.00
    is_active SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'PST'),
    updated_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'PST')
);
