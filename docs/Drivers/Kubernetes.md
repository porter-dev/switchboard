## Source 
The Kubernetes driver can specify a source to pull the base configuration from. For Kubernetes objects, the base configuration must be a valid Kubernetes resource definition, written as a yaml file. This base configuration can be read from:
1. A local file
2. A Github repo, either public or private
3. A Porter Github integration
4. A Porter Manifest integration

### Local

```yaml
source:
  kind: local
  path: /path/to/local/manifest.yaml
```

### Github Repo

```yaml
source:
  kind: github
  repo: porter-dev/porter
  path: ./relative/path.yaml
  token: optional-token
```

### Porter Github Integration

```yaml
target:
  kind: porter
  project: 1234
  github: 4321
```

### Porter Manifest Integration

```yaml
target:
  kind: porter
  project: 1234
  name: my-custom-manifest
```

## Target 
The Kubernetes driver needs a cluster and optionally a namespace as the target to apply a new resource. This Kubernetes target can be determined from 3 sources:
1. The local kubeconfig
2. In-cluster configuration 
3. Porter cluster model

### Local

```yaml
target:
  kind: local
  context: custom-context-name
  path: /custom/path/to/kubeconfig
```

### In-Cluster Configuration

```yaml
target:
  kind: in-cluster
```

### Porter Cluster Model

```yaml
target:
  kind: porter
  project: 1234
  cluster: 4321
```

## Config 

The `config` section for the Kubernetes driver supports very basic variable override. For custom variable override, a more complex templating engine should be used, such as Helm. 

Each variable set in the `config` section will override a "base" variable defined in the `source` section. Thus, a config section might look like the following:

```yaml
config:
  metadata:
    name: user-defined-name
	namespace: custom-namespace
```

## Examples

### Simple Deployment

In the Github repo `github.com/porter-dev/kubernetes-driver-examples`, let's say we have the following base file:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
	    app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
```

A resource group that sets the number of replicas from 3 -> 5 may look like:

```yaml
version: v1
resources:
- name: nginx-deployment
  target:
    kind: local
	path: /path/to/local/manifest
  source:
    kind: github
    repo: porter-dev/porter
    path: ./nginx-deployment.yaml
    token: optional-token
  config:
    spec:
	  replicas: 5
```