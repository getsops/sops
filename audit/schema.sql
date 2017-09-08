CREATE TABLE decrypt_event (
  id SERIAL PRIMARY KEY,
  timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  username TEXT,
  file TEXT
);
