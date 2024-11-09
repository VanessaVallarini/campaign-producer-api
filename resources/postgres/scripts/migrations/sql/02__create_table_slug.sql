CREATE TABLE IF NOT EXISTS slug(
    id  UUID PRIMARY KEY NOT NULL,
    name   VARCHAR(50) UNIQUE NOT NULL,
    status VARCHAR(10) DEFAULT 'ACTIVE',
    cost DECIMAL(5,2) NOT NULL,
    created_by VARCHAR(60),
    updated_by VARCHAR(60),
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX slug_id ON campaign_consumer_api.slug USING btree (id);
CREATE INDEX slug_name ON campaign_consumer_api.slug USING btree (name);