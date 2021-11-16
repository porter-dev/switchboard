package kubernetes

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/porter-dev/switchboard/internal/models"
	"github.com/porter-dev/switchboard/internal/query"
	"github.com/porter-dev/switchboard/pkg/drivers"

	"sigs.k8s.io/yaml"
)

type Driver struct {
	agent       *Agent
	source      *Source
	target      *Target
	base        map[string]interface{}
	output      map[string]interface{}
	lookupTable *map[string]drivers.Driver
}

func NewKubernetesDriver(resource *models.Resource, opts *drivers.SharedDriverOpts) (*Driver, error) {
	driver := &Driver{
		lookupTable: opts.DriverLookupTable,
	}

	source, err := GetSource(resource.Source)

	if err != nil {
		return nil, err
	}

	driver.source = source

	err = driver.initSource(source, opts)

	if err != nil {
		return nil, err
	}

	target, err := GetTarget(resource.Target)

	if err != nil {
		return nil, err
	}

	driver.target = target

	err = driver.initTarget(target)

	if err != nil {
		return nil, err
	}

	return driver, nil
}

func (d *Driver) initSource(source *Source, opts *drivers.SharedDriverOpts) error {
	// read the file and set the base variable
	switch source.Kind {
	case SourceKindLocal:
		path := source.SourceLocal.Path
		base := make(map[string]interface{})

		// if the path is empty, just set the base to the empty map
		if path == "" {
			d.base = base
			return nil
		}

		// check if the filepath is absolute or relative
		if !filepath.IsAbs(source.SourceLocal.Path) {
			path = filepath.Join(opts.BaseDir, path)
		}

		// check if the file exists
		if info, err := os.Stat(path); os.IsNotExist(err) || info.IsDir() {
			return fmt.Errorf("source file specified by \"path\" does not exist or is a directory")
		}

		fileBytes, err := ioutil.ReadFile(path)

		if err != nil {
			return fmt.Errorf("error reading source file specified by \"path\": %v", err)
		}

		// parse the file bytes to yaml
		err = yaml.Unmarshal(fileBytes, &base)

		if err != nil {
			return fmt.Errorf("error parsing source file specified by \"path\" as yaml: %v", err)
		}

		d.base = base
	}

	return nil
}

func (d *Driver) initTarget(target *Target) error {
	// read the kubeconfig and construct an agent
	switch target.Kind {
	case TargetKindLocal:
		kubeBytes, err := GetKubeconfigFromHost(target.TargetLocal.KubeconfigPath, target.TargetLocal.KubeconfigContext)

		if err != nil {
			return fmt.Errorf("could not get kubeconfig: %v", err)
		}

		agent, err := GetAgentFromLocalKubeconfig(kubeBytes)

		if err != nil {
			return fmt.Errorf("could not get kube client: %v", err)
		}

		d.agent = agent
	}

	return nil
}

func (d *Driver) ShouldApply(resource *models.Resource) bool {
	return true
}

func (d *Driver) Apply(resource *models.Resource) (*models.Resource, error) {
	// get the config based on data population
	dataMap := make(map[string]interface{})
	lookupTable := *d.lookupTable

	for _, dependency := range resource.Dependencies {
		depOutput, err := lookupTable[dependency].Output()

		if err != nil {
			return nil, err
		}

		dataMap[dependency] = depOutput
	}

	config, err := query.PopulateQueries(resource.Config, dataMap)

	if err != nil {
		return nil, err
	}

	res, err := d.agent.Apply(&ApplyOpts{
		Config: config,
		Base:   d.base,
		Target: d.target,
	})

	if err != nil {
		return nil, err
	}

	d.output = res

	return resource, nil
}

// Output returns the created Kubernetes configuration, including status section.
func (d *Driver) Output() (map[string]interface{}, error) {
	return d.output, nil
}
