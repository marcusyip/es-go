CREATE TABLE events (
  aggregate_id VARCHAR (50),
  version int,
  event_type VARCHAR (50),
  payload jsonb NOT NULL,
  created_at timestamp without time zone NOT NULL,

  PRIMARY KEY (aggregate_id, version)
);
