package exec

import (
	"fmt"
	"sync"

	"github.com/porter-dev/switchboard/v2/pkg/types"
)

// TODO: this exec func should probably accept channels or something
type ExecFunc func(resource *types.Resource) error

type ExecNode struct {
	isExecFinished bool
	isExecStarted  bool
	execError      error
	parents        []*ExecNode
	resource       *types.Resource
}

func (e *ExecNode) IsFinished() bool {
	return e.isExecFinished
}

func (e *ExecNode) SetFinished() {
	e.isExecFinished = true
}

func (e *ExecNode) IsStarted() bool {
	return e.isExecStarted
}

func (e *ExecNode) SetStarted() {
	e.isExecStarted = true
}

func (e *ExecNode) SetFinishedWithError(err error) {
	e.isExecFinished = true
	e.execError = err
}

func (e *ExecNode) ExecError() error {
	return e.execError
}

func (e *ExecNode) ShouldStart() bool {
	// if the exec has started or finished, return false
	if e.IsStarted() || e.IsFinished() {
		return false
	}

	// if all parents have finished execution, the exec process should start
	parentsFinished := true

	for _, parent := range e.parents {
		parentsFinished = parentsFinished && parent.IsFinished()
	}

	return parentsFinished
}

// GetExecNodes
func GetExecNodes(parsed *types.PorterYAML) ([]*ExecNode, error) {
	// create a map of resource names to exec nodes
	resourceMap := make(map[string]*ExecNode)

	var resources []*types.Resource
	resources = append(resources, parsed.Apps.GetValue()...)
	resources = append(resources, parsed.Addons.GetValue()...)

	for _, resource := range resources {
		resourceMap[resource.Name.GetValue()] = &ExecNode{
			resource: resource,
		}
	}

	// Now that resources are registered, iterate through the resources again
	// to find the dependencies. If a dependency does not exist, throw an error
	res := make([]*ExecNode, 0)

	for _, execNode := range resourceMap {
		for _, dependency := range execNode.resource.DependsOn {
			execNode.parents = append(execNode.parents, resourceMap[dependency.GetValue()])
		}

		res = append(res, execNode)
	}

	return res, nil
}

// Execute simply calls exec on nodes in parallel, in batches. This could be much more
// efficient.
func Execute(nodes []*ExecNode, execFunc ExecFunc) {
	for {
		var wg sync.WaitGroup

		// get the list of nodes which are ready to execute, and execute those nodes
		for _, node := range nodes {
			nodeP := node
			if nodeP.ShouldStart() {
				wg.Add(1)

				go func() {
					defer wg.Done()

					nodeP.SetStarted()

					for _, parentNode := range nodeP.parents {
						if parentNode.ExecError() != nil {
							nodeP.SetFinishedWithError(fmt.Errorf("dependency '%s' failed", parentNode.resource.Name))
							return
						}
					}

					err := execFunc(nodeP.resource)

					if err != nil {
						nodeP.SetFinishedWithError(err)
						return
					}

					nodeP.SetFinished()
				}()
			}
		}

		wg.Wait()

		if allFinished := areAllNodesFinished(nodes); allFinished {
			break
		}
	}
}

func areAllNodesFinished(nodes []*ExecNode) bool {
	areFinished := true

	for _, node := range nodes {
		areFinished = areFinished && node.IsFinished()
	}

	return areFinished
}
