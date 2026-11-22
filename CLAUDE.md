# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this repo is

`bliptrace-contracts` is the **source of truth for the error envelope** — the single language-neutral
wire contract for g8scope (Developers Hub Sdn Bhd). It is **not a service**; it is a **versioned
dependency** pinned by every other repo so the API, web/core and all SDKs agree on one shape.

The canonical artifact is the **JSON Schema** at `schema/v1/envelope.schema.json`. Everything else
(Go structs, Laravel DTOs, SDK payloads) **derives** from it — never the other way round.

**Guardrail:** the envelope design is locked (spec §10). Do not add/rename fields casually — an
envelope change is a cross-repo contract change (see Versioning).

## Commands

| Task | Command |
|---|---|
| Test (validate example + reject bad envelopes) | `make test` / `go test ./...` |
| Single test | `go test -run TestValidateEnvelopeV1` |
| Tidy modules | `make tidy` |

## Layout

```
schema/v1/envelope.schema.json     the contract (draft 2020-12) — edit here
examples/v1/envelope.example.json  a conforming example (kept in sync with the schema)
contracts.go                       go:embed schema + Validate helpers for Go consumers
contracts_test.go                  table-driven validation tests
VERSION                            semver of the contract
docs/versioning.md                 versioning + pinning policy
```

## Load-bearing rules

- **`project` is NOT in the envelope.** It is resolved server-side from the ingest key and stamped
  (anti-spoof). The schema sets `additionalProperties: false` at the top level, so a `project` field
  is rejected. Keep it that way.
- **`event_id`** is the client UUID and the **idempotency key** for retry / spool replay.
- **`runtime.name`** drives the fingerprint normaliser registry downstream — keep it required.
- The example must always conform — the test enforces this. Update the example whenever the schema
  changes.

## Versioning (the whole reason this is a separate repo)

- Breaking change → bump the **major** of `schema_version` (e.g. `1.0` → `2.0`) and `VERSION`, and add
  a new `schema/vN/` directory alongside the old one (don't delete old majors — consumers pin them).
- Additive, backward-compatible change → bump minor.
- Every consumer **pins** a version (Go module tag / Composer constraint). After a bump, re-pin each
  consumer deliberately. Never let consumers float.

## Who derives from this

`bliptrace-api` (Go structs + validator) · `bliptrace-core` / `bliptrace-app` (Laravel DTOs) · all SDKs
(PHP, JS, WordPress, …). See `docs/versioning.md` and the platform spec §10, §18, §20.1.
