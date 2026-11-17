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

## 2026-11-11 · Session 1b · First production deploy

**Done**
- Railway project `vigilant-amazement`, service `urbanism` (Dockerfile), Postgres plugin added, `DATABASE_URL=${{Postgres.DATABASE_URL}}`, `ALLOWED_ORIGINS` = apex + www. Public API: https://urbanism-production.up.railway.app
- Vercel project `urbanism` (root `frontend`, Vite). `VITE_API_URL` = Railway URL (prod + preview). Domain aiabrooklynurbanism.com (Cloudflare DNS → Vercel); apex 308 → www.
- `.vercel/` link lives at repo root (not `frontend/`); run `vercel --prod` from root.

**Verified**
- `/healthz` 200, `/api/v1/state` returns seeded sessions with CORS for www origin; www bundle has Railway URL baked in; legacy CSS/JS served.

**Gotchas**
- Changing a Railway variable auto-redeploys; `railway redeploy` right after says "cannot be redeployed".
- `vercel --prod` from `frontend/` fails ("Root Directory frontend does not exist"); run from repo root.

## 2026-11-17 · Session 2 · Contact page links

**Done**
- Prototype revision (owner request): filled the contact config in `prototype/index.html` and the identical bytes in `frontend/public/legacy/app.js`. Only the empty `of={...}` object changed: website = aiabrooklyn.org/committee/urbanism/, social (Instagram) = www-fallback.instagram.com/aiabkurbanism/, WhatsApp channel = chat.whatsapp.com/I2ZeCLIYhHaDBdu539fogK, email = urbanism@aiabrooklyn.org. `people` left empty.
- Rationale (constitution I): data-only change to the reference UI; prototype and bundle edited in lockstep so parity holds.

**Verified**
- Prototype `<script>` == bundle and `<style>` == app.css (parity logic re-run in PowerShell; Node not installed on this machine, so `npm run build` was not run).
- Not committed, not pushed.

**Next**
- Run `npm run build` where Node is available before pushing.
- Contact note "Links will appear here once the committee's contact details are confirmed." is now stale; needs a prototype copy change if it should go.

## 2026-11-17 · Agent instructions checkpoint

**Done**
- Read the constitution, `CLAUDE.md`, and this handoff; added `AGENTS.md` with project rules and the requested checkpoint practice.

**Verified**
- Reviewed the instructions against the constitution and standing brief. Documentation-only change; build checks were not run.

**Next**
- Existing follow-ups from Session 2 remain open.

## 2026-11-17 · Build and push checkpoint

**Done**
- Prepared contact-page-links for starman011/urbanism at the owner's request, using the configured afreenanjumkhan commit identity.
- Merged the destination branch history; both contact commits had identical trees.
- Added AGENTS.md and obtained local build tools in .build-tools (not committed).

**Verified**
- Frontend npm run build passed: TypeScript, exact prototype parity, and Vite production build.
- Backend go build ./..., go vet ./..., and go test ./... passed with Go 1.27.1.
- npm ci reported four high-severity dependency findings; versions were not changed.

**Next**
- Push this checkpoint and branch to the requested repository; address dependency findings separately.
- The stale contact note from Session 2 remains open.

## 2026-11-17 · Main workflow and constitution 1.1.0

**Done**
- Merged contact-page-links into the latest starman011/urbanism main, preserving existing main changes and commit authorship.
- Amended constitution to v1.1.0: direct work on main is the default; branches and PRs only on explicit owner request. Rationale: the owner is the sole active user and requested removal of unnecessary branch overhead.
- Aligned AGENTS.md and CLAUDE.md with this workflow. Previous branch push succeeded at 0284027.

**Verified**
- Merged code passed frontend build and prototype parity, plus backend build, vet, and tests.
- New commits use afreenanjumkhan's configured identity; historical authorship is preserved.

**Next**
- Push this verified main checkpoint to starman011/urbanism. Contributor display may take time to refresh.
- Existing dependency findings and stale contact note remain follow-ups.
