package worker

import (
	"fmt"
	"os"

	"github.com/porter-dev/switchboard/internal/exec"
	"github.com/porter-dev/switchboard/internal/models"
	"github.com/porter-dev/switchboard/pkg/drivers"
	"github.com/porter-dev/switchboard/pkg/drivers/helm"
	"github.com/porter-dev/switchboard/pkg/drivers/kubernetes"
	"github.com/porter-dev/switchboard/pkg/drivers/terraform"
	"github.com/porter-dev/switchboard/pkg/types"
	"github.com/rs/zerolog"
)

var driversTable map[string]drivers.DriverFunc

type ApplyOpts struct {
	BasePath       string
	Logger         *zerolog.Logger
	ResourceLogger *zerolog.Logger
}

// Apply creates a ResourceGroup
func Apply(group *types.ResourceGroup, opts *ApplyOpts) error {
	// create a map of resource names to drivers
	lookupTable := make(map[string]drivers.Driver)
	stdOut := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})

	sharedDriverOpts := &drivers.SharedDriverOpts{
		BaseDir:           opts.BasePath,
		DriverLookupTable: &lookupTable,
		Logger:            &stdOut,
	}

	execFunc := getExecFunc(sharedDriverOpts)

	resources := make([]*models.Resource, 0)

	for _, resource := range group.Resources {
		modelResource := &models.Resource{
			Name:         resource.Name,
			Driver:       resource.Driver,
			Config:       resource.Config,
			Source:       resource.Source,
			Target:       resource.Target,
			Dependencies: resource.DependsOn,
		}

		resources = append(resources, modelResource)

		var driver drivers.Driver
		var err error

		// switch on the driver type to construct the driver
		if driverFunc, ok := driversTable[resource.Driver]; ok {
			driver, err = driverFunc(modelResource, sharedDriverOpts)

			// TODO: append errors, don't exit here
			if err != nil {
				return err
			}
		} else {
			// TODO: append errors, don't exit here
			err = fmt.Errorf("no driver found with name '%s'", resource.Driver)
			return err
		}

		lookupTable[resource.Name] = driver
	}

	nodes, err := exec.GetExecNodes(&models.ResourceGroup{
		APIVersion: group.Version,
		Resources:  resources,
	})

	if err != nil {
		return err
	}

	return exec.Execute(nodes, execFunc)
}

func getExecFunc(opts *drivers.SharedDriverOpts) exec.ExecFunc {
	return func(resource *models.Resource) error {
		opts.Logger.Info().Msg(
			fmt.Sprintf("running apply for resource %s", resource.Name),
		)

		lookupTable := *opts.DriverLookupTable

		_, err := lookupTable[resource.Name].Apply(resource)

		if err != nil {
			return err
		}

		opts.Logger.Info().Msg(
			fmt.Sprintf("successfully applied resource %s", resource.Name),
		)

		return nil
	}
}

func RegisterDriver(name string, driverFunc drivers.DriverFunc) error {
	if _, ok := driversTable[name]; ok {
		return fmt.Errorf("driver with name '%s' already exists", name)
	}

	driversTable[name] = driverFunc

	return nil
}

func init() {
	driversTable = make(map[string]drivers.DriverFunc)

	RegisterDriver("helm", helm.NewHelmDriver)
	RegisterDriver("kubernetes", kubernetes.NewKubernetesDriver)
	RegisterDriver("terraform", terraform.NewTerraformDriver)
}
