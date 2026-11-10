-- Field Notes schema v1. Mirrors the prototype's client document (version 1).
CREATE TABLE IF NOT EXISTS sessions (
  id            TEXT PRIMARY KEY,
  code          TEXT NOT NULL UNIQUE,
  title         TEXT NOT NULL,
  date          DATE NOT NULL,
  end_date      DATE NOT NULL,
  place         TEXT NOT NULL DEFAULT '',
  neighborhood  TEXT NOT NULL DEFAULT '',
  description   TEXT NOT NULL DEFAULT '',
  art           TEXT NOT NULL DEFAULT 'blocks',
  tags          JSONB NOT NULL DEFAULT '[]',
  notes         JSONB NOT NULL DEFAULT '[]',
  poll          TEXT NOT NULL DEFAULT '',
  options       JSONB NOT NULL DEFAULT '[]',
  sample        BOOLEAN NOT NULL DEFAULT FALSE,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS entries (
  id          TEXT PRIMARY KEY,
  session_id  TEXT NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
  type        TEXT NOT NULL CHECK (type IN ('text','voice','photo')),
  name        TEXT NOT NULL DEFAULT 'Anonymous neighbor',
  hood        TEXT NOT NULL DEFAULT '',
  text        TEXT NOT NULL,
  tags        JSONB NOT NULL DEFAULT '[]',
  location    TEXT NOT NULL DEFAULT '',
  art         TEXT NOT NULL DEFAULT '',
  created     TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS entries_session_idx ON entries(session_id, created);

CREATE TABLE IF NOT EXISTS replies (
  id        TEXT PRIMARY KEY,
  entry_id  TEXT NOT NULL REFERENCES entries(id) ON DELETE CASCADE,
  name      TEXT NOT NULL DEFAULT 'Anonymous neighbor',
  text      TEXT NOT NULL,
  created   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS replies_entry_idx ON replies(entry_id, created);

CREATE TABLE IF NOT EXISTS materials (
  id           TEXT PRIMARY KEY,
  session_id   TEXT NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
  title        TEXT NOT NULL,
  category     TEXT NOT NULL CHECK (category IN ('panel','slides','handout','report','other')),
  description  TEXT NOT NULL DEFAULT '',
  credited_to  TEXT NOT NULL DEFAULT '',
  file_name    TEXT NOT NULL,
  mime_type    TEXT NOT NULL,
  media_id     TEXT NOT NULL,
  size         BIGINT NOT NULL CHECK (size > 0),
  created      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS votes (
  client_id   TEXT NOT NULL,
  session_id  TEXT NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
  option_idx  INT NOT NULL,
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (client_id, session_id)
);

CREATE TABLE IF NOT EXISTS reactions (
  client_id  TEXT NOT NULL,
  entry_id   TEXT NOT NULL REFERENCES entries(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (client_id, entry_id)
);

CREATE TABLE IF NOT EXISTS journal_posts (
  id            TEXT PRIMARY KEY,
  number        TEXT NOT NULL,
  category      TEXT NOT NULL,
  title         TEXT NOT NULL,
  excerpt       TEXT NOT NULL,
  date          DATE NOT NULL,
  read_time     TEXT NOT NULL,
  art           TEXT NOT NULL,
  author        TEXT NOT NULL,
  session_code  TEXT NOT NULL DEFAULT '',
  body          JSONB NOT NULL DEFAULT '[]'
);

CREATE TABLE IF NOT EXISTS schema_migrations (
  version TEXT PRIMARY KEY,
  applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
