Switchboard is a control plane for cloud-native resources. It acts as a management layer for infrastructure and application configuration, built on top of tools like Terraform, Helm, and Kubernetes. 

When used on its own, it serves as a mechanism to simplify complex workflows. When used in conjunction with the Porter API and dashboard, it enables you to build an entirely self-serve platform where developers can:

1. Programmatically create new environments for new infrastructure or applications
2. Interact with exposed configuration variables for infrastructure, through a dashboard, CLI, or API

Simultaneously, it allows platform teams to:

1. Enforce RBAC at any layer, across the entire stack 
2. Support multiple workflows, like GitOps or two-step deploys (?) and adopt flexible ownership models
3. Inject secrets and sensitive data using a set of prebuilt injectors
4. Works with any cloud

Let's dive into an example! We'll see how two of the most common DevOps tools, Terraform and Helm, are natively supported on Switchboard to provide developers with a seamless experience. Let's say a team wants to roll out a new production service on an existing Kubernetes cluster, hosted on AWS, which has the following requirements:
- A dedicated RDS instance, provisioned in an existing VPC
- The creation of a Helm chart for the application, which accepts RDS credentials as environment variables 

Without Switchboard, a typical deployment process might look something like:
1. Task an infrastructure engineer with modifying the VPC to include database CIDRs and call `terraform apply` to create the RDS instance.
2. Infrastructure engineer writes RDS credentials to a secret vault, and notifies the application developer or DevOps team that the RDS instance is ready. 
3. Application developer or DevOps teams writes a `values.yaml` file for the Helm chart, and pushes this code to a shared Github repository containing application manifests. 
4. Applications manifests are applied to the cluster through a GitOps approach, using a tool like ArgoCD. 
5. Developer can view the application through either the ArgoCD dashboard, or monitor the application through an APM tool or custom dashboard in Grafana. 

With Switchboard, the application developer can own the entire flow, with optional input or approval from infrastructure engineers and/or system admins. Here's what the equivalent process looks like on Switchboard!

```
# TODO: video demo
# shows VPC + RDS + application deployment
```


