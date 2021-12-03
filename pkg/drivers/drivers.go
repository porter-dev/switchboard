package drivers

import (
	"github.com/porter-dev/switchboard/internal/query"
	"github.com/porter-dev/switchboard/pkg/models"
	"github.com/rs/zerolog"
)

type SharedDriverOpts struct {
	BaseDir           string
	DriverLookupTable *map[string]Driver
	Logger            *zerolog.Logger
}

type QueryFunc func(data map[string]interface{}, query string) (interface{}, error)

type Driver interface {
	// ShouldApply returns true if the resource should be applied, false otherwise.
	// This enables the driver to pass pre-flight checks or detect if the configuration
	// has changed.
	ShouldApply(resource *models.Resource) bool

	// Apply writes the resource to the target.
	Apply(resource *models.Resource) (*models.Resource, error)

	// Output returns output data from the resource.
	Output() (map[string]interface{}, error)
}

type DriverFunc func(*models.Resource, *SharedDriverOpts) (Driver, error)

type ConstructConfigOpts struct {
	RawConf      map[string]interface{}
	LookupTable  map[string]Driver
	Dependencies []string
}

func ConstructConfig(opts *ConstructConfigOpts) (map[string]interface{}, error) {
	dataMap := make(map[string]interface{})

	for _, dependency := range opts.Dependencies {
		depOutput, err := opts.LookupTable[dependency].Output()

		if err != nil {
			return nil, err
		}

		dataMap[dependency] = depOutput
	}

	return query.PopulateQueries(opts.RawConf, dataMap)
}
