import { boot } from "./sync";

// 1. Sync the client document with the API (or fall back to local-only).
// 2. Then load the untouched prototype bundle, which mounts into #root.
boot().finally(() => {
  const s = document.createElement("script");
  s.src = "/legacy/app.js";
  document.body.appendChild(s);
});
