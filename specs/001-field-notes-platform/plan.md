# Implementation Plan: Field Notes Platform (001)

## Summary

Keep the prototype as the UI. Wrap it with a boot adapter that syncs its
localStorage document with a Go API. Persist shared slices in PostgreSQL.

## Technical Context

- Frontend: Vite shell + `src/sync.ts` adapter; `public/legacy/` = verbatim bundle.
- Backend: Go 1.27, chi, pgx; layers per constitution II.
- DB: PostgreSQL (Railway add-on; local via Postgres.app or docker-compose).
- Deploy: Vercel (root `frontend/`, `VITE_API_URL`, custom domain aiabrooklynurbanism.com),
  Railway (Dockerfile, `DATABASE_URL`, `ALLOWED_ORIGINS=https://aiabrooklynurbanism.com,https://www.aiabrooklynurbanism.com`).

## Architecture

```
frontend/
  index.html            shell (head verbatim from prototype)
  public/legacy/app.css prototype stylesheet (verbatim, fonts embedded)
  public/legacy/app.js  prototype React bundle (verbatim)
  src/main.ts           boot: sync then load bundle
  src/sync.ts           localStorage <-> API adapter (diff on setItem, poll 45s)
  src/api.ts            fetch client, X-Client-Id
backend/
  cmd/server            wiring + graceful shutdown
  internal/config       env
  internal/domain       entities, errors, repository interfaces
  internal/service      validation + rules
  internal/repository/postgres  pgx impl, embedded migrations + seed
  internal/controller/http      chi handlers, DTO = domain JSON
  internal/middleware   request id, logger, recover, CORS, client id, body limit
```

## Sync Contract

- Boot: `GET /api/v1/state` → merged into `aia-field-notes-v1` → bundle mounts.
- Writes: bundle calls `localStorage.setItem(key, doc)`; adapter diffs vs last
  doc and issues `POST /sessions|/entries|/replies|/materials`,
  `PUT /sessions/{id}/vote`, `PUT|DELETE /entries/{id}/like`.
- Refresh: poll every 45 s, write merged doc, dispatch `storage` event (bundle re-reads).

## Phases

1. **Phase 1 (done)**: scaffold, API, adapter, seed, deploy configs.
2. **Phase 2**: media bytes to object storage (Railway volume or S3-compatible),
   `/materials/{id}/file`; vote tallies surfaced; admin export (CSV/JSON) for thesis work.
3. **Phase 3**: componentize UI from prototype with visual regression tests
   (Playwright screenshots vs prototype) so parity is enforced by CI.

## Risks

- Bundle is minified; adapter relies on the `aia-field-notes-v1` key and the
  `storage` listener. Pin the bundle; any new prototype export must be re-verified.
- Client id is spoofable; acceptable for phase 1 (no accounts).
