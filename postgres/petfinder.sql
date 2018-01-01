CREATE TABLE IF NOT EXISTS breed (
  id serial,
  name text,
  animal text,
  created_on timestamp with time zone,
  updated_on timestamp with time zone,
  PRIMARY KEY (id)
);

CREATE INDEX on breed (animal);

CREATE TABLE IF NOT EXISTS shelter (
  id text,
  name text,
  longitude text,
  latitude text,
  address1 text,
  address2 text,
  city text,
  state text,
  country text,
  phone text,
  email text,
  zip text,
  fax text,
  created_on timestamp with time zone,
  updated_on timestamp with time zone,
  PRIMARY KEY (id)
);

CREATE INDEX on shelter (city);
CREATE INDEX on shelter (state);
CREATE INDEX on shelter (zip);

CREATE TABLE IF NOT EXISTS pet (
  id text,
  shelter_id text REFERENCES shelter (id),
  shelterpet_id text,
  status text,
  name text,
  sex text,
  age text,
  size text,
  mix text,
  address1 text,
  address2 text,
  city text,
  state text,
  phone text,
  email text,
  zip text,
  fax text,
  description text,
  last_update timestamp with time zone,
  created_on timestamp with time zone,
  updated_on timestamp with time zone,
  PRIMARY KEY (id)
);

CREATE INDEX on shelter (city);
CREATE INDEX on shelter (state);
CREATE INDEX on shelter (zip);

CREATE TABLE IF NOT EXISTS pet_option (
  id serial,
  pet_id text REFERENCES pet (id),
  option text,
  PRIMARY KEY (id)
);

CREATE INDEX on pet_option (pet_id);

CREATE TABLE IF NOT EXISTS pet_breed (
  id serial,
  breed_id integer REFERENCES breed (id),
  pet_id text REFERENCES pet (id),
  PRIMARY KEY (id)
);

CREATE INDEX on pet_breed (breed_id);
CREATE INDEX on pet_breed (pet_id);

CREATE TABLE IF NOT EXISTS pet_photo (
  id serial,
  photo_id text,
  pet_id text REFERENCES pet (id),
  url text,
  size text,
  PRIMARY KEY (id)
);

CREATE INDEX on pet_photo (pet_id);
CREATE INDEX on pet_photo (size);

