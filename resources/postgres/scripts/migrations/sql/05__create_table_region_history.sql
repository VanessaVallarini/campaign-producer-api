CREATE TABLE IF NOT EXISTS region_history(
    id  UUID PRIMARY KEY NOT NULL,
    region_id  UUID NOT NULL,
    status VARCHAR(10) DEFAULT 'ACTIVE',
    description VARCHAR(100),
    created_by VARCHAR(60),
    created_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (region_id) REFERENCES region(id)
);