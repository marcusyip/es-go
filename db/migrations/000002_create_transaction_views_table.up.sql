CREATE TABLE IF NOT EXISTS transaction_views (
  id varchar(255) NOT NULL,
  status varchar(32) NOT NULL,
  currency varchar(16) NOT NULL,
  amount numeric(36, 18) NOT NULL,
  done_by varchar(64) NOT NULL,
  created_at timestamp without time zone NOT NULL,
  updated_at timestamp without time zone NOT NULL,

  PRIMARY KEY (id)
);
