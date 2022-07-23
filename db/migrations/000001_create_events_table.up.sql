CREATE TABLE events (
  parent_id VARCHAR (50),
  aggregate_id VARCHAR (50),
  version int,
  event_type VARCHAR (50),
  payload jsonb NOT NULL,
  created_at timestamp without time zone NOT NULL,

  PRIMARY KEY (aggregate_id, version)
);

CREATE INDEX idx_events_on_parent_id_aggregate_id_version ON events (parent_id, aggregate_id, version)
