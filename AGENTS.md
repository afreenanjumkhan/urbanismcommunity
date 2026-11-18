# Agent instructions

At the start of each working session, read `CLAUDE.md`, `handoff.md`, and
`.specify/memory/constitution.md`. Follow the constitution; it supersedes
ad-hoc practice. Consult the relevant feature spec, plan, and tasks before
implementation.

- Treat `prototype/index.html` as the visual and interaction contract. Do not
  directly edit `frontend/public/legacy/app.js` or `app.css`; revise the prototype
  first and regenerate with `frontend/scripts/sync-prototype.mjs`, then verify parity.
- Preserve backend layering: controller → service → repository, with domain as
  the shared contract. Keep SQL in repositories and business rules in services.
- Keep frontend and backend separately deployable and communicate through the
  versioned JSON API. Preserve local-only behavior when the API is unavailable.
- Preserve the append-only record, mark seed content `sample: true`, and follow
  the constitution's anonymous-by-default and minimal-data requirements.
- Follow the spec-driven workflow and the development and verification rules
  in the constitution, including required checks before pushing.
- Work directly on `main` for this single-user project; separate branches and
  pull requests are required only when explicitly requested by the owner.
- Whenever creating or changing something, append a concise checkpoint to
  `handoff.md` at a natural stopping point, and always do so before ending the
  working session. Record what changed, what was verified, what remains, and
  any surprises. Preserve prior entries.
