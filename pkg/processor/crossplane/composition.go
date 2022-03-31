package crossplane

import (
	"io"
	"strings"

	"github.com/arttor/helmify/pkg/helmify"
	yamlformat "github.com/arttor/helmify/pkg/yaml"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var compositionGVC = schema.GroupVersionKind{
	Group:   "apiextensions.crossplane.io",
	Version: "v1",
	Kind:    "Composition",
}

// NewComposition creates processor for k8s Composition resource.
func NewComposition() helmify.Processor {
	return &composition{}
}

type composition struct {
}

func (c composition) Process(appMeta helmify.AppMetadata, obj *unstructured.Unstructured) (bool, helmify.Template, error) {
	if obj.GroupVersionKind() != compositionGVC {
		return false, nil, nil
	}
	name := obj.GetName()
	name = strings.Replace(name, appMeta.ChartName()+"-", "", -1)
	obj.SetName(name)

	obj.SetNamespace("")

	body, err := yamlformat.Marshal(obj.Object, 0)
	if err != nil {
		return true, nil, err
	}
	return true, &compositionresult{
		name: obj.GetName() + "-composition.yaml",
		data: []byte(body),
	}, nil
}

type compositionresult struct {
	name string
	data []byte
}

func (r *compositionresult) Filename() string {
	return r.name
}

func (r *compositionresult) Values() helmify.Values {
	return helmify.Values{}
}

func (r *compositionresult) Write(writer io.Writer) error {
	_, err := writer.Write(r.data)
	return err
}
