CREATE TABLE accounts (
 email TEXT PRIMARY KEY,
 role TEXT NOT NULL CHECK (role IN ('board','super_admin')),
 active BOOLEAN NOT NULL DEFAULT false,
 password_hash TEXT NOT NULL DEFAULT '',
 invite_hash TEXT NOT NULL DEFAULT '',
 invite_expires TIMESTAMPTZ,
 created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE account_audit (
 id BIGSERIAL PRIMARY KEY,
 actor TEXT NOT NULL,
 target TEXT NOT NULL,
 action TEXT NOT NULL,
 created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
