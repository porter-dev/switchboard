package terraform

import (
	"fmt"

	"github.com/porter-dev/switchboard/internal/objutils"
)

const (
	SourceKindLocal string = "local"
)

type Source struct {
	*SourceLocal

	Kind string
}

type SourceLocal struct {
	Path string
}

func GetSource(genericSource map[string]interface{}) (*Source, error) {
	res := &Source{}
	var err error

	res.Kind, err = objutils.GetNestedString(genericSource, "kind")

	if err != nil {
		return nil, err
	}

	if res.Kind == "" {
		return nil, fmt.Errorf("source parameter \"kind\" must be set")
	}

	switch res.Kind {
	case SourceKindLocal:
		res.SourceLocal = &SourceLocal{}
		res.SourceLocal.Path, err = objutils.GetNestedString(genericSource, "path")

		if err != nil {
			return nil, err
		}
	}

	return res, nil
}
