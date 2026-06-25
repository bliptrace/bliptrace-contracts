# OTLP ↔ envelope mapping (P3)

> Reference mapping between the g8scope error envelope (§10) and the OpenTelemetry data model, for
> the P3 **OTLP-compatible ingest endpoint** (features §3). This is a design reference, not a frozen
> contract — the concrete wiring lands with the api's OTLP endpoint and may refine these choices.

## Model choice

An error maps most naturally to an OpenTelemetry **LogRecord** that carries the exception
semantic-convention attributes (`exception.*`). When trace IDs are present, the LogRecord is
correlated to a span; an equivalent representation is a span **exception event**. g8scope ingests the
LogRecord form.

## Field mapping (envelope → OTLP)

| Envelope | OTLP target | Notes |
|---|---|---|
| `occurred_at` | `LogRecord.TimeUnixNano` | RFC3339 → epoch nanos |
| `level` | `LogRecord.SeverityNumber` + `SeverityText` | see severity table |
| `exception.type` | attribute `exception.type` | |
| `exception.message` | attribute `exception.message` | |
| `exception.frames[]` | attribute `exception.stacktrace` | rendered per runtime; `in_app` preserved in `g8scope.frames` |
| `event_id` | attribute `g8scope.event_id` (+ `LogRecord` dedup key) | idempotency key |
| `fingerprint_hint` | attribute `g8scope.fingerprint_hint` | optional grouping override |
| `runtime.name` / `runtime.version` | `process.runtime.name` / `process.runtime.version` (Resource) | drives normaliser |
| `sdk.name` / `sdk.version` | `telemetry.sdk.name` / `telemetry.sdk.version` (Resource) | |
| `environment` | `deployment.environment` (Resource) | |
| `release` | `service.version` (Resource) | |
| `context.request.method` / `.url` | `http.request.method` / `url.full` | |
| `context.user.id` | `enduser.id` | PII — scrubbed |
| `context.custom.*` | attributes (prefixed) | PII — scrubbed |
| `tags.*` | attributes | e.g. `tenant` |
| `trace.trace_id` / `span_id` | `LogRecord.TraceId` / `SpanId` | W3C hex |

> `project` has **no** OTLP source field — it is resolved server-side from the ingest key (anti-spoof),
> exactly as on the native path.

## Severity mapping

| Envelope `level` | OTLP `SeverityText` | `SeverityNumber` range |
|---|---|---|
| `debug` | DEBUG | 5–8 |
| `info` | INFO | 9–12 |
| `warning` | WARN | 13–16 |
| `error` | ERROR | 17–20 |
| `fatal` | FATAL | 21–24 |

## Reverse (OTLP → envelope), for the ingest endpoint

Incoming OTLP LogRecords/Spans carrying `exception.*` are translated to an envelope: Resource
attributes → `runtime`/`sdk`/`environment`/`release`; `exception.*` → `exception`; `TraceId`/`SpanId`
→ `trace`; remaining attributes → `tags`/`context`. The envelope then flows through the **same**
validate → scrub → fingerprint → sample pipeline as the native path — OTLP is an alternate front door,
not a separate pipeline.

## Open questions (resolve when building the api OTLP endpoint)

- Stacktrace fidelity: parse `exception.stacktrace` string back into structured frames, or require
  the `g8scope.frames` attribute for full grouping quality?
- Metrics/traces beyond errors are out of scope here — this mapping covers error capture only.
