package worker

import (
	"fmt"
	"os"

	"github.com/porter-dev/switchboard/internal/exec"
	"github.com/porter-dev/switchboard/internal/query"
	"github.com/porter-dev/switchboard/pkg/drivers"
	"github.com/porter-dev/switchboard/pkg/models"
	"github.com/porter-dev/switchboard/pkg/types"
	"github.com/rs/zerolog"
)

type hookWithName struct {
	WorkerHook
	name string
}
type Worker struct {
	driversTable  map[string]drivers.DriverFunc
	hooks         []hookWithName
	defaultDriver string
}

func NewWorker() *Worker {
	return &Worker{
		driversTable:  make(map[string]drivers.DriverFunc),
		hooks:         make([]hookWithName, 0),
		defaultDriver: "",
	}
}

func (w *Worker) RegisterDriver(name string, driverFunc drivers.DriverFunc) error {
	if _, ok := w.driversTable[name]; ok {
		return fmt.Errorf("driver with name '%s' already exists", name)
	}

	w.driversTable[name] = driverFunc

	return nil
}

func (w *Worker) SetDefaultDriver(name string) error {
	if _, ok := w.driversTable[name]; !ok {
		return fmt.Errorf("attempting to set default driver with name '%s' that does not exist", name)
	}

	w.defaultDriver = name

	return nil
}

type WorkerHook interface {
	PreApply() error
	DataQueries() map[string]interface{}
	PostApply(populatedData map[string]interface{}) error
	OnError(err error)
}

func (w *Worker) RegisterHook(name string, hook WorkerHook) error {
	w.hooks = append(w.hooks, hookWithName{
		WorkerHook: hook,
		name:       name,
	})

	return nil
}

// Apply creates a ResourceGroup
func (w *Worker) Apply(group *types.ResourceGroup, opts *types.ApplyOpts) error {
	// run any pre-apply hooks
	for _, hook := range w.hooks {
		err := hook.WorkerHook.PreApply()

		if err != nil {
			return fmt.Errorf("error running hook '%s': %v", hook.name, err)
		}
	}

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
		if len(w.driversTable) == 0 {
			return fmt.Errorf("no drivers registered")
		} else if resource.Driver == "" {
			driver, err = w.driversTable[w.defaultDriver](modelResource, sharedDriverOpts)

			// TODO: append errors, don't exit here
			if err != nil {
				return err
			}
		} else if driverFunc, ok := w.driversTable[resource.Driver]; ok {
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
		w.runErrorHooks(err)
		return err
	}

	err = exec.Execute(nodes, execFunc)

	if err != nil {
		w.runErrorHooks(err)
		return err
	}

	// TODO: place in separate method, case on no hooks registered
	// get all output data if there are post-apply hooks
	allOutputData := make(map[string]interface{})

	for _, resource := range group.Resources {
		resourceOutput, err := lookupTable[resource.Name].Output()

		if err != nil {
			w.runErrorHooks(err)
			return err
		}

		allOutputData[resource.Name] = resourceOutput
	}

	// run any post-apply hooks
	for _, hook := range w.hooks {
		// get the data to query
		dataQueries := hook.WorkerHook.DataQueries()
		dataRes, err := query.PopulateQueries(dataQueries, allOutputData)

		if err != nil {
			return fmt.Errorf("error running hook '%s': %v", hook.name, err)
		}

		err = hook.WorkerHook.PostApply(dataRes)

		if err != nil {
			return fmt.Errorf("error running hook '%s': %v", hook.name, err)
		}
	}

	return nil
}

func (w *Worker) runErrorHooks(err error) {
	for _, hook := range w.hooks {
		hook.WorkerHook.OnError(err)
	}
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
