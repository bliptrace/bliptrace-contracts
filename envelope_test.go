package contracts

import (
	"os"
	"testing"
)

func TestValidateEnvelopeVersionRouting(t *testing.T) {
	example, err := os.ReadFile("examples/v1/envelope.example.json")
	if err != nil {
		t.Fatalf("read example: %v", err)
	}

	cases := []struct {
		name    string
		raw     []byte
		wantErr bool
	}{
		{"v1 example routes to v1 schema", example, false},
		{"unsupported major rejected", []byte(`{
			"schema_version":"9.0","event_id":"9f1c0a7e-1f2b-4d3e-9a8b-2c4d5e6f7a8b",
			"occurred_at":"2026-06-25T09:14:32.521Z","runtime":{"name":"php"},
			"level":"error","exception":{"type":"E","frames":[]}}`), true},
		{"missing schema_version rejected", []byte(`{
			"event_id":"9f1c0a7e-1f2b-4d3e-9a8b-2c4d5e6f7a8b",
			"occurred_at":"2026-06-25T09:14:32.521Z","runtime":{"name":"php"},
			"level":"error","exception":{"type":"E","frames":[]}}`), true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateEnvelope(tc.raw)
			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected valid, got: %v", err)
			}
		})
	}
}

func TestUnmarshalEnvelope(t *testing.T) {
	raw, err := os.ReadFile("examples/v1/envelope.example.json")
	if err != nil {
		t.Fatalf("read example: %v", err)
	}
	env, err := UnmarshalEnvelope(raw)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	checks := []struct {
		name string
		got  any
		want any
	}{
		{"schema_version", env.SchemaVersion, "1.1"},
		{"runtime name", env.Runtime.Name, "php"},
		{"level", env.Level, "error"},
		{"exception type", env.Exception.Type, "RuntimeException"},
		{"frame in_app", env.Exception.Frames[0].InApp, true},
		{"previous type", env.Exception.Previous[0].Type, "PDOException"},
		{"breadcrumb type", env.Breadcrumbs[0].Type, "http"},
		{"tag tenant", env.Tags["tenant"], "umpsa"},
		{"trace id", env.Trace.TraceID, "4bf92f3577b34da6a3ce929d0e0e4736"},
	}
	for _, c := range checks {
		t.Run(c.name, func(t *testing.T) {
			if c.got != c.want {
				t.Fatalf("%s = %v, want %v", c.name, c.got, c.want)
			}
		})
	}
}

func TestSupportedMajors(t *testing.T) {
	if len(SupportedMajors) == 0 {
		t.Fatal("SupportedMajors empty")
	}
	for _, m := range SupportedMajors {
		if _, err := EnvelopeSchema(m); err != nil {
			t.Fatalf("supported major v%d has no schema: %v", m, err)
		}
	}
}
