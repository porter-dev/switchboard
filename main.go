package main

import (
	"io/ioutil"
	"os"

	"github.com/porter-dev/switchboard/pkg/drivers/helm"
	"github.com/porter-dev/switchboard/pkg/drivers/kubernetes"
	"github.com/porter-dev/switchboard/pkg/drivers/terraform"
	"github.com/porter-dev/switchboard/pkg/parser"
	"github.com/porter-dev/switchboard/pkg/types"
	"github.com/porter-dev/switchboard/pkg/worker"
)

func main() {
	// read the local resource bytes
	fileBytes, err := ioutil.ReadFile("./test-resource-1.yaml")

	if err != nil {
		panic(err)
	}

	resGroup, err := parser.ParseRawBytes(fileBytes)

	if err != nil {
		panic(err)
	}

	basePath, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	worker := worker.NewWorker()
	worker.RegisterDriver("helm", helm.NewHelmDriver)
	worker.RegisterDriver("kubernetes", kubernetes.NewKubernetesDriver)
	worker.RegisterDriver("terraform", terraform.NewTerraformDriver)

	err = worker.Apply(resGroup, &types.ApplyOpts{
		BasePath: basePath,
	})

	if err != nil {
		panic(err)
	}
}
