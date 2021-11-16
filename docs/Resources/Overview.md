## Overview

Resources represent underlying objects on the target platform -- like Helm charts, Kubernetes deployments, or Terraform modules. Resources are grouped as a `ResourceGroup`, which are a collection of cloud resources that are meant to be deployed together. For example, if an application (represented as a Helm chart) depends on an RDS database and an S3 bucket, there would be three resources:
1. An RDS database, created using the `terraform` driver. 
2. An S3 bucket, created using the `terraform` driver. 
3. A Helm chart that depends on both the RDS and S3 bucket, created using the `helm` driver. 


### Dependencies

Resources can be dependent on other resources, which means that the parent resources get applied before the child resource. Dependencies are computed using two mechanisms:

- `explicit` declarations use the `depends_on` field to declare parent -> child relationships
- `implicit` declarations use [[Resources/Overview#Variable Injection|variable injection]] to determine parent -> child relationships

Internally, the dependency graph is represented as a directed acyclic graph. Any dependency cycles will return an error before apply.

### Variable Injection

Each resource has an output that can be referenced by a [[Resources/Overview#Dependencies|dependent resource]]. Any value in the `config` section of the resource can be set via variable injection from a different resource. There are currently 2 supported query languages:

- `jsonpath` queries
- `jq` queries

For an example, take a look at [[Resource Reference#RDS Helm Chart|this resource group]], in which a dependent application reads data from an RDS resource.