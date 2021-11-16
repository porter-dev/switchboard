package drivers

import (
	"github.com/porter-dev/switchboard/internal/models"
)

type SharedDriverOpts struct {
	BaseDir           string
	DriverLookupTable *map[string]Driver
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
