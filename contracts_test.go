package contracts

import (
	"os"
	"testing"
)

// Table-driven per platform convention. Proves the schema compiles, the canonical
// example conforms, and representative bad envelopes are rejected.
func TestValidateEnvelopeV1(t *testing.T) {
	example, err := os.ReadFile("examples/v1/envelope.example.json")
	if err != nil {
		t.Fatalf("read example: %v", err)
	}

	cases := []struct {
		name    string
		raw     []byte
		wantErr bool
	}{
		{"canonical example conforms", example, false},
		{"missing required event_id", []byte(`{
			"schema_version":"1.0","occurred_at":"2026-06-25T09:14:32.521Z",
			"runtime":{"name":"php"},"level":"error",
			"exception":{"type":"E","frames":[]}}`), true},
		{"invalid level enum", []byte(`{
			"schema_version":"1.0","event_id":"9f1c0a7e-1f2b-4d3e-9a8b-2c4d5e6f7a8b",
			"occurred_at":"2026-06-25T09:14:32.521Z","runtime":{"name":"php"},
			"level":"catastrophe","exception":{"type":"E","frames":[]}}`), true},
		{"server-stamped project rejected (anti-spoof)", []byte(`{
			"schema_version":"1.0","event_id":"9f1c0a7e-1f2b-4d3e-9a8b-2c4d5e6f7a8b",
			"occurred_at":"2026-06-25T09:14:32.521Z","runtime":{"name":"php"},
			"level":"error","exception":{"type":"E","frames":[]},
			"project":"someone-elses-tenant"}`), true},
		{"minimal valid envelope", []byte(`{
			"schema_version":"1.0","event_id":"9f1c0a7e-1f2b-4d3e-9a8b-2c4d5e6f7a8b",
			"occurred_at":"2026-06-25T09:14:32.521Z","runtime":{"name":"go"},
			"level":"fatal","exception":{"type":"panic","frames":[]}}`), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateEnvelopeV1Bytes(tc.raw)
			if tc.wantErr && err == nil {
				t.Fatalf("expected validation error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected valid, got error: %v", err)
			}
		})
	}
}

func TestEnvelopeSchemaV1Embedded(t *testing.T) {
	b, err := EnvelopeSchemaV1()
	if err != nil {
		t.Fatalf("read embedded schema: %v", err)
	}
	if len(b) == 0 {
		t.Fatal("embedded schema is empty")
	}
}
