# urbanism — Field Notes by AIA Brooklyn Urbanism

A place for neighbors to record feedback on urbanism sessions as they happen,
catch up on sessions they missed, and add their point of view. The record
becomes a resource for community-driven placemaking.

Production: https://aiabrooklynurbanism.com

- `prototype/` reference UI (immutable; the app renders it verbatim)
- `frontend/` Vite + TypeScript shell around the prototype bundle → Vercel
- `backend/` Go API (chi · pgx · PostgreSQL), three-tier → Railway
- `specs/` spec-kit artifacts · `CLAUDE.md` agent brief · `handoff.md` session log

See `CLAUDE.md` for local run and deploy steps.
