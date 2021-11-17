package terraform

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/porter-dev/switchboard/internal/models"
	"github.com/porter-dev/switchboard/pkg/drivers"
)

type Driver struct {
	source      *Source
	output      map[string]interface{}
	lookupTable *map[string]drivers.Driver
	varFilePath string
	tf          *tfexec.Terraform
}

func NewTerraformDriver(resource *models.Resource, opts *drivers.SharedDriverOpts) (*Driver, error) {
	driver := &Driver{
		lookupTable: opts.DriverLookupTable,
	}

	source, err := GetSource(resource.Source)

	if err != nil {
		return nil, err
	}

	err = driver.initSource(source, opts)

	if err != nil {
		return nil, err
	}

	driver.source = source

	// construct the var file path
	if filepath.IsAbs(source.Path) {
		driver.varFilePath = filepath.Join(source.Path, "tfvars.json")
	} else {
		driver.varFilePath = filepath.Join(opts.BaseDir, source.Path, "tfvars.json")
	}

	return driver, nil
}

func (d *Driver) initSource(source *Source, opts *drivers.SharedDriverOpts) error {
	// read the file and set the base variable
	switch source.Kind {
	case SourceKindLocal:
		tf, err := tfexec.NewTerraform(source.SourceLocal.Path, "terraform")

		if err != nil {
			log.Fatalf("error running NewTerraform: %s", err)
		}

		// TODO: don't set these to os stdout or stderr necessary, we probably want a json parser
		// of sorts
		// tf.SetStdout(os.Stdout)
		tf.SetStderr(os.Stderr)

		d.tf = tf
	}

	return nil
}

func (d *Driver) ShouldApply(resource *models.Resource) bool {
	return true
}

func (d *Driver) Apply(resource *models.Resource) (*models.Resource, error) {
	config, err := drivers.ConstructConfig(&drivers.ConstructConfigOpts{
		RawConf:      resource.Config,
		LookupTable:  *d.lookupTable,
		Dependencies: resource.Dependencies,
	})

	if err != nil {
		return nil, err
	}

	// TODO: write config as var file in local path directory
	file, err := json.Marshal(config)

	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(d.varFilePath, file, 0644)

	if err != nil {
		return nil, err
	}

	err = d.tf.Init(context.Background())
	if err != nil {
		log.Fatalf("error running Init: %s", err)
	}

	err = d.tf.Apply(context.Background(), tfexec.VarFile(d.varFilePath))

	if err != nil {
		log.Fatalf("error running Apply: %s", err)
	}

	return resource, nil
}

// Output returns the created TF output
func (d *Driver) Output() (map[string]interface{}, error) {
	output, err := d.tf.Output(context.Background())

	if err != nil {
		return nil, err
	}

	keyToVal := make(map[string]interface{})

	for key, meta := range output {
		keyToVal[key] = meta.Value
	}

	// not the most efficient, but marshal to json and decode again
	res := make(map[string]interface{})

	rawBytes, err := json.Marshal(keyToVal)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(rawBytes, &res)

	if err != nil {
		return nil, err
	}

	return res, nil
}
