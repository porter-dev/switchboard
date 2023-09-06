package validator

import (
	"fmt"

	"github.com/porter-dev/switchboard/v2/pkg/types"
)

type ValidatorError struct {
	yamlNodeMetadata *types.YAMLNodeMetadata
	err              error
}

func NewValidatorError(metadata *types.YAMLNodeMetadata, err error) *ValidatorError {
	return &ValidatorError{
		yamlNodeMetadata: metadata,
		err:              err,
	}
}

func (e *ValidatorError) Error() string {
	return fmt.Sprintf("porter.yaml:%d:%d : %s", e.yamlNodeMetadata.Line, e.yamlNodeMetadata.Column, e.err.Error())
}
