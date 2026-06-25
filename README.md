# g8scope-contracts

The **source of truth** for the g8scope **error envelope** — the language-neutral wire contract every
SDK, the ingest API, and the fingerprinter agree on. Part of the [g8scope](https://g8suite.com)
platform by Developers Hub Sdn Bhd.

> This repo is **not deployed**. It is a **versioned dependency** that every consumer pins, so the
> API, web/core and all SDKs never silently diverge. See the platform specification §10, §20.1.

## What's here

| Path | What |
|---|---|
| `schema/v1/envelope.schema.json` | The canonical envelope JSON Schema (draft 2020-12) — **the contract** |
| `examples/v1/envelope.example.json` | A conforming example envelope |
| `contracts.go` | Embeds the schema + `Validate` helpers for Go consumers |
| `VERSION` | Current contract version (semver) |
| `docs/versioning.md` | Versioning + pinning policy |

## Who consumes it (and how)

| Consumer | Pull mechanism | Derives |
|---|---|---|
| `g8scope-api` (Go) | Tagged Go module (`import`) or codegen from schema | Envelope structs + runtime validator |
| `g8scope-core` / `g8scope-app` (Laravel) | Composer (path / private registry) | Envelope DTOs |
| SDKs (PHP, JS, WP, …) | Pin the schema version | Payload they POST (the universal `curl`-able envelope) |

## Anti-spoof rule

`project` is **intentionally not in the envelope**. It is resolved server-side from the ingest key
and stamped, so a client cannot spoof another tenant. The schema rejects a `project` field
(`additionalProperties: false`).

## Use from Go

```go
import contracts "github.com/developers-hub-my/g8scope-contracts"

if err := contracts.ValidateEnvelopeV1Bytes(raw); err != nil {
    // reject the envelope
}
raw, _ := contracts.EnvelopeSchemaV1() // embedded schema bytes
```

## Test

```bash
make test    # go test ./...  — validates the example + rejects bad/spoofed envelopes
```

## Versioning

Breaking envelope changes bump the **major** of `schema_version` and `VERSION`, and add a new
`schema/vN/` directory. Consumers pin a version. See `docs/versioning.md`.

---

_Developers Hub Sdn Bhd · devhub.my · g8scope-contracts_
