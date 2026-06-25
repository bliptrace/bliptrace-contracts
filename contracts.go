// Package contracts is the source of truth for the g8scope error envelope.
//
// The canonical artifact is the JSON Schema under schema/. This package embeds it
// and offers a Validate helper so Go consumers (g8scope-api) can validate envelopes
// and code-generate structs from the same single source. Laravel consumers
// (g8scope-core / g8scope-app) pull the same schema via Composer. Pin a version in
// every consumer so they never silently diverge.
//
// See the platform specification §10 and docs/.
package contracts

import (
	"bytes"
	"embed"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

//go:embed schema/v1/envelope.schema.json
var schemaFS embed.FS

// EnvelopeSchemaV1Path is the embedded path of the v1 envelope schema.
const EnvelopeSchemaV1Path = "schema/v1/envelope.schema.json"

// EnvelopeSchemaV1 returns the raw v1 envelope JSON Schema bytes.
func EnvelopeSchemaV1() ([]byte, error) {
	return schemaFS.ReadFile(EnvelopeSchemaV1Path)
}

// compileV1 compiles the embedded v1 envelope schema.
func compileV1() (*jsonschema.Schema, error) {
	raw, err := EnvelopeSchemaV1()
	if err != nil {
		return nil, fmt.Errorf("read schema: %w", err)
	}
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("parse schema: %w", err)
	}
	c := jsonschema.NewCompiler()
	if err := c.AddResource(EnvelopeSchemaV1Path, doc); err != nil {
		return nil, fmt.Errorf("add schema: %w", err)
	}
	return c.Compile(EnvelopeSchemaV1Path)
}

// ValidateEnvelopeV1 validates an already-unmarshalled JSON value (map[string]any,
// as produced by jsonschema.UnmarshalJSON or encoding/json into any) against the v1
// envelope schema. Returns nil when the value conforms.
func ValidateEnvelopeV1(instance any) error {
	sch, err := compileV1()
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
