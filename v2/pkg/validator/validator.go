package validator

import (
	"bytes"
	"embed"
	"fmt"

	"github.com/porter-dev/switchboard/v2/pkg/types"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

//go:embed schemas/*.schema.json
var schemas embed.FS

func ValidatePorterYAML(parsed *types.ParsedPorterYAML) error {
	switch parsed.PorterYAML.Version.GetValue() {
	case "v2":
		return validateV2(parsed)
	}

	return fmt.Errorf("invalid porter.yaml version: %s", parsed.PorterYAML.Version.GetValue())
}

func validateV2(parsed *types.ParsedPorterYAML) error {
	// start by first validating against the schema
	porterYAMLSchemaBytes, err := schemas.ReadFile("schemas/porter_v2.schema.json")

	if err != nil {
		return fmt.Errorf("error reading porter_v2.schema.json: %w", err)
	}

	porterYAMLBuildSchemaBytes, err := schemas.ReadFile("schemas/porter_v2.build.schema.json")

	if err != nil {
		return fmt.Errorf("error reading porter_v2.build.schema.json: %w", err)
	}

	compiler := jsonschema.NewCompiler()
	compiler.AddResource("porter_v2.schema.json", bytes.NewReader(porterYAMLSchemaBytes))
	compiler.AddResource("porter_v2.build.schema.json", bytes.NewReader(porterYAMLBuildSchemaBytes))

	compiledSchema, err := compiler.Compile("porter_v2.schema.json")

	if err != nil {
		return fmt.Errorf("error compiling porter_v2.schema.json: %w", err)
	}

	err = compiledSchema.Validate(parsed.Raw)

	if err != nil {
		return fmt.Errorf("error validating porter.yaml against schema: %w", err)
	}

	err = validateV2Variables(parsed.PorterYAML.Variables)

	if err != nil {
		return err
	}

	err = validateV2Apps(parsed.PorterYAML.Apps)

	if err != nil {
		return err
	}

	err = validateV2Addons(parsed.PorterYAML.Addons)

	if err != nil {
		return err
	}

	return nil
}

func validateV2Variables(variables []*types.Variable) error {
	// let us do some basic checking against variables here
	//   - check that the name is unique

	vars := make(map[string]bool)

	for _, variable := range variables {
		if _, ok := vars[variable.Name.GetValue()]; ok {
			// duplicate variable name
			return fmt.Errorf("duplicate variable name: %s", variable.Name.GetValue())
		} else {
			vars[variable.Name.GetValue()] = true
		}
	}

	return nil
}

func validateV2Apps(apps []*types.AppResource) error {
	return nil
}

func validateV2Addons(addons []*types.AddonResource) error {
	return nil
}