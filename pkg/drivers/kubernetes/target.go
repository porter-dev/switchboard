package kubernetes

import "fmt"

const (
	TargetKindLocal string = "local"
)

type Target struct {
	*TargetLocal

	Kind      string
	Namespace string
}

type TargetLocal struct {
	KubeconfigPath    string
	KubeconfigContext string
}

func GetTarget(genericTarget map[string]interface{}) (*Target, error) {
	res := &Target{}

	// look for a target kind, which is required
	targetKind, targetKindExists := genericTarget["kind"]

	if targetKindExists {
		if targetKindStr, ok := targetKind.(string); ok {
			res.Kind = targetKindStr
		}
	}

	if res.Kind == "" {
		return nil, fmt.Errorf("target parameter \"kind\" must be set")
	}

	// look for a target namespace
	targetNS, targetNSExists := genericTarget["namespace"]

	if targetNSExists {
		if targetNSStr, ok := targetNS.(string); ok {
			res.Namespace = targetNSStr
		}
	}

	// if the target namespace does not exist, set it as "default"
	if res.Namespace == "" {
		res.Namespace = "default"
	}

	switch res.Kind {
	case TargetKindLocal:
		// if the target kind is local, the kubeconfig path and context can be optionally set
		res.TargetLocal = &TargetLocal{}

		kubePath, kubePathExists := genericTarget["kubeconfig_path"]
		kubeContext, kubeContextExists := genericTarget["kubeconfig_context"]

		if kubePathExists {
			if kubePathStr, ok := kubePath.(string); ok {
				res.TargetLocal.KubeconfigPath = kubePathStr
			}
		}

		if kubeContextExists {
			if kubeContextStr, ok := kubeContext.(string); ok {
				res.TargetLocal.KubeconfigContext = kubeContextStr
			}
		}
	}

	return res, nil
}
