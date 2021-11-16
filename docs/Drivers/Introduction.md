## Introduction

Drivers are mechanisms to interface with underlying tools, like Kubernetes APIs, Terraform, and Helm. The primary function of drivers is to extend these tools to make difficult tasks much simpler: tasks such as implementing Git-based workflows, role-based access control, and secrets management. 

All drivers are configured through three primary fields: `source`, `target`, and `config`:
- `source` represents the source of the base configuration or templates used when creating a resource. 
- `target` represents where the resource should be created, such as an AWS account, a Kubernetes cluster, or a Docker registry. 
- `config` represents the input values passed to the templating engine. 

This may sound a bit abstract: for a better sense of these fields, let's consider a simple example of reading a set of Kubernetes manifests from a Github repository and applying these manifests to a cluster. Here's an example Porter resource group that creates an NGINX deployment on a Kubernetes cluster: 

```yaml
version: v1
resources:
- name: nginx-deployment
  driver: kubernetes
  source:
    kind: github
    repo: porter-dev/kubernetes-driver-examples
    path: ./nginx-deployment.yaml
  target:
    kind: local_kubeconfig
  config:
    spec:
	  replicas: 5
```

In this case, these fields represent the following:
- `source` declares that the base template should be read from the file path `./nginx-deployment.yaml` in the `github.com/porter-dev/kubernetes-driver-examples` repository. 
- `target` declares that the resource should be created on the cluster which is declared in the local kubeconfig's current context
- `config` declares that we should overwrite the number of replicas declared in github.com/porter-dev/kubernetes-driver-examples from 3 replicas -> 5 replicas. 

If we run `porter apply -f apply-example.yaml`, we'll see the resource created! 

We already see how this might make implementing a Git-based workflow much simpler. The `kubernetes` driver has a built-in mechanism to read from a Github filesystem, and has support for specifying a Github reference, like a tag or branch. 