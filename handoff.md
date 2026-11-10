# handoff.md — Session checkpoints

Append one block per session. Newest at bottom. Keep each block short:
what changed, what was verified, what is next, anything surprising.

---

## 2026-11-11 · Session 1 · Bootstrap

**Done**
- Repo init, remote `origin` = github.com/starman011/urbanism, spec-kit scaffold (Claude integration), rtk 0.45 installed via brew.
- Prototype `index.html` analyzed: React bundle, hash routing, state doc `aia-field-notes-v1` (sessions, entries, replies, liked, bookmarks, attending, votes, notebook, materials), media in IndexedDB `aia-field-notes-media`.
- Prototype moved to `prototype/index.html`; CSS/JS split verbatim into `frontend/public/legacy/`.
- Backend Go API: domain / service / repository(postgres) / controller(http) / middleware; migrations + seed embedded; Dockerfile + railway.toml.
- Frontend: Vite shell, `src/sync.ts` adapter (boot pull, setItem diff push, 45 s poll), `vercel.json`.
- Docs: constitution v1.0.0, spec/plan/tasks 001, CLAUDE.md, README. Target domain: aiabrooklynurbanism.com.

**Verified**
- `go build ./... && go vet ./...` OK; `tsc --noEmit` OK. Local run status at end of session block below.

**Next**
- T12 service unit tests; T13 Playwright parity test; T14 create Railway + Vercel projects, env vars, attach domain.
- Phase 2: media bytes server-side; research export endpoint.

**Gotchas**
- Bundle validator rejects `null` arrays; API always returns `[]`.
