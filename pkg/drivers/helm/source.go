package helm

import "github.com/porter-dev/switchboard/internal/objutils"

type SourceKind string

const (
	SourceKindRepository SourceKind = "repository"
	SourceKindLocal      SourceKind = "local"
)

type SourceLocal struct {
	Path string
}

type SourceRepository struct {
	ChartRepoURL string
	ChartName    string
	ChartVersion string
}
type Source struct {
	*SourceLocal
	*SourceRepository

	Kind SourceKind
}

func GetSource(genericSource map[string]interface{}) (*Source, error) {
	res := &Source{}
	var err error

	// read source kind
	kind, err := objutils.GetNestedString(genericSource, "kind")

	if err != nil {
		return nil, err
	}

	res.Kind = SourceKind(kind)

	switch SourceKind(kind) {
	case SourceKindLocal:
		res.SourceLocal = &SourceLocal{}

		res.SourceLocal.Path, err = objutils.GetNestedString(genericSource, "path")

		if err != nil {
			return nil, err
		}
	case SourceKindRepository:
		res.SourceRepository = &SourceRepository{}
		res.SourceRepository.ChartName, err = objutils.GetNestedString(genericSource, "chart_name")

		if err != nil {
			return nil, err
		}

		res.SourceRepository.ChartRepoURL, err = objutils.GetNestedString(genericSource, "chart_repository")

		if err != nil {
			return nil, err
		}

		res.SourceRepository.ChartVersion, err = objutils.GetNestedString(genericSource, "chart_version")

		if err != nil {
			return nil, err
		}
	}

	return res, nil
}
