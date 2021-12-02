package worker

import (
	"fmt"
	"os"

	"github.com/porter-dev/switchboard/internal/exec"
	"github.com/porter-dev/switchboard/internal/models"
	"github.com/porter-dev/switchboard/pkg/drivers"
	"github.com/porter-dev/switchboard/pkg/types"
	"github.com/rs/zerolog"
)

type Worker struct {
	driversTable map[string]drivers.DriverFunc
}

func NewWorker() *Worker {
	return &Worker{
		driversTable: make(map[string]drivers.DriverFunc),
	}
}

func (w *Worker) RegisterDriver(name string, driverFunc drivers.DriverFunc) error {
	if _, ok := w.driversTable[name]; ok {
		return fmt.Errorf("driver with name '%s' already exists", name)
	}

	w.driversTable[name] = driverFunc

	return nil
}

// Apply creates a ResourceGroup
func (w *Worker) Apply(group *types.ResourceGroup, opts *types.ApplyOpts) error {
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
		if driverFunc, ok := w.driversTable[resource.Driver]; ok {
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
