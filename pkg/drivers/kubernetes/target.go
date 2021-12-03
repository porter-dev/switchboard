package kubernetes

import (
	"fmt"

	"github.com/porter-dev/switchboard/utils/objutils"
)

const (
	TargetKindLocal string = "local"
)

type Target struct {
	*TargetLocal

	Kind      string
	Namespace string
	Agent     *Agent
}

type TargetLocal struct {
	KubeconfigPath    string
	KubeconfigContext string
}

func GetTarget(genericTarget map[string]interface{}) (*Target, error) {
	res := &Target{}
	var err error

	// look for a target kind, which is required
	res.Kind, err = objutils.GetNestedString(genericTarget, "kind")

	if err != nil {
		return nil, err
	}

	if res.Kind == "" {
		return nil, fmt.Errorf("target parameter \"kind\" must be set")
	}

	// look for a target namespace
	res.Namespace, err = objutils.GetNestedString(genericTarget, "namespace")

	// if the target namespace does not exist, set it as "default"
	if res.Namespace == "" {
		res.Namespace = "default"
	}

	switch res.Kind {
	case TargetKindLocal:
		// if the target kind is local, the kubeconfig path and context can be optionally set
		res.TargetLocal = &TargetLocal{}

		res.TargetLocal.KubeconfigPath, _ = objutils.GetNestedString(genericTarget, "kubeconfig_path")
		res.TargetLocal.KubeconfigContext, _ = objutils.GetNestedString(genericTarget, "kubeconfig_context")

		agent, err := GetAgentFromHost(
			res.TargetLocal.KubeconfigPath,
			res.TargetLocal.KubeconfigContext,
			res.Namespace,
		)

		if err != nil {
			return nil, fmt.Errorf("could not get kube client: %v", err)
		}

		res.Agent = agent
	}

	return res, nil
}
