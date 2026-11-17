# Field Notes (AIA Brooklyn Urbanism) Constitution

The city is a collective work. This platform records what neighbors say about
their streets, in the moment they say it, so that the record can inform
community-driven placemaking and urban development.

## Core Principles

### I. The Prototype Is the UI (NON-NEGOTIABLE)
`prototype/index.html` is the visual and interaction contract. The shipped UI
MUST be pixel- and behavior-identical to it. The bundle (`frontend/public/legacy/app.js`)
and stylesheet (`frontend/public/legacy/app.css`) are served verbatim and MUST
NOT be edited. Integration happens around the bundle (boot adapter, storage
hooks), never inside it. Any visual change requires a new prototype first.

### II. Three-Tier Backend, Strict Layering
Go backend layers: `controller/http` (transport) → `service` (business rules)
→ `repository` (persistence) with `domain` as the shared contract and
`middleware` for cross-cutting concerns. Dependencies point inward only:
controller imports service and domain; service imports domain; repository
imports domain. No SQL outside `repository`; no HTTP types outside `controller`
and `middleware`. Business rules (validation, limits, defaults) live in
`service` and are unit-testable without a database.

### III. UI Is a Separate Deployable
Frontend (`frontend/`, Vercel) and backend (`backend/`, Railway) are separate
entities with separate builds, environments, and release cadences. They
communicate only over the versioned JSON API (`/api/v1`). The frontend MUST
degrade to local-only (the prototype's behavior) when the API is unreachable.

### IV. The Record Is Append-Only and Honest
Feedback is a historical record. Sessions, entries, and replies are appended,
never silently edited or deleted by the system. Seed content is flagged
`sample: true` and stays distinguishable from real contributions. Personal
slices (bookmarks, attending, notebook) remain on-device; shared slices
(sessions, entries, replies, votes, appreciation, material metadata) are
persisted server-side.

### V. Anonymous by Default, Minimal Data
No accounts in phase 1. An anonymous per-device client id scopes votes and
appreciation. The API stores no emails, IPs beyond transient logs, or tracking
identifiers. Every stored field maps to something the prototype shows.

### VI. Spec-Driven, Session-Logged
Every feature starts in `specs/<n>-<slug>/spec.md`, then `plan.md`, then
`tasks.md` (spec-kit flow). `CLAUDE.md` is the standing brief for any agent
session; `handoff.md` is the running checkpoint log and MUST be appended at the
end of every working session.

## Technology Constraints

- Frontend: React (prototype bundle) under Vite + TypeScript; Node 24.
- Backend: Go ≥ 1.27, chi router, pgx/v5, PostgreSQL (Railway).
- Deploy: Vercel (frontend root `frontend/`, domain aiabrooklynurbanism.com), Railway (Dockerfile in `backend/`).
- Tooling: spec-kit (`specify`), rtk for token-efficient agent command output.
- API: JSON, versioned under `/api/v1`, arrays never `null`, errors as `{"error": msg}`.

## Development Workflow

- Single-user workflow: work, commit, and push directly on `main` by default.
  Separate feature branches and pull requests are not required; create them
  only when the owner explicitly requests them. Retain spec references for features.
- `go build ./... && go vet ./... && go test ./...` and `npm run build` MUST pass before push.
- Schema changes are new numbered files in `backend/internal/repository/postgres/migrations/`.
- Never commit `.env`; `.env.example` is the contract.

## Governance

This constitution supersedes ad-hoc practice. Amendments require a short
rationale appended to `handoff.md` and a version bump below. Reviews check
layering (II), prototype parity (I), and the append-only rule (IV) first.

**Version**: 1.1.0 | **Ratified**: 2026-11-11 | **Last Amended**: 2026-11-17
