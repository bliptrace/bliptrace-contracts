// Package contracts is the source of truth for the g8scope error envelope.
//
// The canonical artifact is the JSON Schema under schema/. This package embeds it
// and offers Validate helpers + typed Go structs so Go consumers (g8scope-api) can
// validate envelopes and unmarshal them into types from the same single source.
// Laravel consumers (g8scope-core / g8scope-app) pull the same schema via Composer.
// Pin a version in every consumer so they never silently diverge.
//
// Multi-major coexistence: schema/vN/ directories live side by side. ValidateEnvelope
// reads schema_version, routes to the matching major's schema, and rejects unknown
// majors — so a v2 can roll out while consumers still pin v1. See docs/versioning.md
// and the platform specification §10.
package contracts

import (
	"bytes"
	"embed"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

//go:embed schema
var schemaFS embed.FS

// EnvelopeSchemaV1Path is the embedded path of the v1 envelope schema.
const EnvelopeSchemaV1Path = "schema/v1/envelope.schema.json"

// SupportedMajors lists the envelope schema majors this build can validate.
var SupportedMajors = []int{1}

var (
	compileMu sync.Mutex
	compiled  = map[int]*jsonschema.Schema{}
)

// schemaPath returns the embedded schema path for a major version.
func schemaPath(major int) string {
	return fmt.Sprintf("schema/v%d/envelope.schema.json", major)
}

// EnvelopeSchema returns the raw envelope JSON Schema bytes for a major version.
func EnvelopeSchema(major int) ([]byte, error) {
	return schemaFS.ReadFile(schemaPath(major))
}

// EnvelopeSchemaV1 returns the raw v1 envelope JSON Schema bytes.
func EnvelopeSchemaV1() ([]byte, error) { return EnvelopeSchema(1) }

// compile compiles (and caches) the embedded schema for a major version.
func compile(major int) (*jsonschema.Schema, error) {
	compileMu.Lock()
	defer compileMu.Unlock()
	if sch, ok := compiled[major]; ok {
		return sch, nil
	}
	path := schemaPath(major)
	raw, err := schemaFS.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unsupported schema major v%d: %w", major, err)
	}
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("parse schema v%d: %w", major, err)
	}
	c := jsonschema.NewCompiler()
	if err := c.AddResource(path, doc); err != nil {
		return nil, fmt.Errorf("add schema v%d: %w", major, err)
	}
	sch, err := c.Compile(path)
	if err != nil {
		return nil, err
	}
	compiled[major] = sch
	return sch, nil
}

// majorFromVersion parses the major from a "MAJOR.MINOR" schema_version string.
func majorFromVersion(v string) (int, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0, fmt.Errorf("missing schema_version")
	}
	head := v
	if i := strings.IndexByte(v, '.'); i >= 0 {
		head = v[:i]
	}
	major, err := strconv.Atoi(head)
	if err != nil {
		return 0, fmt.Errorf("invalid schema_version %q: %w", v, err)
	}
	return major, nil
}

// ValidateEnvelopeV1 validates an already-unmarshalled JSON value against the v1
// envelope schema. Returns nil when the value conforms.
func ValidateEnvelopeV1(instance any) error {
	sch, err := compile(1)
	if err != nil {
		return err
	}
	return sch.Validate(instance)
}

// ValidateEnvelopeV1Bytes validates raw JSON bytes against the v1 envelope schema.
func ValidateEnvelopeV1Bytes(raw []byte) error {
	inst, err := jsonschema.UnmarshalJSON(bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("parse instance: %w", err)
	}
	return ValidateEnvelopeV1(inst)
}

// ValidateEnvelope validates raw JSON bytes by routing on the envelope's own
// schema_version to the matching major schema. Unknown/unsupported majors are
// rejected. This is the entry point the ingest API should use.
func ValidateEnvelope(raw []byte) error {
	inst, err := jsonschema.UnmarshalJSON(bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("parse instance: %w", err)
	}
	obj, ok := inst.(map[string]any)
	if !ok {
		return fmt.Errorf("envelope must be a JSON object")
	}
	ver, _ := obj["schema_version"].(string)
	major, err := majorFromVersion(ver)
	if err != nil {
		return err
	}
	sch, err := compile(major)
	if err != nil {
		return fmt.Errorf("schema_version %q: %w", ver, err)
	}
	return sch.Validate(inst)
}
