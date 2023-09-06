package main

import (
	"io/ioutil"
	"os"

	"github.com/porter-dev/switchboard/v2/pkg/parser"
	"github.com/porter-dev/switchboard/v2/pkg/validator"
)

func main() {
	bytes, err := ioutil.ReadFile(os.Args[1])

	if err != nil {
		panic(err)
	}

	parsed, err := parser.ParseRawBytes(bytes)

	if err != nil {
		panic(err)
	}

	err = validator.ValidatePorterYAML(parsed)

	if err != nil {
		panic(err)
	}
}
