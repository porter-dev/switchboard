package kubernetes

import "fmt"

const (
	SourceKindNone  string = "none"
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

	sourceKind, sourceKindExists := genericSource["kind"]

	if sourceKindExists {
		if sourceKindStr, ok := sourceKind.(string); ok {
			res.Kind = sourceKindStr
		}
	}

	if res.Kind == "" {
		res.Kind = SourceKindNone
	}

	switch res.Kind {
	case SourceKindLocal:
		sourcePath, sourcePathExists := genericSource["path"]

		if !sourcePathExists {
			return nil, fmt.Errorf("source parameter \"path\" must be set when using \"local\" kind")
		}

		if sourcePathStr, ok := sourcePath.(string); ok {
			res.SourceLocal = &SourceLocal{
				Path: sourcePathStr,
			}
		} else if !ok {
			return nil, fmt.Errorf("source parameter \"path\" is not of type \"string\"")
		}
	}

	return res, nil
}
