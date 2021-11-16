package main

import (
	"io/ioutil"
	"os"

	"github.com/porter-dev/switchboard/pkg/parser"
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

	err = worker.Apply(resGroup, &worker.ApplyOpts{
		BasePath: basePath,
	})

	if err != nil {
		panic(err)
	}
}
