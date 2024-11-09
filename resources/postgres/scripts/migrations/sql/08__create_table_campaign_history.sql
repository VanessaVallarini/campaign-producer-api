CREATE TABLE IF NOT EXISTS campaign_history(
    id  UUID PRIMARY KEY NOT NULL,
    campaign_id  UUID NOT NULL,
    status VARCHAR(10) DEFAULT 'ACTIVE',
    description VARCHAR(100),
    created_by VARCHAR(60),
    created_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (campaign_id) REFERENCES campaign(id)
);