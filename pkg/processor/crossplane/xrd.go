package crossplane

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/arttor/helmify/pkg/helmify"
	yamlformat "github.com/arttor/helmify/pkg/yaml"
	v1 "github.com/crossplane/crossplane/apis/apiextensions/v1"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/yaml"
)

const xrdTeml = `apiVersion: apiextensions.crossplane.io/v1
kind: CompositeResourceDefinition
metadata:
  name: %[1]s
  annotations:
    controller-gen.kubebuilder.io/version: v0.7.0
  labels:
  {{- include "%[2]s.labels" . | nindent 4 }}
spec:
%[3]s`

var xrdGVC = schema.GroupVersionKind{
	Group:   "apiextensions.crossplane.io",
	Version: "v1",
	Kind:    "CompositeResourceDefinition",
}

// NewXRD creates processor for k8s CompositeResourceDefinition resource.
func NewXRD() helmify.Processor {
	return &xrd{}
}

type xrd struct {
}

func (x xrd) Process(appMeta helmify.AppMetadata, obj *unstructured.Unstructured) (bool, helmify.Template, error) {
	if obj.GroupVersionKind() != xrdGVC {
		return false, nil, nil
	}

	name := obj.GetName()
	name = strings.Replace(name, appMeta.ChartName()+"-", "", -1)
	obj.SetName(name)

	specUnstr, ok, err := unstructured.NestedMap(obj.Object, "spec")
	if err != nil || !ok {
		return true, nil, errors.Wrap(err, "unable to create xrd template")
	}

	spec := v1.CompositeResourceDefinitionSpec{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(specUnstr, &spec)
	if err != nil {
		return true, nil, errors.Wrap(err, "unable to cast to xrd spec")
	}

	versions, _ := yaml.Marshal(spec)
	versions = yamlformat.Indent(versions, 2)
	versions = bytes.TrimRight(versions, "\n ")

	res := fmt.Sprintf(xrdTeml, obj.GetName(), appMeta.ChartName(), string(versions))
	return true, &xrdresult{
		name: obj.GetName() + "-xrd.yaml",
		data: []byte(res),
	}, nil
}

type xrdresult struct {
	name string
	data []byte
}

func (r *xrdresult) Filename() string {
	return r.name
}

func (r *xrdresult) Values() helmify.Values {
	return helmify.Values{}
}

func (r *xrdresult) Write(writer io.Writer) error {
	_, err := writer.Write(r.data)
	return err
}
