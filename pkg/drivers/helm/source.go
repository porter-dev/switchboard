package helm

import "github.com/porter-dev/switchboard/internal/objutils"

type Source struct {
	ChartRepoURL string
	ChartName    string
	ChartVersion string
}

func GetSource(genericSource map[string]interface{}) (*Source, error) {
	res := &Source{}
	var err error

	res.ChartName, err = objutils.GetNestedString(genericSource, "chart_name")

	if err != nil {
		return nil, err
	}

	res.ChartRepoURL, err = objutils.GetNestedString(genericSource, "chart_repository")

	if err != nil {
		return nil, err
	}

	res.ChartVersion, err = objutils.GetNestedString(genericSource, "chart_version")

	if err != nil {
		return nil, err
	}

	return res, nil
}
