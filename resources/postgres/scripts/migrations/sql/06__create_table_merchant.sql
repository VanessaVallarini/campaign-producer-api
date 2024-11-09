CREATE TABLE IF NOT EXISTS merchant (
    id  UUID PRIMARY KEY NOT NULL,
    owner_id  UUID NOT NULL,
    region_id  UUID NOT NULL,
    slugs UUID[] NOT NULL,
    name   VARCHAR(50) NOT NULL,
    status VARCHAR(10) DEFAULT 'ACTIVE',
    created_by VARCHAR(60),
    updated_by VARCHAR(60),
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (owner_id) REFERENCES owner(id),
    FOREIGN KEY (region_id) REFERENCES region(id)
);

CREATE INDEX merchant_id ON campaign_consumer_api.merchant USING btree (id);
CREATE INDEX merchant_name ON campaign_consumer_api.merchant USING btree (name);