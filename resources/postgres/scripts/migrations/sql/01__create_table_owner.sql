CREATE TABLE IF NOT EXISTS owner (
    id  UUID PRIMARY KEY NOT NULL,
    email VARCHAR(50) UNIQUE NOT NULL,
    status VARCHAR(10) DEFAULT 'ACTIVE',
    created_by VARCHAR(60),
    updated_by VARCHAR(60),
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX owner_id ON campaign_consumer_api.owner USING btree (id);
CREATE INDEX owner_email ON campaign_consumer_api.owner USING btree (email);