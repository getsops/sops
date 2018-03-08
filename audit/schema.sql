CREATE TYPE action AS ENUM ('decrypt', 'encrypt', 'rotate');
CREATE TABLE audit_event (
  id SERIAL PRIMARY KEY,
  timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  action action,
  file TEXT,
  username TEXT,
  details jsonb
);


CREATE ROLE sops WITH NOSUPERUSER INHERIT NOCREATEROLE NOCREATEDB LOGIN PASSWORD 'sops';

GRANT INSERT ON audit_event TO sops;
GRANT USAGE ON audit_event_id_seq TO sops;
