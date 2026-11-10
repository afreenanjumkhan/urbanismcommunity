# CLAUDE.md — Field Notes (AIA Brooklyn Urbanism)

Standing brief for any Claude Code session. Read this, then `handoff.md`
(latest checkpoint), then `.specify/memory/constitution.md`.

## What this is
Community feedback platform: neighbors log field notes on urbanism sessions,
catch up on missed ones, reply, vote, share materials. Record feeds future
placemaking work. Owner: AIA Brooklyn Urbanism Committee.
Production target: https://aiabrooklynurbanism.com

## Hard rules (see constitution)
1. UI = `prototype/index.html`, served verbatim from `frontend/public/legacy/`. Never edit the bundle.
2. Backend layering: controller → service → repository; `domain` shared; SQL only in repository.
3. Frontend and backend are separate deployables (Vercel / Railway). Talk only via `/api/v1`.
4. Append-only record; seed rows carry `sample: true`.
5. Log a checkpoint in `handoff.md` at the end of every session.

## Layout
```
prototype/index.html      reference UI (immutable)
frontend/                 Vite shell + sync adapter (TypeScript)
backend/                  Go API (chi, pgx, Postgres)
specs/001-field-notes-platform/{spec,plan,tasks}.md
.specify/                 spec-kit memory/templates/scripts
.claude/skills/speckit-* spec-kit slash skills (/speckit-specify, -plan, -tasks, -implement)
handoff.md                session checkpoint log
```

## Run locally
```
# db: Postgres.app on 5432 (or: docker compose up -d db)
createdb urbanism            # once
# backend
cd backend && cp .env.example .env   # set DATABASE_URL
go run ./cmd/server          # :8080, migrates + seeds on boot
# frontend
cd frontend && cp .env.example .env && npm install && npm run dev   # :5173
```

## Verify before push
`cd backend && go build ./... && go vet ./... && go test ./...`
`cd frontend && npm run build`

## API (v1)
`GET /healthz` · `GET /state` · `GET|POST /sessions` · `GET /sessions/{code}` ·
`PUT /sessions/{id}/vote` · `GET /sessions/{id}/tally` · `GET|POST /entries` ·
`PUT|DELETE /entries/{id}/like` · `GET|POST /replies` · `GET|POST /materials` · `GET /journal`
Header `X-Client-Id` (8–64 `[A-Za-z0-9_-]`) scopes votes/likes.

## Deploy
- Railway: service from `backend/` (Dockerfile). Env: `DATABASE_URL`, `APP_ENV=production`,
  `ALLOWED_ORIGINS=https://aiabrooklynurbanism.com,https://www.aiabrooklynurbanism.com`.
- Vercel: project `urbanism`, root dir `frontend/`, Vite, domain aiabrooklynurbanism.com (www canonical). Env: `VITE_API_URL=https://urbanism-production.up.railway.app`. Deploy with `vercel --prod` from repo root.
- Railway: project `vigilant-amazement`, service `urbanism`, public https://urbanism-production.up.railway.app. Variable change = auto redeploy.

## Session history
See `handoff.md`. Session 1 (2026-11-11): project bootstrapped end-to-end.
