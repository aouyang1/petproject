CREATE TABLE IF NOT EXISTS breed (
  id serial,
  name text,
  animal text,
  PRIMARY KEY (id)
);

CREATE INDEX on breed (animal);
