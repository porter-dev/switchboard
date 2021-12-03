package kubernetes

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"

	"context"
	"fmt"

	"github.com/porter-dev/switchboard/utils/objutils"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Agent struct {
	RESTClientGetter genericclioptions.RESTClientGetter
	Clientset        kubernetes.Interface
	DynamicClientset dynamic.Interface
}

type ApplyOpts struct {
	Config map[string]interface{}
	Base   map[string]interface{}
	Target *Target
}

func (a *Agent) Apply(opts *ApplyOpts) (map[string]interface{}, error) {
	// override the base config with the specified resource's config
	obj := objutils.CoalesceValues(opts.Base, opts.Config)
	gvr, err := a.getGroupVersionResource(obj)

	if err != nil {
		return nil, fmt.Errorf("could not get API group, version, or resource: %v", err)
	}

	dynResource := a.DynamicClientset.Resource(*gvr).Namespace(opts.Target.Namespace)

	// attempt to get the resource
	name, err := getObjectName(obj)

	if err != nil {
		return nil, fmt.Errorf("could not get object name: %v", err)
	}

	_, err = dynResource.Get(context.TODO(), name, metav1.GetOptions{})
	var res map[string]interface{}

	// check if the error is a resource "NotFound" error
	if err != nil && errors.IsNotFound(err) {
		// create the resource
		unstructObj, err := dynResource.Create(context.TODO(), &unstructured.Unstructured{
			Object: obj,
		}, metav1.CreateOptions{})

		if err != nil {
			return nil, err
		}

		res = unstructObj.Object
	} else if err != nil {
		return nil, fmt.Errorf("error getting the resource: %v", err)
	} else {
		// update the resource
		unstructObj, err := dynResource.Update(context.TODO(), &unstructured.Unstructured{
			Object: obj,
		}, metav1.UpdateOptions{})

		if err != nil {
			return nil, err
		}

		res = unstructObj.Object
	}

	return res, nil
}

func (a *Agent) getGroupVersionResource(obj map[string]interface{}) (*schema.GroupVersionResource, error) {
	// get the apiVersion and kind from the object
	apiVersion, apiVersionExists := obj["apiVersion"]

	if !apiVersionExists {
		return nil, fmt.Errorf("apiVersion field must be set")
	}

	apiVersionStr, ok := apiVersion.(string)

	if !ok {
		return nil, fmt.Errorf("apiVersion field is not a string")
	}

	kind, kindExists := obj["kind"]

	if !kindExists {
		return nil, fmt.Errorf("kind field must be set")
	}

	kindStr, ok := kind.(string)

	if !ok {
		return nil, fmt.Errorf("kind field is not a string")
	}

	// parse the object for the object group, version, kind
	gvk := schema.FromAPIVersionAndKind(apiVersionStr, kindStr)

	// use the gvk and restmapper to construct a gvr
	restMapper, err := a.RESTClientGetter.ToRESTMapper()

	if err != nil {
		return nil, fmt.Errorf("error in REST mapper: %v", err)
	}

	mapping, err := restMapper.RESTMapping(schema.GroupKind{
		Group: gvk.Group,
		Kind:  gvk.Kind,
	}, gvk.Version)

	if err != nil {
		return nil, fmt.Errorf("error creating a REST mapping: %v", err)
	}

	return &mapping.Resource, nil
}

func getObjectName(obj map[string]interface{}) (string, error) {
	metadataInt, metadataExists := obj["metadata"]

	if !metadataExists {
		return "", fmt.Errorf("metadata field must be set")
	}

	metadata, ok := metadataInt.(map[string]interface{})

	if !ok {
		return "", fmt.Errorf("metadata field could not be converted to a map[string]interface{}")
	}

	nameInt, nameExists := metadata["name"]

	if !nameExists {
		return "", fmt.Errorf("name field must be set")
	}

	name, ok := nameInt.(string)

	if !ok {
		return "", fmt.Errorf("metadata.name field is not a string")
	}

	return name, nil
}
