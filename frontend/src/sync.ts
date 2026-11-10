// Sync adapter: keeps the prototype's localStorage document and the API in step
// without touching the UI bundle. Shared slices (sessions, entries, replies,
// materials metadata, votes, liked) round-trip to the server; personal slices
// (bookmarks, attending, notebook) stay on this device.
import { api } from "./api";
import type { ClientDoc, ServerState } from "./types";

export const DOC_KEY = "aia-field-notes-v1";
const POLL_MS = 45_000;

let last: ClientDoc | null = null;
let queue: Promise<void> = Promise.resolve();

function readLocal(): ClientDoc | null {
  try {
    const raw = localStorage.getItem(DOC_KEY);
    return raw ? (JSON.parse(raw) as ClientDoc) : null;
  } catch {
    return null;
  }
}

function merge(local: ClientDoc | null, server: ServerState): ClientDoc {
  const byId = <T extends { id: string }>(a: T[], b: T[]) => {
    const seen = new Set(a.map((x) => x.id));
    return [...a, ...b.filter((x) => !seen.has(x.id))];
  };
  // Local-only materials (bytes live in IndexedDB on this device) are kept.
  const localMaterials = local?.materials ?? [];
  return {
    version: 1,
    sessions: server.sessions,
    entries: server.entries,
    replies: server.replies,
    materials: byId(localMaterials, server.materials),
    votes: { ...(local?.votes ?? {}), ...server.votes },
    liked: Array.from(new Set([...(local?.liked ?? []), ...server.liked])),
    bookmarks: local?.bookmarks ?? [],
    attending: local?.attending ?? [],
    notebook: local?.notebook ?? [],
  };
}

const origSetItem = Storage.prototype.setItem;

function enqueue(job: () => Promise<unknown>) {
  queue = queue.then(() => job().then(() => undefined, (e) => console.warn("[sync]", e)));
}

function diffAndPush(prev: ClientDoc | null, next: ClientDoc) {
  const ids = <T extends { id: string }>(a: T[] | undefined) => new Set((a ?? []).map((x) => x.id));
  const newOf = <T extends { id: string }>(a: T[] | undefined, b: T[] | undefined) => {
    const have = ids(a);
    return (b ?? []).filter((x) => !have.has(x.id));
  };
  for (const s of newOf(prev?.sessions, next.sessions)) enqueue(() => api.createSession(s));
  for (const e of newOf(prev?.entries, next.entries)) enqueue(() => api.createEntry(e));
  for (const r of newOf(prev?.replies, next.replies)) enqueue(() => api.createReply(r));
  for (const m of newOf(prev?.materials, next.materials)) enqueue(() => api.createMaterial(m));
  for (const [sid, opt] of Object.entries(next.votes ?? {})) {
    if (prev?.votes?.[sid] !== opt) enqueue(() => api.vote(sid, opt));
  }
  const before = new Set(prev?.liked ?? []);
  const after = new Set(next.liked ?? []);
  for (const id of after) if (!before.has(id)) enqueue(() => api.like(id, true));
  for (const id of before) if (!after.has(id)) enqueue(() => api.like(id, false));
}

function install() {
  Storage.prototype.setItem = function (key: string, value: string) {
    origSetItem.call(this, key, value);
    if (this !== localStorage || key !== DOC_KEY) return;
    try {
      const next = JSON.parse(value) as ClientDoc;
      diffAndPush(last, next);
      last = next;
    } catch (e) {
      console.warn("[sync] unreadable document", e);
    }
  };
}

async function pull(): Promise<boolean> {
  const server = await api.state();
  const merged = merge(readLocal(), server);
  last = merged;
  origSetItem.call(localStorage, DOC_KEY, JSON.stringify(merged));
  return true;
}

function notifyApp() {
  // The bundle reloads its state on a storage event for DOC_KEY.
  window.dispatchEvent(new StorageEvent("storage", { key: DOC_KEY, storageArea: localStorage }));
}

export async function boot(): Promise<void> {
  if (!api.enabled) {
    console.info("[sync] VITE_API_URL not set; running local-only like the prototype.");
    return;
  }
  install();
  try {
    await pull();
  } catch (e) {
    console.warn("[sync] API unreachable at boot; using local document.", e);
    last = readLocal();
  }
  setInterval(async () => {
    try {
      await queue;
      await pull();
      notifyApp();
    } catch {
      /* offline; try again next tick */
    }
  }, POLL_MS);
}
