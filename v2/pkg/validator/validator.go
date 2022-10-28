package validator

import (
	"bytes"
	"embed"
	"fmt"
	"strconv"
	"strings"

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
		validationErr, ok := err.(*jsonschema.ValidationError)

		if ok {
			errStr := ""

			for _, c := range validationErr.Causes {
				absInstanceName := getAbsoluteInstanceName(c)
				absMessage := getAbsoluteMessage(c)

				line, col, err := getLineColForSchemaInstanceName(absInstanceName, parsed)

				if err != nil {
					return fmt.Errorf("error validating porter.yaml: %w", err)
				}

				errStr += fmt.Sprintf("\n- porter.yaml:%d:%d : %s", line, col, absMessage)
			}

			return fmt.Errorf("error validating porter.yaml against schema: %s", errStr)
		}

		return fmt.Errorf("error validating porter.yaml against schema")
	}

	err = validateV2Variables(parsed.PorterYAML.Variables.GetValue())

	if err != nil {
		return err
	}

	err = validateV2Apps(parsed.PorterYAML.Apps.GetValue())

	if err != nil {
		return err
	}

	err = validateV2Addons(parsed.PorterYAML.Addons.GetValue())

	if err != nil {
		return err
	}

	return nil
}

func getAbsoluteInstanceName(err *jsonschema.ValidationError) string {
	if len(err.Causes) == 0 {
		return err.InstanceLocation
	}

	return getAbsoluteInstanceName(err.Causes[0])
}

func getAbsoluteMessage(err *jsonschema.ValidationError) string {
	if len(err.Causes) == 0 {
		return err.Message
	}

	return getAbsoluteMessage(err.Causes[0])
}

func getLineColForSchemaInstanceName(instanceName string, parsed *types.ParsedPorterYAML) (int, int, error) {
	tree := strings.Split(strings.TrimPrefix(instanceName, "/"), "/")

	if len(tree) == 0 {
		return -1, -1, fmt.Errorf("invalid instance name: %s", instanceName)
	}

	switch tree[0] {
	case "variables":
		if len(tree) < 2 {
			return -1, -1, fmt.Errorf("invalid instance name: %s", instanceName)
		}

		buildIdx, err := strconv.Atoi(tree[1])

		if err != nil {
			return -1, -1, fmt.Errorf("invalid instance name: %s", instanceName)
		}

		if buildIdx < 0 || buildIdx >= len(parsed.PorterYAML.Variables.GetValue()) {
			return -1, -1, fmt.Errorf("invalid instance name: %s", instanceName)
		}

		if len(tree) > 2 {
			switch tree[2] {
			case "name":
				return parsed.PorterYAML.Variables.GetValue()[buildIdx].Name.GetLine(),
					parsed.PorterYAML.Variables.GetValue()[buildIdx].Name.GetColumn(), nil
			case "value":
				return parsed.PorterYAML.Variables.GetValue()[buildIdx].Value.GetLine(),
					parsed.PorterYAML.Variables.GetValue()[buildIdx].Value.GetColumn(), nil
			case "random":
				return parsed.PorterYAML.Variables.GetValue()[buildIdx].Random.GetLine(),
					parsed.PorterYAML.Variables.GetValue()[buildIdx].Random.GetColumn(), nil
			case "once":
				return parsed.PorterYAML.Variables.GetValue()[buildIdx].Once.GetLine(),
					parsed.PorterYAML.Variables.GetValue()[buildIdx].Once.GetColumn(), nil
			}
		}

	case "env_groups":

	case "builds":
		if len(tree) < 2 {
			return -1, -1, fmt.Errorf("invalid instance name: %s", instanceName)
		}

		buildIdx, err := strconv.Atoi(tree[1])

		if err != nil {
			return -1, -1, fmt.Errorf("invalid instance name: %s", instanceName)
		}

		if len(tree) > 2 {
			switch tree[2] {
			case "method":
				return parsed.PorterYAML.Builds.GetValue()[buildIdx].Method.GetLine(),
					parsed.PorterYAML.Builds.GetValue()[buildIdx].Method.GetColumn(), nil

			}
		} else {
			return parsed.PorterYAML.Builds.GetLine(), parsed.PorterYAML.Builds.GetColumn(), nil
		}
	case "apps":

	case "addons":
	}

	return -1, -1, fmt.Errorf("invalid instance name: %s", instanceName)
}

func validateV2Variables(variables []*types.Variable) error {
	// let us do some basic checking against variables here
	//   - check that the name is unique

	vars := make(map[string]bool)

	for _, variable := range variables {
		if _, ok := vars[variable.Name.GetValue()]; ok {
			// duplicate variable name
			return NewValidatorError(variable.Name.GetYAMLNodeMetadata(),
				fmt.Errorf("duplicate variable name: %s", variable.Name.GetValue()))
		} else {
			vars[variable.Name.GetValue()] = true
		}
	}

	return nil
}

func validateV2Apps(apps []*types.YAMLNode[*types.Resource]) error {
	return nil
}

func validateV2Addons(addons []*types.YAMLNode[*types.Resource]) error {
	return nil
}
