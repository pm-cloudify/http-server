CREATE TABLE jobs (
    id SERIAL PRIMARY KEY,
    file_id INT NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'created' CHECK (status IN('created', 'executing', 'failed', 'done')),
    query VARCHAR(512),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_change TIMESTAMP NOT NULL DEFAULT NOW()
);