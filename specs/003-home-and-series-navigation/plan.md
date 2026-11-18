# 003 · Plan

Prototype-first (constitution I). No hand edits to `frontend/public/legacy/`.

- `frontend/scripts/sync-prototype.mjs` applies anchored, idempotent swaps to the
  canonical bundle: nav order, per-tab series number, default route `/home`, a
  `/home` route that mounts an empty `#fn-home` container, page title, logo link.
  A swap fails loudly if its anchor disappears from a new prototype export.
- `prototype/enhancements.js` renders the home page into `#fn-home` with DOM APIs
  (`textContent` only, so user-written notes are never parsed as HTML). React
  never renders children there, so it leaves the injected content alone.
- `prototype/enhancements.css` adds the red series palette (`--red-1`…`--red-7`),
  folder/topic tab styling and motion, and home page layout.
- Verify with `npm run build` (tsc, parity, Vite) and headless-browser screenshots.
