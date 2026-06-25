# Proposal: type-specific context for richer issue detail (envelope 1.1)

**Status:** Draft · **Target:** envelope `schema_version` 1.1 (additive, backward-compatible) · **Repo VERSION:** 1.0.0 → 1.1.0

## Why

The g8scope-app issue-detail page (g8scope-app#5) currently renders generic
frames + breadcrumbs. The legacy nadi.pro detail page showed per-runtime panels
(exception / HTTP / DB query) with code context. To unlock those panels
(g8scope-app#46–#50), the **envelope must carry the type-specific fields** — the
web can only render what the SDK sends.

## What's already in v1 (no change needed)

- **Code context** — `frame.context.{pre, line, post}` already exists. So
  **g8scope-app#50 (stack-frame code context) is NOT contract-blocked** — it is
  purely an SDK-populate + web-render task. Syntax highlighting is web-only.
- **Exception basics** — `exception.{type, message, frames}` → covers the core
  of g8scope-app#46.
- **Request basics** — `context.request.{method, url, headers}` → partial
  g8scope-app#47.
- **Breadcrumbs, tags, trace, user/custom context** — present.

## Additive fields (all OPTIONAL → backward-compatible, minor bump)

| Field | For | Unlocks |
|---|---|---|
| `exception.previous[]` (`{type, message, frames?}`) | chained/previous exceptions | g8scope-app#46 (full) |
| `context.request.{status_code, route, query, body, duration_ms}` | HTTP request reconstruction | g8scope-app#47 |
| `context.query` (`{connection, sql, bindings, duration_ms}`) | DB query that was running | g8scope-app#48 |

`context` is already `additionalProperties: true`, so these are documentation +
validation of an agreed shape rather than a hard gate. Naming them lets every
SDK and the web agree on the keys.

## Out of scope (deliberately deferred)

- Queue / command / log panels (g8scope-app#49, P2) — add `context.{job, command, log}`
  blocks in a later minor when those runtimes are prioritised.
- RUM / web-vitals / system metrics — these are g8scope **P3** (APM/RUM), not the
  error envelope.

## Roll-out (per CLAUDE.md cross-repo rule)

1. Land this schema as `1.1` (additive). Regenerate Go structs (`contracts.go`,
   `envelope.go`) and PHP DTOs (`src/`).
2. Tag `1.1.0`; re-pin consumers (`g8scope-api`, `g8scope-core`/SDK, `g8scope-app`).
3. SDK (g8scope-laravel) populates the new fields.
4. Web renders the per-type panels (g8scope-app#46–#48) + code context (#50).

This doc + the draft schema on this branch are the kickoff; the version bump and
consumer re-pin happen once the shape is reviewed.
