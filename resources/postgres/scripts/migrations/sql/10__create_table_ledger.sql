CREATE TABLE IF NOT EXISTS ledger (
    id  UUID PRIMARY KEY NOT NULL,
    spent_id  UUID NOT NULL,
    campaign_id  UUID NOT NULL,
    merchant_id  UUID NOT NULL,
    slug_name VARCHAR(50) NOT NULL,
    region_name VARCHAR(50) NOT NULL,
    user_id UUID NOT NULL,
    session_id UUID NOT NULL,
    event_type VARCHAR(30) NOT NULL,
    cost DECIMAL(5,2) NOT NULL,
    ip VARCHAR(30) NOT NULL,
    lat FLOAT NOT NULL,
    long FLOAT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    event_time TIMESTAMP WITH TIME ZONE NOT NULL,
    FOREIGN KEY (spent_id) REFERENCES spent(id),
    FOREIGN KEY (campaign_id) REFERENCES campaign(id),
    FOREIGN KEY (merchant_id) REFERENCES merchant(id),
    FOREIGN KEY (slug_name) REFERENCES slug(name),
    FOREIGN KEY (region_name) REFERENCES region(name)
);

CREATE INDEX ledger_spent_id ON campaign_consumer_api.ledger USING btree (spent_id);
CREATE INDEX ledger_campaign_id ON campaign_consumer_api.ledger USING btree (campaign_id);
CREATE INDEX ledger_merchant_id ON campaign_consumer_api.ledger USING btree (merchant_id);
CREATE INDEX ledger_slug_name ON campaign_consumer_api.ledger USING btree (slug_name);