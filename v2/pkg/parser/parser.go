package parser

import (
	"fmt"

	"github.com/porter-dev/switchboard/v2/pkg/types"
	"gopkg.in/yaml.v3"
)

func ParseRawBytes(raw []byte) (*types.ParsedPorterYAML, error) {
	porterYAML := &types.PorterYAML{}

	err := yaml.Unmarshal(raw, porterYAML)

	if err != nil {
		return nil, fmt.Errorf("error unmarshalling porter.yaml: %w", err)
	}

	data := make(map[string]interface{})

	err = yaml.Unmarshal(raw, &data)

	if err != nil {
		return nil, fmt.Errorf("error unmarshalling porter.yaml raw data: %w", err)
	}

	return &types.ParsedPorterYAML{
		PorterYAML: porterYAML,
		Raw:        data,
	}, nil
}
