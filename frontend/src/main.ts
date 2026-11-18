import { privacy } from "./privacy";
import { boot, markSessionSynced } from "./sync";
import { auth, authCall } from "./auth";
import type { Session } from "./types";
import { UI_VERSION } from './ui-version';

Object.assign(window,{fieldNotes:{privacy,auth,createSession:async(s:Session)=>{await authCall('/sessions','POST',s);markSessionSynced(s.id)}}});
const root=document.getElementById('root')!;
root.innerHTML='<div class="fn-boot" role="status" aria-live="polite">Opening the field notes…<div class="fn-skeleton" aria-hidden="true"></div><div class="fn-skeleton" aria-hidden="true"></div></div>';

// 1. Sync the client document with the API (or fall back to local-only).
// 2. Then load the untouched prototype bundle, which mounts into #root.
boot().finally(() => {
  const s = document.createElement("script");
  s.src = `/legacy/app.js?v=${UI_VERSION}`;
  s.onerror=()=>{root.textContent='The journal could not load. Check your connection and reload this page.'};
  document.body.appendChild(s);
});
