# media/ — all site images, audio, and documents

Everything here deploys with the frontend and is served at
`https://www.aiabrooklynurbanism.com/media/<path>`.

```
media/
  sessions/<CODE>/   photos, audio, PDFs for one session  (e.g. sessions/ATL24/walk-01.jpg)
  journal/<post-id>/ images for a journal post            (e.g. journal/city-at-eye-level/cover.jpg)
  gallery/           general gallery / fieldwork photos   (e.g. gallery/2026-10-gowanus-edge.jpg)
```

Reference from anywhere as `/media/sessions/ATL24/walk-01.jpg`.

Rules
- Lowercase, hyphenated file names; keep originals under 2 MB (resize to max 2000 px).
- Formats: jpg, png, webp, svg, mp3, m4a, pdf.
- Commit + push → Vercel deploys. Served with long cache (`vercel.json`), so rename a file instead of replacing it in place.
- Uploads made in the browser by neighbors are NOT here (they live in that browser until phase 2 storage lands).
