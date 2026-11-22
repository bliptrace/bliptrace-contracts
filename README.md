# bliptrace-contracts

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
| `contracts.go` / `envelope.go` | Go: embedded schema, version-routing validator, typed `Envelope` structs |
| `src/Schema.php` | PHP/Composer: locate + decode the schema (`G8Scope\Contracts\Schema`) |
| `composer.json` | Composer package `g8scope/contracts` |
| `VERSION` | Current contract version (semver) |
| `docs/versioning.md` · `docs/otlp-mapping.md` | Versioning + pinning policy · OTLP ↔ envelope mapping (P3) |

## Who consumes it (and how)

| Consumer | Pull mechanism | Derives |
|---|---|---|
| `bliptrace-api` (Go) | Tagged Go module (`import`) or codegen from schema | Envelope structs + runtime validator |
| `bliptrace-core` / `bliptrace-app` (Laravel) | Composer (path / private registry) | Envelope DTOs |
| SDKs (PHP, JS, WP, …) | Pin the schema version | Payload they POST (the universal `curl`-able envelope) |

## Anti-spoof rule

`project` is **intentionally not in the envelope**. It is resolved server-side from the ingest key
and stamped, so a client cannot spoof another tenant. The schema rejects a `project` field
(`additionalProperties: false`).

## Use from Go (bliptrace-api / bliptrace-mcp)

```go
import contracts "github.com/bliptrace/bliptrace-contracts"

// Version-routing validate (reads schema_version, picks the matching major):
if err := contracts.ValidateEnvelope(raw); err != nil {
    // reject the envelope
}

// Validate + unmarshal into typed structs:
env, err := contracts.UnmarshalEnvelope(raw)
_ = env.Exception.Type

schemaBytes, _ := contracts.EnvelopeSchema(1) // embedded schema bytes for a major
```

## Use from Laravel (bliptrace-core / bliptrace-app)

Add the repo as a Composer VCS source, then require it (private repo — git auth applies):

```jsonc
// composer.json
"repositories": [
  { "type": "vcs", "url": "https://github.com/bliptrace/bliptrace-contracts" }
]
```

```bash
composer require g8scope/contracts:^1.0
```

```php
use G8Scope\Contracts\Schema;

$path  = Schema::path(1);              // absolute path to the v1 schema
$array = Schema::decode(1);            // schema as an array (derive DTOs from this)
$major = Schema::majorFromVersion('1.0'); // => 1
```

Envelope DTOs are **derived** from this schema in `bliptrace-core` — the envelope is never redefined.

## Test

```bash
make test    # go test ./...  — validates the example + rejects bad/spoofed envelopes
```

## Versioning

Breaking envelope changes bump the **major** of `schema_version` and `VERSION`, and add a new
`schema/vN/` directory. Consumers pin a version. See `docs/versioning.md`.

---

_Developers Hub Sdn Bhd · devhub.my · bliptrace-contracts_
