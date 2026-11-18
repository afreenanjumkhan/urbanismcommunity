# 003 · Home page and series navigation

Owner request (2026-11-17): add a dynamic, interactive home page; move About to
the start of the navigation; number the folder tabs as a correct series; give the
folder tabs and the journal topic tabs a varied palette of the logo's red.

## Requirements
- FR-1 New route `#/home`, the default route and the logo's target. It reads the
  shared document (`aia-field-notes-v1`) and shows only real data.
- FR-2 Home includes an interactive block map: blocks respond to the pointer, a
  block opens a real field note, and a keyboard-accessible button cycles notes.
- FR-3 Home shows record counts (sessions, field notes, neighborhoods, replies),
  the next session (or the most recent one), and links to every folder.
- FR-4 Navigation order: Home, About, Journal, Sessions, Archives, Records,
  Contact. Each tab shows its position (01–07); the old archive count is removed.
- FR-5 Folder tabs, journal topic tabs and home folder cards use tints of the
  logo red (#B12F30); the active item uses the logo red. Text keeps WCAG AA contrast.
- FR-6 Motion (tab rise, hover growth, block lift, idle wave, count-up) is
  disabled under `prefers-reduced-motion`.

## Out of scope
Server changes; new stored fields (constitution V).
