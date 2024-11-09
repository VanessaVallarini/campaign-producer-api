CREATE TABLE IF NOT EXISTS campaign(
    id  UUID PRIMARY KEY NOT NULL,
    merchant_id  UUID NOT NULL,
    status VARCHAR(10) DEFAULT 'ACTIVE',
    budget DECIMAL(5,2) NOT NULL,
    created_by VARCHAR(60),
    updated_by VARCHAR(60),
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (merchant_id) REFERENCES merchant(id)
);

CREATE INDEX campaign_id ON campaign_consumer_api.campaign USING btree (id);
