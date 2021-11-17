package helm

import (
	"github.com/porter-dev/switchboard/internal/objutils"
	"github.com/porter-dev/switchboard/pkg/drivers/kubernetes"
	"github.com/rs/zerolog"
)

type Target struct {
	*kubernetes.Target

	agent *Agent

	// Helm-specific fields
	Name string
}

func GetTarget(genericTarget map[string]interface{}, logger *zerolog.Logger) (*Target, error) {
	res := &Target{}

	kubeTarget, err := kubernetes.GetTarget(genericTarget)

	if err != nil {
		return nil, err
	}

	res.Target = kubeTarget

	// get the Helm agent from the kube Agent
	agent, err := GetAgent(&GetAgentOpts{
		Namespace: kubeTarget.Namespace,
		Storage:   "secret",
		Agent:     kubeTarget.Agent,
		Logger:    logger,
	})

	if err != nil {
		return nil, err
	}

	res.Name, err = objutils.GetNestedString(genericTarget, "name")

	if err != nil {
		return nil, err
	}

	res.agent = agent

	return res, nil
}
