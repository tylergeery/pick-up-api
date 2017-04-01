CREATE TABLE user_tokens (
    user_id integer PRIMARY KEY,
    refresh_token VARCHAR(250),
    access_token VARCHAR(250),
    created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'PST'),
    updated_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'PST')
);

INSERT INTO user_tokens (user_id, refresh_token, created_at)
VALUES (1,'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwZXJtaXNzaW9ucyI6MSwidG9rZW5UeXBlIjoicmVmcmVzaCIsInRzIjoxNDkwNTkyMDcyLCJ0eXBlIjoidXNlciIsInVzZXJJZCI6MH0.6fnQTQ9iqdp-3duuL4A_WQ6tiyHVUvmY3bKYx1VFIms', NOW());
