package terraform

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/porter-dev/switchboard/internal/models"
	"github.com/porter-dev/switchboard/pkg/drivers"

	hcljson "github.com/hashicorp/hcl2/hcl/json"
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

	if source.VarMethod == VarMethodFile {
		// construct the var file path
		if filepath.IsAbs(source.Path) {
			driver.varFilePath = filepath.Join(source.Path, "tfvars.json")
		} else {
			driver.varFilePath = filepath.Join(opts.BaseDir, source.Path, "tfvars.json")
		}
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

	err = d.tf.Init(context.Background())

	if err != nil {
		return nil, err
	}

	applyOpts, err := d.getVarOpts(config)

	if err != nil {
		return nil, err
	}

	err = d.tf.Apply(context.Background(), applyOpts...)

	if err != nil {
		return nil, err
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

// setVars sets variables for the Terraform process through
// either a var file or env variables.
func (d *Driver) getVarOpts(config map[string]interface{}) ([]tfexec.ApplyOption, error) {
	applyOpts := make([]tfexec.ApplyOption, 0)

	if d.source.VarMethod == VarMethodFile {
		applyOpts = append(applyOpts, tfexec.VarFile(d.varFilePath))
	}

	switch d.source.VarMethod {
	case VarMethodEnv:
		for key, val := range config {
			valBytes, err := toHCL(val)

			if err != nil {
				return nil, err
			}

			applyOpts = append(applyOpts, tfexec.Var(fmt.Sprintf("%s=%s", key, string(valBytes))))
		}
	case VarMethodFile:
		file, err := json.Marshal(config)

		if err != nil {
			return nil, err
		}

		err = ioutil.WriteFile(d.varFilePath, file, 0644)

		if err != nil {
			return nil, err
		}

		applyOpts = append(applyOpts, tfexec.VarFile(d.varFilePath))
	}

	return applyOpts, nil
}

func toHCL(val interface{}) ([]byte, error) {
	jsonValBytes, err := json.Marshal(val)

	if err != nil {
		return []byte{}, err
	}

	switch val.(type) {
	case map[string]interface{}:
	case []interface{}:
		// this is pretty annoying, but we first marshal into json and
		// from there HCL
		hclFile, err := hcljson.Parse(jsonValBytes, "")

		if err != nil {
			return []byte{}, err
		}

		return hclFile.Bytes, nil
	}

	// in the default case (string, int), we just return raw json,
	// it should be valid
	return jsonValBytes, nil
}
