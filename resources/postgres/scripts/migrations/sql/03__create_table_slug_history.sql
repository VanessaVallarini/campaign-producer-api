CREATE TABLE IF NOT EXISTS slug_history(
    id  UUID PRIMARY KEY NOT NULL,
    slug_id  UUID NOT NULL,
    status VARCHAR(10) DEFAULT 'ACTIVE',
    description VARCHAR(100),
    created_by VARCHAR(60),
    created_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (slug_id) REFERENCES slug(id)
);