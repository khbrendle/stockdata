CREATE DATABASE stocks;
\c stocks

CREATE SCHEMA symbolstream;
SET search_path TO symbolstream;

CREATE TABLE response_raw (
  id VARCHAR(255) PRIMARY KEY,
  date_time TIMESTAMP,
  symbol VARCHAR(20) NOT NULL,
  response JSON,
  error TEXT
);

CREATE TABLE message (
  response_raw_id VARCHAR(255) NOT NULL,
  symbol VARCHAR(20) NOT NULL,
  message_id BIGINT,
  message JSON,
  CONSTRAINT message_PK PRIMARY KEY (response_raw_id, message_id)
);

-- will need trigger on response_raw

-- CREATE TABLE message_detail (
--   symbol VARCHAR(20) NOT NULL,
--   message_id BIGINT,
--   message JSON,
--   CONSTRAINT message_PK PRIMARY KEY (response_raw_id, message_id)
-- );



CREATE USER webapp WITH ENCRYPTED PASSWORD 'webapp';
GRANT INSERT ON ALL TABLES IN SCHEMA symbolstream TO webapp;
GRANT USAGE ON SCHEMA symbolstream TO webapp;
ALTER ROLE webapp SET search_path TO symbolstream;
