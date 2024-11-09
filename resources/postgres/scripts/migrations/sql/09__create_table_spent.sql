CREATE TABLE IF NOT EXISTS spent(
    id  UUID PRIMARY KEY NOT NULL,
    campaign_id  UUID NOT NULL,
    merchant_id  UUID NOT NULL,
    bucket VARCHAR(50) NOT NULL,
    total_spent DECIMAL(14,2) NOT NULL,
    total_clicks integer NOT NULL,
    total_impressions integer NOT NULL,
    FOREIGN KEY (campaign_id) REFERENCES campaign(id),
    FOREIGN KEY (merchant_id) REFERENCES merchant(id)
);

CREATE INDEX spent_merchant ON campaign_consumer_api.spent USING btree (merchant_id, bucket);

ALTER TABLE campaign_consumer_api.spent
ADD CONSTRAINT unique_merchant_bucket
UNIQUE (merchant_id, bucket);