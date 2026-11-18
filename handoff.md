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

## 2026-11-17 · Accessible board platform implementation checkpoint

**Done**
- Main push at 568532c succeeded. Started spec 002 on main.
- Revised prototype with owner palette, mobile menu, focus/motion support, image states, consent and board dialogs; regenerated legacy assets with parity checks.
- Added JWT/password services, invitation activation, super-admin board controls, account audit migration, and backend protection for session creation.
- Amended constitution to 1.2.0: owner explicitly requested accounts and a new visual prototype; clarified regeneration and anonymous public access with optional persistence.
- Google sign-in is high priority in backlog.md per owner clarification.

**Verified so far**
- Frontend build/parity and initial backend tests pass. Browser preview works outside the sandbox on 127.0.0.1:5174; mobile menu and create-session sign-in dialog inspected.

**Next**
- Complete security tests, accessibility fixes, README/SVG diagrams, and final validation. Configure private production auth secrets before live account activation.

## 2026-11-17 · Accessible board platform final checkpoint

**Done**
- Completed spec 002 UI, consent, JWT accounts and board authorization. Preserved journal/folder design with the requested red/ink/paper palette; increased reading text and fixed 320/390px overflow.
- Added README architecture/engineering/setup/privacy documentation, two SVG UML diagrams, browser screenshots, and explicit verification limits. Google sign-in remains high-priority backlog.
- Versioned legacy asset URLs and removed year-long immutable caching for fixed filenames so returning visitors receive revised assets.
- Local history was externally rewritten (reflog: "rewrite authorship", base d28a522); requested GitHub main remains 568532c. Their starting trees are identical. Preserve local history and publish the new snapshot as a normal child of GitHub main, without force-pushing.

**Verified**
- Frontend build, TypeScript and exact prototype parity pass; consent behavior test passes.
- Backend build, vet and tests pass, including new token expiry/tampering/revocation, role, bootstrap and anonymous-session-denial tests.
- Browser review at 1440, 390 and 320px; menu opening/Escape/focus, session sign-in gate, dialog focus return, consent reopening and necessary-only selection checked. No journal images missing alt or buttons missing accessible names in DOM inspection. Screenshots in docs/.

**Next / limits**
- Configure AUTH_JWT_SECRET and AUTH_BOOTSTRAP_TOKEN privately on Railway, run the migration in staging, activate akhannyc72@gmail.com, then remove the bootstrap token. No production secrets or credentials were created.
- No local PostgreSQL/Docker: real-database lifecycle verification remains a staging step. Full screen-reader/legal review and throttled-network media checks remain; no compliance certification is claimed.
- See backlog.md for Google sign-in, recovery, dependency audit and organizational review follow-ups. A GitHub push does not verify a production deployment.

## 2026-11-17 · Vector brand and folder-notch checkpoint

**Done**
- Recreated the supplied logo as SVG (white B/slashes, black AIA circle and Brooklyn lettering on Signal Red), stored in frontend/public/media/brand/aia-brooklyn.svg and embedded into the canonical prototype.
- Simplified header page navigation to only the raised label/notch silhouette; removed the separate miniature-folder body/index decorations. Retained labels, active state, focus rings, and mobile disclosure behavior.
- Regenerated legacy assets and versioned URLs. Updated spec 002 refinement notes.

**Verified**
- Frontend build and prototype parity pass. Browser review at 1440px and 320px confirms loaded SVG, correct tab appearance and no horizontal overflow; mobile menu Escape behavior retained.
- Header screenshot saved in docs/brand-header-preview.jpg.
- Required backend build, vet and tests also pass; SVG parses successfully as XML.

**Next**
- Publish to requested GitHub main while preserving the externally rewritten local history (previous source publication: faf83df). Existing auth configuration and staging checks remain as documented above.

## 2026-11-17 · AIA typography revision

**Done**
- Owner asked for fonts similar to the official AIA sites. aiabrooklyn.org sets headings/nav/labels in Knockout (H&FJ) and text in Architype (MCKL, custom for AIA); both are licensed, so open-license stand-ins are used: Oswald for the Knockout role, Figtree for Architype (closest match in a side-by-side render).
- `prototype/enhancements.css`: typography block — `--serif` now Oswald (all 54 heading/title rules), reading-text rules and former Lexend Deca accents on Figtree, serif-tuned negative tracking reset, nav/buttons/eyebrows/status on Oswald.
- Fonts self-hosted (no Google request): `prototype/fonts/*.woff2` (latin, ~70 KB) + OFL licenses; `sync-prototype.mjs` embeds them as data URIs like the logo.

**Verified**
- `npm run build` (tsc, prototype parity, Vite) and privacy test pass. Headless-Edge screenshots of journal, contact, about, sessions (desktop) and journal (390 px) look right.
- Node here is Adobe Photoshop's bundled v22 (copied to a scratch dir), not the constitution's Node 24.

**Next**
- Pre-existing, not caused by this change: at 390 px the journal's margin-note card and topic tabs overflow the right edge.
- Not committed or pushed.

## 2026-11-17 · Spec 003: home page and series navigation

**Done**
- New `#/home` route (default; logo target) via anchored, idempotent swaps in `sync-prototype.mjs`; the bundle mounts an empty `#fn-home` and passes live app state through a ref (`fnState`). `prototype/enhancements.js` renders it: interactive block map (pointer lift, idle wave, tap a block for a real field note, keyboard "Show another field note"), record counts with count-up, next/most recent session, numbered folder cards.
- Nav order Home, About, Journal, Sessions, Archives, Records, Contact; each tab shows its series number 01–07 (the old archive-count superscript, which read "04", is gone).
- Red series palette `--red-1`…`--red-7` (logo #B12F30 tints, ink contrast >= 4.5:1) on folder tabs (rise-in and hover growth), journal topic tabs (current = logo red) and home folders.
- Spec, plan, tasks in `specs/003-home-and-series-navigation/`.

**Verified**
- `npm run build` (tsc, parity, Vite), privacy test; sync run twice applies each swap once; bundle parses.
- Screenshots: home and journal at 1400 px, home at a true 390 px (iframe), reduced-motion home (static, counts 4/13/9/0).
- Correction: the 390 px overflow noted in the typography checkpoint was a headless-Edge artifact (its window cannot go below ~500 px, so shots were cropped). Measured at a true 390 px, home and journal do not overflow.

**Next**
- Not committed or pushed. Hover/pointer motion verified by code review, not by an interactive browser session.
