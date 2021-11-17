package helm

import (
	"io/ioutil"

	"helm.sh/helm/v3/pkg/action"

	"github.com/porter-dev/switchboard/pkg/drivers/kubernetes"
	"github.com/rs/zerolog"
)

type GetConfigOpts struct {
	// inherits the same config options as Kubernetes
	*kubernetes.GetConfigOpts
}

type GetAgentOpts struct {
	// The default namespace of the Helm agent
	Namespace string

	// The storage backend of the Helm agent
	Storage string

	// A Kubernetes agent
	Agent *kubernetes.Agent

	// Zerolog logger
	Logger *zerolog.Logger
}

// GetAgent creates a new Helm Agent
func GetAgent(opts *GetAgentOpts) (*Agent, error) {
	actionConf := &action.Configuration{}

	silentLogger := zerolog.New(ioutil.Discard)

	if err := actionConf.Init(opts.Agent.RESTClientGetter, opts.Namespace, opts.Storage, silentLogger.Printf); err != nil {
		return nil, err
	}

	// use k8s agent to create Helm agent
	return &Agent{
		ActionConfig: actionConf,
		K8sAgent:     opts.Agent,
		Logger:       opts.Logger,
	}, nil
}
