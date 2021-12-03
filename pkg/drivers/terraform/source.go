package terraform

import (
	"fmt"

	"github.com/porter-dev/switchboard/utils/objutils"
)

type VarMethod string

const (
	VarMethodFile VarMethod = "varfile"
	VarMethodEnv  VarMethod = "varenv"
)

const (
	SourceKindLocal string = "local"
)

type Source struct {
	*SourceLocal

	Kind      string
	VarMethod VarMethod
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

	varMethod, err := objutils.GetNestedString(genericSource, "var_method")

	if err != nil || varMethod == "" || VarMethodEnv == VarMethod(varMethod) {
		res.VarMethod = VarMethodEnv
	} else if VarMethodFile == VarMethod(varMethod) {
		res.VarMethod = VarMethodFile
	} else {
		return nil, fmt.Errorf("unknown var method %s", varMethod)
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
