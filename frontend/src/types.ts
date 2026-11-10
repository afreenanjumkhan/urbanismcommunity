// Mirrors the prototype's localStorage document (key aia-field-notes-v1, version 1).
export interface Session {
  id: string; code: string; title: string; date: string; endDate: string;
  place: string; neighborhood: string; desc: string; art: string;
  tags: string[]; notes: string[]; poll: string; options: string[]; sample?: boolean;
}
export interface Entry {
  id: string; sessionId: string; type: "text" | "voice" | "photo"; name: string; hood: string;
  text: string; tags: string[]; created: string; location?: string; art?: string;
}
export interface Reply { id: string; entryId: string; name: string; text: string; created: string }
export interface Material {
  id: string; sessionId: string; title: string; category: string; description: string; creditedTo: string;
  fileName: string; mimeType: string; mediaId: string; size: number; created: string;
}
export interface NotebookNote { id: string; text: string; created: string }

export interface ClientDoc {
  version: 1;
  sessions: Session[]; entries: Entry[]; replies: Reply[];
  liked: string[]; bookmarks: string[]; attending: string[];
  votes: Record<string, number>; notebook: NotebookNote[]; materials?: Material[];
}

export interface ServerState {
  version: 1; sessions: Session[]; entries: Entry[]; replies: Reply[]; materials: Material[];
  votes: Record<string, number>; liked: string[];
}
