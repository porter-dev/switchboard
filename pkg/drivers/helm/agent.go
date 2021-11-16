package helm

import (
	"fmt"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"

	"github.com/porter-dev/switchboard/pkg/drivers/helm/loader"
	"github.com/porter-dev/switchboard/pkg/drivers/kubernetes"
	"github.com/rs/zerolog"
)

// Agent is a Helm agent for performing helm operations
type Agent struct {
	ActionConfig *action.Configuration
	K8sAgent     *kubernetes.Agent
	Logger       *zerolog.Logger

	release *release.Release
}

type ApplyOpts struct {
	Config map[string]interface{}
	Source *Source
	Target *Target
}

func (a *Agent) Apply(opts *ApplyOpts) (*release.Release, error) {
	err := a.loadRelease(opts.Source, opts.Target)

	if err != nil {
		// if error is not nil, we create the chart
		return a.installChart(opts.Source, opts.Target, opts.Config)
	}

	return a.upgradeRelease(opts.Source, opts.Target, opts.Config)
}

// GetRelease returns the info of a release.
func (a *Agent) loadRelease(
	source *Source,
	target *Target,
) error {
	// note: namespace is already known by the RESTClientGetter.
	cmd := action.NewGet(a.ActionConfig)
	cmd.Version = 0

	release, err := cmd.Run(target.Name)

	if err != nil {
		// TODO: case on whether the release exists or not
		return err
	}

	a.release = release

	if release.Chart != nil && release.Chart.Metadata != nil {
		loadDependencies(release.Chart)
	}

	return nil
}

func (a *Agent) upgradeRelease(
	source *Source,
	target *Target,
	values map[string]interface{},
) (*release.Release, error) {
	ch := a.release.Chart
	cmd := action.NewUpgrade(a.ActionConfig)
	cmd.Namespace = target.Namespace

	res, err := cmd.Run(target.Name, ch, values)

	if err != nil {
		return nil, fmt.Errorf("Upgrade failed: %v", err)
	}

	return res, nil
}

// installChart installs a new chart
func (a *Agent) installChart(
	source *Source,
	target *Target,
	values map[string]interface{},
) (*release.Release, error) {
	cmd := action.NewInstall(a.ActionConfig)
	cmd.ReleaseName = target.Name
	cmd.Namespace = target.Namespace
	cmd.Timeout = 300

	// load the chart
	chart, err := loader.LoadChartPublic(source.ChartRepoURL, source.ChartName, source.ChartVersion)

	if err != nil {
		return nil, err
	}

	return cmd.Run(chart, values)
}

func loadDependencies(chart *chart.Chart) {
	for _, dep := range chart.Metadata.Dependencies {
		depExists := false

		for _, currDep := range chart.Dependencies() {
			// we just case on name for now -- there might be edge cases we're missing
			// but this will cover 99% of cases
			if dep != nil && currDep != nil && dep.Name == currDep.Name() {
				depExists = true
				break
			}
		}

		if !depExists {
			depChart, err := loader.LoadChartPublic(dep.Repository, dep.Name, dep.Version)

			if err == nil {
				chart.AddDependency(depChart)
			}
		}
	}

	return
}
