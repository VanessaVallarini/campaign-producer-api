CREATE SCHEMA campaign_consumer_api;

GRANT ALL PRIVILEGES ON DATABASE "campaign-consumer-api" TO "postgres";

GRANT USAGE ON SCHEMA campaign_consumer_api TO "postgres";
ALTER USER "postgres" SET search_path = 'campaign_consumer_api';


SET SCHEMA 'campaign_consumer_api';
ALTER DEFAULT PRIVILEGES
    IN SCHEMA campaign_consumer_api
GRANT SELECT, UPDATE, INSERT, DELETE ON TABLES
    TO "postgres";

ALTER DEFAULT PRIVILEGES
    IN SCHEMA campaign_consumer_api
GRANT USAGE ON SEQUENCES
    TO "postgres";