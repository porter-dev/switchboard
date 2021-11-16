package kubernetes

import (
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/discovery"
	diskcached "k8s.io/client-go/discovery/cached/disk"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	// this line will register plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

type SourceType string

const (
	SourceTypeLocal SourceType = "local"
)

type GetConfigLocalOpts struct {
	Kubeconfig []byte
}

type GetConfigOpts struct {
	*GetConfigLocalOpts

	SourceType SourceType
}

func GetAgentFromHost(kubeconfigPath, context, defaultNamespace string) (*Agent, error) {
	cmdConf, err := GetClientCmdFromHost(kubeconfigPath, context, defaultNamespace)

	if err != nil {
		return nil, err
	}

	getter := &LocalRESTClientGetter{defaultNamespace, cmdConf}

	restConf, err := getter.ToRESTConfig()

	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(restConf)

	if err != nil {
		return nil, err
	}

	client, err := dynamic.NewForConfig(restConf)

	if err != nil {
		return nil, err
	}

	return &Agent{getter, clientset, client}, nil
}

type LocalRESTClientGetter struct {
	namespace string
	cmdConf   clientcmd.ClientConfig
}

func (l *LocalRESTClientGetter) ToRESTConfig() (*rest.Config, error) {
	restConf, err := l.cmdConf.ClientConfig()

	if err != nil {
		return nil, err
	}

	rest.SetKubernetesDefaults(restConf)
	return restConf, nil
}

// ToDiscoveryClient returns a CachedDiscoveryInterface using a computed RESTConfig
// It's required to implement the interface genericclioptions.RESTClientGetter
func (l *LocalRESTClientGetter) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	// From: k8s.io/cli-runtime/pkg/genericclioptions/config_flags.go > func (*configFlags) ToDiscoveryClient()
	restConf, err := l.ToRESTConfig()

	if err != nil {
		return nil, err
	}

	restConf.Burst = 100
	defaultHTTPCacheDir := filepath.Join(homedir.HomeDir(), ".kube", "http-cache")

	// takes the parentDir and the host and comes up with a "usually non-colliding" name for the discoveryCacheDir
	parentDir := filepath.Join(homedir.HomeDir(), ".kube", "cache", "discovery")
	// strip the optional scheme from host if its there:
	schemelessHost := strings.Replace(strings.Replace(restConf.Host, "https://", "", 1), "http://", "", 1)
	// now do a simple collapse of non-AZ09 characters.  Collisions are possible but unlikely.  Even if we do collide the problem is short lived
	safeHost := regexp.MustCompile(`[^(\w/\.)]`).ReplaceAllString(schemelessHost, "_")
	discoveryCacheDir := filepath.Join(parentDir, safeHost)

	return diskcached.NewCachedDiscoveryClientForConfig(restConf, discoveryCacheDir, defaultHTTPCacheDir, time.Duration(10*time.Minute))
}

// ToRESTMapper returns a mapper
func (l *LocalRESTClientGetter) ToRESTMapper() (meta.RESTMapper, error) {
	// From: k8s.io/cli-runtime/pkg/genericclioptions/config_flags.go > func (*configFlags) ToRESTMapper()
	discoveryClient, err := l.ToDiscoveryClient()
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(discoveryClient)
	expander := restmapper.NewShortcutExpander(mapper, discoveryClient)
	return expander, nil
}

// ToRawKubeConfigLoader creates a clientcmd.ClientConfig from the raw kubeconfig found in
// the OutOfClusterConfig. It does not implement loading rules or overrides.
func (l *LocalRESTClientGetter) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	return l.cmdConf
}
