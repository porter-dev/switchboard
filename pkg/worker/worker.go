package worker

import (
	"fmt"
	"time"

	"github.com/porter-dev/switchboard/internal/exec"
	"github.com/porter-dev/switchboard/internal/models"
	"github.com/porter-dev/switchboard/pkg/drivers"
	"github.com/porter-dev/switchboard/pkg/drivers/helm"
	"github.com/porter-dev/switchboard/pkg/drivers/kubernetes"
	"github.com/porter-dev/switchboard/pkg/types"
)

type ApplyOpts struct {
	BasePath string
}

// Apply creates a ResourceGroup
func Apply(group *types.ResourceGroup, opts *ApplyOpts) error {
	// create a map of resource names to drivers
	driverLookupTable := make(map[string]drivers.Driver)

	sharedDriverOpts := &drivers.SharedDriverOpts{
		BaseDir:           opts.BasePath,
		DriverLookupTable: &driverLookupTable,
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
		switch resource.Driver {
		case "kubernetes":
			driver, err = kubernetes.NewKubernetesDriver(modelResource, sharedDriverOpts)
		case "helm":
			driver, err = helm.NewHelmDriver(modelResource, sharedDriverOpts)
		}

		// TODO: append errors, don't exit here
		if err != nil {
			return err
		}

		driverLookupTable[resource.Name] = driver
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
		fmt.Println("RUNNING EXEC FOR", resource.Name)

		lookupTable := *opts.DriverLookupTable

		_, err := lookupTable[resource.Name].Apply(resource)

		// TODO: remove sleep statement
		time.Sleep(3 * time.Second)

		return err
	}
}
