-- Add urls table

CREATE TABLE urls(
  id serial PRIMARY KEY,
  original VARCHAR UNIQUE NOT NULL,
  shorten VARCHAR ( 7 ) UNIQUE,
  expired_at TIMESTAMP NOT NULL,
  created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX expired_at_index ON urls (expired_at);

---- create above / drop below ----

DROP TABLE urls;
