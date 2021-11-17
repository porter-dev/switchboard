package terraform

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-exec/tfexec"
)

type Agent struct {
	//
}

func Apply() {
	workingDir := "/path/to/working/dir"
	tf, err := tfexec.NewTerraform(workingDir, "")
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	}

	err = tf.Init(context.Background())
	if err != nil {
		log.Fatalf("error running Init: %s", err)
	}

	err = tf.Apply(context.Background())

	if err != nil {
		log.Fatalf("error running Apply: %s", err)
	}
}
