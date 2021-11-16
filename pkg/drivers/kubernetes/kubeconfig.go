package kubernetes

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"
)

// GetClientCmdFromHost returns a clientcmd from the host. Finds default kubeconfig and context if not set
func GetClientCmdFromHost(kubeconfigPath, context, defaultNamespace string) (clientcmd.ClientConfig, error) {
	kubeconfigPath, err := resolveKubeconfigPath(kubeconfigPath)

	if err != nil {
		return nil, err
	}

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.ExplicitPath = kubeconfigPath

	clientConf := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{
		Context: clientcmdapi.Context{
			Namespace: defaultNamespace,
		},
	})

	rawConf, err := clientConf.RawConfig()

	if err != nil {
		return nil, err
	}

	if context == "" {
		if rawConf.CurrentContext == "" {
			return nil, fmt.Errorf("could not detect current context: at least one context must be specified")
		}

		context = rawConf.CurrentContext
	}

	populateCertificateRefs(&rawConf)
	populateOIDCPluginCerts(&rawConf)

	conf, err := stripAndValidateClientContexts(&rawConf, context, []string{context}, defaultNamespace)

	if err != nil {
		return nil, err
	}

	return conf, err
}

// resolveKubeconfigPath finds the path to a kubeconfig, first searching for the
// passed string, then in the home directory, then as an env variable.
func resolveKubeconfigPath(kubeconfigPath string) (string, error) {
	envVarName := clientcmd.RecommendedConfigPathEnvVar

	if kubeconfigPath != "" {
		if _, err := os.Stat(kubeconfigPath); os.IsNotExist(err) {
			// the specified kubeconfig does not exist, throw error
			return "", fmt.Errorf("kubeconfig not found: %s does not exist", kubeconfigPath)
		}
	}

	if kubeconfigPath == "" {
		if os.Getenv(envVarName) == "" {
			if home := homedir.HomeDir(); home != "" {
				kubeconfigPath = filepath.Join(home, ".kube", "config")
			}
		} else {
			kubeconfigPath = os.Getenv(envVarName)
		}
	}

	return kubeconfigPath, nil
}

func stripAndValidateClientContexts(
	rawConf *clientcmdapi.Config,
	currentContext string,
	allowedContexts []string,
	defaultNamespace string,
) (clientcmd.ClientConfig, error) {
	// grab a copy to get the pointer and set clusters, authinfos, and contexts to empty
	copyConf := rawConf.DeepCopy()

	copyConf.Clusters = make(map[string]*api.Cluster)
	copyConf.AuthInfos = make(map[string]*api.AuthInfo)
	copyConf.Contexts = make(map[string]*api.Context)
	copyConf.CurrentContext = currentContext

	// put allowed clusters in a map
	aContextMap := createAllowedContextMap(allowedContexts)

	for contextName, context := range rawConf.Contexts {
		userName := context.AuthInfo
		clusterName := context.Cluster
		authInfo, userFound := rawConf.AuthInfos[userName]
		cluster, clusterFound := rawConf.Clusters[clusterName]

		// make sure the cluster is "allowed"
		_, isAllowed := aContextMap[contextName]

		if userFound && clusterFound && isAllowed {
			copyConf.Clusters[clusterName] = cluster
			copyConf.AuthInfos[userName] = authInfo
			copyConf.Contexts[contextName] = context
		}
	}

	// validate the copyConf and create a ClientConfig
	err := clientcmd.Validate(*copyConf)

	if err != nil {
		return nil, err
	}

	clientConf := clientcmd.NewDefaultClientConfig(*copyConf, &clientcmd.ConfigOverrides{
		Context: clientcmdapi.Context{
			Namespace: defaultNamespace,
		},
	})

	return clientConf, nil
}

func populateCertificateRefs(config *clientcmdapi.Config) {
	for _, cluster := range config.Clusters {
		refs := clientcmd.GetClusterFileReferences(cluster)
		for _, str := range refs {
			// only write certificate if the file reference is CA
			if *str != cluster.CertificateAuthority {
				break
			}

			fileBytes, err := ioutil.ReadFile(*str)

			if err != nil {
				continue
			}

			cluster.CertificateAuthorityData = fileBytes
			cluster.CertificateAuthority = ""
		}
	}

	for _, authInfo := range config.AuthInfos {
		refs := clientcmd.GetAuthInfoFileReferences(authInfo)
		for _, str := range refs {
			if *str == "" {
				continue
			}

			var refType int

			if authInfo.ClientCertificate == *str {
				refType = 0
			} else if authInfo.ClientKey == *str {
				refType = 1
			} else if authInfo.TokenFile == *str {
				refType = 2
			}

			fileBytes, err := ioutil.ReadFile(*str)

			if err != nil {
				continue
			}

			if refType == 0 {
				authInfo.ClientCertificateData = fileBytes
				authInfo.ClientCertificate = ""
			} else if refType == 1 {
				authInfo.ClientKeyData = fileBytes
				authInfo.ClientKey = ""
			} else if refType == 2 {
				authInfo.Token = base64.StdEncoding.EncodeToString(fileBytes)
				authInfo.TokenFile = ""
			}
		}
	}
}

func populateOIDCPluginCerts(config *clientcmdapi.Config) {
	for _, authInfo := range config.AuthInfos {
		if authInfo.AuthProvider != nil && authInfo.AuthProvider.Name == "oidc" {
			if ca, ok := authInfo.AuthProvider.Config["idp-certificate-authority"]; ok && ca != "" {
				fileBytes, err := ioutil.ReadFile(ca)

				if err != nil {
					continue
				}

				authInfo.AuthProvider.Config["idp-certificate-authority-data"] = base64.StdEncoding.EncodeToString(fileBytes)
				delete(authInfo.AuthProvider.Config, "idp-certificate-authority")
			}
		}
	}
}

// createAllowedContextMap creates a dummy map from context name to context name
func createAllowedContextMap(contexts []string) map[string]string {
	aContextMap := make(map[string]string)

	for _, context := range contexts {
		aContextMap[context] = context
	}

	return aContextMap
}
