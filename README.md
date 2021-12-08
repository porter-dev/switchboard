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

## Hooks

Hooks can be added to the worker when calling the package:

```go

type TestHook struct{}

func (t *TestHook) PreApply() error {
	fmt.Println("RUNNING PRE APPLY")
	return nil
}

func (t *TestHook) DataQueries() map[string]interface{} {
	return map[string]interface{}{
		"first": "{ .test-deployment.spec.replicas }",
	}
}

func (t *TestHook) PostApply(populatedData map[string]interface{}) error {
	fmt.Println("POPULATED DATA IS", populatedData)
	return nil
}

```

Registered via:

```go
worker.RegisterHook("test", &TestHook{})
```
