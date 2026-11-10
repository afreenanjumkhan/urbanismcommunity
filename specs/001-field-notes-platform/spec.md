# Feature Specification: Field Notes Platform (001)

**Status**: Implemented (phase 1) | **Created**: 2026-11-11 | **Branch**: main
**Target**: https://aiabrooklynurbanism.com

## Purpose

A public space where Brooklyn neighbors record feedback on urbanism sessions
(assemblies, fieldwork walks, labs) at the moment they happen, catch up on
sessions they missed, and add their own point of view. The accumulated record
becomes evidence for community-driven placemaking.

## User Scenarios

1. **Catch up on a session**: A neighbor who missed a walk opens the session by
   its code, reads the notes and other neighbors' field notes, and listens to or
   views what was shared.
2. **Add a field note**: A participant adds a text, voice, or photo note tagged
   by theme and neighborhood; it appears in the conversation for everyone.
3. **Respond**: Neighbors reply to a note and appreciate it; a poll on each
   session captures one priority per device.
4. **Create a session**: A committee member creates a new session with a code,
   date, place, notes, and poll; others join by the code.
5. **Share materials**: Panels, slides, handouts, and reports are attached to a
   session (metadata server-side; file bytes on-device in phase 1).
6. **Journal**: Editorial pieces by the committee link back to sessions.

## Functional Requirements

- FR-1 UI identical to `prototype/index.html` (constitution I).
- FR-2 Sessions, entries, replies, votes, appreciation, and material metadata
  persist server-side and are visible across devices.
- FR-3 Bookmarks, attending plans, and notebook stay on-device.
- FR-4 Offline/API-down: the app still loads with the last local document.
- FR-5 Input limits: title 120, text 4000, name 80, tags ≤ 8, options ≤ 6.
- FR-6 Session code: 3–12 alphanumerics, unique, uppercase.
- FR-7 One vote per device per session; changeable.
- FR-8 Seed content flagged `sample: true`.

## Out of Scope (phase 1)

Accounts/auth, moderation tools, media bytes on the server, analytics dashboard,
email notifications. See `plan.md` for phase 2.

## Key Entities

Session, Entry, Reply, Material, Vote, Reaction, JournalPost. Shapes mirror the
prototype's `aia-field-notes-v1` document (see `backend/internal/domain/models.go`).

## Success Criteria

- Two browsers see each other's new entries within 45 s without reload.
- `prototype/index.html` and `http://localhost:5173` render identically.
- Backend starts from an empty database and serves seeded sessions.
