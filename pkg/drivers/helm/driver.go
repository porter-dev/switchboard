package helm

import (
	"github.com/porter-dev/switchboard/internal/models"
	"github.com/porter-dev/switchboard/pkg/drivers"
)

type Driver struct {
	source      *Source
	target      *Target
	output      map[string]interface{}
	lookupTable *map[string]drivers.Driver
}

func NewHelmDriver(resource *models.Resource, opts *drivers.SharedDriverOpts) (*Driver, error) {
	driver := &Driver{
		lookupTable: opts.DriverLookupTable,
	}

	source, err := GetSource(resource.Source)

	if err != nil {
		return nil, err
	}

	driver.source = source

	target, err := GetTarget(resource.Target)

	if err != nil {
		return nil, err
	}

	driver.target = target

	return driver, nil
}

func (d *Driver) ShouldApply(resource *models.Resource) bool {
	return true
}

func (d *Driver) Apply(resource *models.Resource) (*models.Resource, error) {
	config, err := drivers.ConstructConfig(&drivers.ConstructConfigOpts{
		RawConf:      resource.Config,
		LookupTable:  *d.lookupTable,
		Dependencies: resource.Dependencies,
	})

	if err != nil {
		return nil, err
	}

	rel, err := d.target.agent.Apply(&ApplyOpts{
		Config: config,
		Target: d.target,
	})

	if err != nil {
		return nil, err
	}

	d.output = rel.Config

	return resource, nil
}

// Output returns the created Kubernetes configuration, including status section.
func (d *Driver) Output() (map[string]interface{}, error) {
	return d.output, nil
}
