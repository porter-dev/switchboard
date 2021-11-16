DONE:
- Create simple example of resource in a file, read it in, and apply it to a cluster -- Kubernetes driver
	- Finish `source` config with kind `local`, param `path`
	- Finish `target` config with kind `local`, param `context` and `namespace`
	- Clean up code + remove unnecessary code, repackage how clientsets are loaded and such. For example, `local` package shouldn't return an Agent or depend on kubernetes. Can write an interface for loading based on the target context. 

DONE:
- Simple graph-based execution model with parallelism. Will refine more later

DONE:
- Get test running that uses output data from dependent sources, and we parse those other sources. Can be really simple query model that just supports nested mapped fields and slices for now. 
	- Can use jsonpath: https://pkg.go.dev/k8s.io/client-go/util/jsonpath#New -- looks like we need to register New, call Parse, and then call FindResults. 
- Use output data from previous resources. Need a model where output data can be queried. Perhaps we can reuse the jq query engine here?

TODO:
- 3-way strategic merge patches with reconciliation -- need to write docs and re-write Kubernetes apply operation to use this.
- Work on other two lifecycle commands:
	- Update 
	- Delete
- Write the Helm driver, should be pretty quick for local 
- Formalize driver interface, and rewrite main program to work from driver interface