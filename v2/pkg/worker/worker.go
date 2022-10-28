package worker

import (
	"fmt"

	"github.com/porter-dev/switchboard/v2/internal/exec"
	"github.com/porter-dev/switchboard/v2/pkg/types"
)

// type hookWithName struct {
// 	WorkerHook
// 	name string
// }

type Worker struct {
	driversTable map[string]types.Driver
	// hooks         []hookWithName
	defaultDriver string
}

func NewWorker() *Worker {
	return &Worker{
		driversTable: make(map[string]types.Driver),
		// hooks:         make([]hookWithName, 0),
		defaultDriver: "",
	}
}

func (w *Worker) RegisterDriver(name string, driver types.Driver) error {
	if _, ok := w.driversTable[name]; ok {
		return fmt.Errorf("driver with name '%s' already exists", name)
	}

	w.driversTable[name] = driver

	return nil
}

func (w *Worker) SetDefaultDriver(name string) error {
	if _, ok := w.driversTable[name]; !ok {
		return fmt.Errorf("attempting to set default driver with name '%s' that does not exist", name)
	}

	w.defaultDriver = name

	return nil
}

// type WorkerHook interface {
// 	PreApply() error
// 	DataQueries() map[string]interface{}
// 	PostApply(populatedData map[string]interface{}) error
// 	OnError(err error)
// 	OnConsolidatedErrors(allErrors map[string]error)
// }

// func (w *Worker) RegisterHook(name string, hook WorkerHook) error {
// 	w.hooks = append(w.hooks, hookWithName{
// 		WorkerHook: hook,
// 		name:       name,
// 	})

// 	return nil
// }

// Apply creates the apps and addons
func (w *Worker) Apply(parsed *types.PorterYAML) error {
	// // run any pre-apply hooks
	// for _, hook := range w.hooks {
	// 	err := hook.WorkerHook.PreApply()

	// 	if err != nil {
	// 		return fmt.Errorf("error running hook '%s': %v", hook.name, err)
	// 	}
	// }

	// // create a map of resource names to drivers
	// lookupTable := make(map[string]drivers.Driver)
	// stdOut := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})

	// sharedDriverOpts := &drivers.SharedDriverOpts{
	// 	BaseDir:           opts.BasePath,
	// 	DriverLookupTable: &lookupTable,
	// 	Logger:            &stdOut,
	// }

	execFunc := getExecFunc(w.driversTable)

	// var resources []*types.Resource

	// for _, resource := range Resources {
	// 	modelResource := &models.Resource{
	// 		Name:         resource.Name,
	// 		Driver:       resource.Driver,
	// 		Config:       resource.Config,
	// 		Source:       resource.Source,
	// 		Target:       resource.Target,
	// 		Dependencies: resource.DependsOn,
	// 	}

	// 	resources = append(resources, modelResource)

	// 	// var driver drivers.Driver
	// 	var err error

	// 	// switch on the driver type to construct the driver
	// 	if len(w.driversTable) == 0 {
	// 		return fmt.Errorf("no drivers registered")
	// 	} else if resource.Driver == "" {
	// 		driver, err = w.driversTable[w.defaultDriver](modelResource, sharedDriverOpts)

	// 		// TODO: append errors, don't exit here
	// 		if err != nil {
	// 			return err
	// 		}
	// 	} else if driverFunc, ok := w.driversTable[resource.Driver]; ok {
	// 		driver, err = driverFunc(modelResource, sharedDriverOpts)

	// 		// TODO: append errors, don't exit here
	// 		if err != nil {
	// 			return err
	// 		}
	// 	} else {
	// 		// TODO: append errors, don't exit here
	// 		err = fmt.Errorf("no driver found with name '%s'", resource.Driver)
	// 		return err
	// 	}

	// 	lookupTable[resource.Name] = driver
	// }

	// depResolver := exec.NewDependencyResolver(resources)
	// err := depResolver.Resolve()

	// if err != nil {
	// 	w.runErrorHooks(err)
	// 	return err
	// }

	nodes, err := exec.GetExecNodes(parsed)

	if err != nil {
		// w.runErrorHooks(err)
		return err
	}

	exec.Execute(nodes, execFunc)

	// allErrors := make(map[string]error)

	// for _, node := range nodes {
	// 	if node.ExecError() != nil {
	// 		allErrors[node.ResourceName()] = node.ExecError()
	// 	}
	// }

	// if len(allErrors) > 0 {
	// 	for _, hook := range w.hooks {
	// 		hook.OnConsolidatedErrors(allErrors)
	// 	}

	// 	return fmt.Errorf("errors were encountered with one or more resources")
	// }

	// // TODO: place in separate method, case on no hooks registered
	// // get all output data if there are post-apply hooks
	// allOutputData := make(map[string]interface{})

	// for _, resource := range group.Resources {
	// 	resourceOutput, err := lookupTable[resource.Name].Output()

	// 	if err != nil {
	// 		w.runErrorHooks(err)
	// 		return err
	// 	}

	// 	allOutputData[resource.Name] = resourceOutput
	// }

	// // run any post-apply hooks
	// for _, hook := range w.hooks {
	// 	// get the data to query
	// 	dataQueries := hook.WorkerHook.DataQueries()
	// 	dataRes, err := query.PopulateQueries(dataQueries, allOutputData)

	// 	if err != nil {
	// 		return fmt.Errorf("error running hook '%s': %v", hook.name, err)
	// 	}

	// 	err = hook.WorkerHook.PostApply(dataRes)

	// 	if err != nil {
	// 		return fmt.Errorf("error running hook '%s': %v", hook.name, err)
	// 	}
	// }

	return nil
}

// func (w *Worker) runErrorHooks(err error) {
// 	for _, hook := range w.hooks {
// 		hook.WorkerHook.OnError(err)
// 	}
// }

func getExecFunc(driverLookupTable map[string]types.Driver) exec.ExecFunc {
	return func(resource *types.Resource) error {
		// opts.Logger.Info().Msg(
		// 	fmt.Sprintf("running apply for resource %s", resource.Name),
		// )

		// lookupTable := *opts.DriverLookupTable

		err := driverLookupTable[resource.Name.GetValue()].Apply(resource)

		if err != nil {
			return err
		}

		// opts.Logger.Info().Msg(
		// 	fmt.Sprintf("successfully applied resource %s", resource.Name),
		// )

		return nil
	}
}
