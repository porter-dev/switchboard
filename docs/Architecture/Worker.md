## Overview

The worker is responsible for reading a `ResourceGroup` and managing the specified drivers to apply the resource. The worker parses the `ResourceGroup` into an `ExecTree`. It does this by computing all dependencies (either `symbolic` or `declared`) and detecting things like circular dependencies. It then goes through each node in the execution tree to determine if the node needs to be re-applied (or destroyed). If a branch of the `ExecTree` does not need to be re-applied, no operation is performed on that branch. 

## Apply Change Detection

`Apply` operations are checked against a previous run to determine if the worker needs to take action for a certain node in the `ExecTree`. Once an `Apply` operation does take place on a `Node`, the operation is idempotent: it performs the exact same action on each run. At the point where a resource is applied, change detection is the *responsibility of the underlying driver*. Each driver can handle drift detection and reconciliation in its own way. 

NOTE: this is the case, but go into more detail on three-way drift detection. Each `ShouldApply` is a way to determine if there are patches to be made. Some more details: https://kubernetes.io/docs/tasks/manage-kubernetes-objects/declarative-config/#merge-patch-calculation

To highlight this difference, let's say the resource is a `terraform` `rds` module, in which the underlying TF module declares 20 Terraform resources. In this case, the Switchboard worker will pass the same set of `config` values to Terraform, and that Terraform module will check against its own state to determine which of those 20 resources should be created, updated, deleted, etc. 

This means that the worker state is *only used to speed up apply operations*. Thus, the worker's previous state is compared against the current `ResourceGroup`, and the following checks are run:
1. If the `source` has changed. This varies depending on the source type, but generally for registries and Github we check the references, and some drivers compute hashes of local files. 
2. If the `target` has changed. If the `target` changes, we apply the new target, and once the new target has been created, we destroy the old target. 
3. If the `config` has changed. We verify this using cryptographically secure hashes.


## Failed Apply or Destroy Operations

The worker will try to handle failed Apply operations gracefully according to the following rules:
1. Apply operations for failed resources will always be retried.
2. If a parent fails, child resources will not be applied. 
3. If an Apply operation fails and the target is changed, failure to remove the previous target is not a fatal error. 
4. Execution will not halt until all execution branches have been attempted. 
