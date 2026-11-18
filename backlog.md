# Backlog

## High priority — Google sign-in
Add Google Identity Services with server-side signature, audience, issuer and expiry verification.
Link only verified identities; preserve server-side board admission and the super-admin role.
Configure OAuth client ID and allowed origins with the owner. Do not auto-promote arbitrary Google users.

## Follow-ups
- Staging PostgreSQL migration and end-to-end account lifecycle verification; configure production auth secrets and activate the owner.
- Password recovery and verified email delivery; document trusted super-admin recovery.
- Full screen-reader audit, physical-device coverage, and throttled-network media checks.
- Review the existing four npm high-severity audit findings and upgrade with regression checks.
- Organization review of privacy notice, retention, incident response, and applicable US/state laws.
