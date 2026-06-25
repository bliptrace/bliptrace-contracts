package contracts

import (
	"bytes"
	"encoding/json"
)

// Envelope is the typed Go representation of the v1 error envelope. It is hand-kept
// in lockstep with schema/v1/envelope.schema.json (the source of truth). Consumers
// that prefer raw maps can ignore these types and use ValidateEnvelope directly.
//
// NOTE: there is intentionally no Project field — the project is resolved server-side
// from the ingest key and stamped, never trusted from the payload (anti-spoof).
type Envelope struct {
	SchemaVersion   string            `json:"schema_version"`
	EventID         string            `json:"event_id"`
	OccurredAt      string            `json:"occurred_at"`
	Runtime         Runtime           `json:"runtime"`
	SDK             *SDK              `json:"sdk,omitempty"`
	Environment     string            `json:"environment,omitempty"`
	Release         *string           `json:"release,omitempty"`
	Level           string            `json:"level"`
	Exception       Exception         `json:"exception"`
	FingerprintHint *string           `json:"fingerprint_hint,omitempty"`
	Context         map[string]any    `json:"context,omitempty"`
	Breadcrumbs     []Breadcrumb      `json:"breadcrumbs,omitempty"`
	Tags            map[string]string `json:"tags,omitempty"`
	Trace           *Trace            `json:"trace,omitempty"`
}

// Runtime drives the per-runtime fingerprint normaliser.
type Runtime struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

// SDK identifies the client that produced the envelope.
type SDK struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

// Exception is the captured error and its stack frames.
type Exception struct {
	Type     string              `json:"type"`
	Message  string              `json:"message,omitempty"`
	Frames   []Frame             `json:"frames"`
	Previous []PreviousException `json:"previous,omitempty"` // chained causes, added in 1.1
}

// PreviousException is one link in the exception cause chain (added in 1.1).
type PreviousException struct {
	Type    string  `json:"type,omitempty"`
	Message string  `json:"message,omitempty"`
	Frames  []Frame `json:"frames,omitempty"`
}

// Frame is a single stack frame.
type Frame struct {
	File     string        `json:"file,omitempty"`
	Line     int           `json:"line,omitempty"`
	Function string        `json:"function,omitempty"`
	InApp    bool          `json:"in_app,omitempty"`
	Context  *FrameContext `json:"context,omitempty"`
}

// FrameContext is the surrounding source for a frame.
type FrameContext struct {
	Pre  []string `json:"pre,omitempty"`
	Line string   `json:"line,omitempty"`
	Post []string `json:"post,omitempty"`
}

// Breadcrumb is one event in the trail leading to the error.
type Breadcrumb struct {
	Timestamp string         `json:"timestamp"`
	Type      string         `json:"type"`
	Category  string         `json:"category,omitempty"`
	Level     string         `json:"level,omitempty"`
	Message   string         `json:"message,omitempty"`
	Data      map[string]any `json:"data,omitempty"`
}

// Trace carries optional W3C IDs for log<->trace correlation.
type Trace struct {
	TraceID string `json:"trace_id,omitempty"`
	SpanID  string `json:"span_id,omitempty"`
}

// UnmarshalEnvelope validates raw JSON against the matching schema major, then
// unmarshals it into a typed Envelope. Use this when you want both safety and types.
func UnmarshalEnvelope(raw []byte) (*Envelope, error) {
	if err := ValidateEnvelope(raw); err != nil {
		return nil, err
	}
	var e Envelope
	dec := json.NewDecoder(bytes.NewReader(raw))
	if err := dec.Decode(&e); err != nil {
		return nil, err
	}
	return &e, nil
}
