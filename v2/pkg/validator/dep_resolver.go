package validator

import (
	"fmt"

	"github.com/porter-dev/switchboard/v2/pkg/types"
)

type dependencyResolver struct {
	resources  []*types.Resource
	graph      map[string][]string
	resolved   map[string]bool
	unresolved map[string]bool
}

func newDependencyResolver(resources []*types.Resource) *dependencyResolver {
	return &dependencyResolver{
		resources:  resources,
		graph:      make(map[string][]string),
		resolved:   make(map[string]bool),
		unresolved: make(map[string]bool),
	}
}

func (r *dependencyResolver) Resolve() error {
	if len(r.resources) > 0 {
		// construct dependency graph
		for _, resource := range r.resources {
			// check for duplicate resource
			if _, ok := r.graph[resource.Name.GetValue()]; ok {
				return fmt.Errorf("duplicate app/addon name detected: '%s'", resource.Name.GetValue())
			}

			r.graph[resource.Name.GetValue()] = []string{}

			for _, dep := range resource.DependsOn {
				r.graph[resource.Name.GetValue()] = append(r.graph[resource.Name.GetValue()], dep.GetValue())
			}
		}

		for _, resource := range r.resources {
			err := r.depResolve(resource.Name.GetValue())

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *dependencyResolver) depResolve(name string) error {
	r.unresolved[name] = true

	for _, dep := range r.graph[name] {
		if _, ok := r.graph[dep]; !ok {
			return fmt.Errorf("for app/addon '%s': invalid dependency '%s'", name, dep)
		}

		if _, ok := r.resolved[dep]; !ok {
			if _, ok = r.unresolved[dep]; ok {
				return fmt.Errorf("circular depedency detected: '%s' -> '%s'", name, dep)
			}
			err := r.depResolve(dep)
			if err != nil {
				return err
			}
		}
	}

	r.resolved[name] = true
	delete(r.unresolved, name)

	return nil
}
