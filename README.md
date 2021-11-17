## Quick Start

Start by building the Switchboard binary:

```
make build-cli-dev
```

Then run an example TF and Helm deployment:

```
./bin/switchboard apply ./examples/terraform/test-resource-1.yaml
```

This will deploy a TF resource with a dependent deployment, and should print something like this to the console:

```
INF running apply for resource rds
INF successfully applied resource rds
INF running apply for resource tf-deployment
INF successfully applied resource tf-deployment
```
