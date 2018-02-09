CREATE TABLE decrypt_event (
  id SERIAL PRIMARY KEY,
  timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  username TEXT,
  file TEXT
);

CREATE RULE decrypt_event_delete_protection AS ON DELETE TO decrypt_event DO INSTEAD NOTHING;

CREATE TABLE encrypt_event (
  id SERIAL PRIMARY KEY,
  timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  username TEXT,
  file TEXT
);

CREATE RULE encrypt_event_delete_protection AS ON DELETE TO encrypt_event DO INSTEAD NOTHING;

CREATE TABLE rotate_event (
  id SERIAL PRIMARY KEY,
  timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  username TEXT,
  file TEXT
);

CREATE RULE rotate_event_delete_protection AS ON DELETE TO rotate_event DO INSTEAD NOTHING;

CREATE ROLE sops WITH NOSUPERUSER INHERIT NOCREATEROLE NOCREATEDB LOGIN PASSWORD 'sops';

GRANT INSERT ON decrypt_event TO sops;
GRANT USAGE ON decrypt_event_id_seq TO sops;
GRANT INSERT ON encrypt_event TO sops;
GRANT USAGE ON encrypt_event_id_seq TO sops;
GRANT INSERT ON rotate_event TO sops;
GRANT USAGE ON rotate_event_id_seq TO sops;
