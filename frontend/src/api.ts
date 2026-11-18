import type { Entry, Material, Reply, ServerState, Session } from "./types";
import { authorization } from './auth';

const BASE = (import.meta.env.VITE_API_URL as string | undefined)?.replace(/\/$/, "") ?? "";
const CLIENT_KEY = "aia-field-notes-client";

export function clientId(): string {
  let id = localStorage.getItem(CLIENT_KEY);
  if (!id) {
    id = crypto.randomUUID().replace(/-/g, "");
    localStorage.setItem(CLIENT_KEY, id);
  }
  return id;
}

async function call<T>(method: string, path: string, body?: unknown, timeoutMs = 6000): Promise<T> {
  const ctl = new AbortController();
  const t = setTimeout(() => ctl.abort(), timeoutMs);
  try {
    const res = await fetch(`${BASE}/api/v1${path}`, {
      method,
      headers: { "Content-Type": "application/json", "X-Client-Id": clientId(), ...authorization() },
      body: body === undefined ? undefined : JSON.stringify(body),
      signal: ctl.signal,
    });
    if (!res.ok) throw new Error(`${method} ${path} -> ${res.status}`);
    return res.status === 204 ? (undefined as T) : ((await res.json()) as T);
  } finally {
    clearTimeout(t);
  }
}

export const api = {
  enabled: BASE !== "",
  state: () => call<ServerState>("GET", "/state"),
  createSession: (s: Session) => call<Session>("POST", "/sessions", s),
  createEntry: (e: Entry) => call<Entry>("POST", "/entries", e),
  createReply: (r: Reply) => call<Reply>("POST", "/replies", r),
  createMaterial: (m: Material) => call<Material>("POST", "/materials", m),
  vote: (sessionId: string, option: number) => call<void>("PUT", `/sessions/${sessionId}/vote`, { option }),
  like: (entryId: string, on: boolean) => call<void>(on ? "PUT" : "DELETE", `/entries/${entryId}/like`),
};
