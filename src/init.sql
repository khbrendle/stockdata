CREATE DATABASE stocks;
\c stocks

CREATE SCHEMA symbolstream;
SET search_path TO symbolstream;

CREATE TABLE data (
  id VARCHAR(255) PRIMARY KEY,
  date_time TIMESTAMP,
  symbol VARCHAR(20) NOT NULL,
  response JSON,
  error TEXT
);

CREATE USER webapp WITH ENCRYPTED PASSWORD 'webapp';
GRANT INSERT ON ALL TABLES IN SCHEMA symbolstream TO webapp;
GRANT USAGE ON SCHEMA symbolstream TO webapp;
ALTER ROLE webapp SET search_path TO symbolstream;
