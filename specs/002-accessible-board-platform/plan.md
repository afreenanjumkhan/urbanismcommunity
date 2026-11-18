# Implementation plan

Revise the canonical prototype and regenerate its legacy assets using a checked extraction script.
Keep readable enhancement sources alongside the prototype. Keep API integration in TypeScript.
Add PostgreSQL accounts/invitations and service-layer password/JWT/role checks. HTTP controllers
handle transport and rate limits. Require configured secrets; fail closed when absent.
Persist only token hashes for invitations and salted password hashes; keep browser JWTs in memory.
Use existing sync for public records; authenticated session creation must await server confirmation.
Review mobile and desktop layouts, navigation, consent rejection, loading failures, auth denial,
and run frontend build/parity plus backend tests/build/vet. Legal review remains organizational.

Brand refinement: keep the SVG asset in public/media/brand and embed it in the
canonical prototype during regeneration. Override navigation decoration with a
single shaped tab per page; keep the existing mobile disclosure and focus behavior.
