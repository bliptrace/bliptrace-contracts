# Versioning & pinning policy

> Derived from the platform specification §10 and §20.1. The envelope is a cross-repo contract; this
> policy is what keeps consumers from silently diverging.

## Two version numbers, kept in lockstep

| Where | Meaning |
|---|---|
| `schema_version` (inside the envelope) | The wire version a payload declares, e.g. `"1.0"`. Major bump = breaking. |
| `VERSION` (repo file) + git tag | The release of this repo. Semver. Consumers pin a git tag / Composer constraint. |

## Change classes

| Change | Action |
|---|---|
| **Breaking** (remove/rename a field, tighten a type, change semantics) | Bump **major** of `schema_version` and `VERSION`. Add a new `schema/vN/` directory; keep old majors. |
| **Additive backward-compatible** (new optional field) | Bump **minor**. Same `schema/vN/` dir. |
| **Editorial** (description, example) | Bump **patch**. |

## Multi-major coexistence

Old majors are **not deleted** — consumers may still pin `v1` while `v2` rolls out:

```
schema/
  v1/envelope.schema.json
  v2/envelope.schema.json   # added on a breaking change
```

The ingest API can accept multiple `schema_version` majors during a migration window and route to the
matching validator.

## How consumers pin

| Consumer | Pin |
|---|---|
| `bliptrace-api`, `bliptrace-mcp` (Go) | Go module tag: `require github.com/bliptrace/bliptrace-contracts vX.Y.Z` |
| `bliptrace-core` / `bliptrace-app` (Laravel) | Composer constraint (path or private registry) |
| SDKs | Embed/declare the `schema_version` they target |

## Release checklist

1. Edit `schema/vN/envelope.schema.json`.
2. Update `examples/vN/envelope.example.json` so it still conforms (the test enforces this).
3. Update `VERSION`.
4. `go test ./...` green.
5. Tag the release; re-pin consumers deliberately (one PR per consumer).
