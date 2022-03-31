package crossplane

import (
	"testing"

	"github.com/arttor/helmify/internal"
	"github.com/arttor/helmify/pkg/metadata"
	"github.com/stretchr/testify/assert"
)

const (
	strComposition = `apiVersion: apiextensions.crossplane.io/v1
kind: Composition
metadata:
  name: default
  creationTimestamp: null
  annotations:
    controller-gen.kubebuilder.io/version: v0.7.0
spec:
  compositeTypeRef:
    apiVersion: olap.kyligence.io/v1alpha1
    kind: KyligenceRDSInstance
`
)

func Test_composition_Process(t *testing.T) {
	var testInstance composition

	t.Run("processed", func(t *testing.T) {
		obj := internal.GenerateObj(strComposition)
		processed, _, err := testInstance.Process(metadata.New("cloud-operator"), obj)
		assert.NoError(t, err)
		assert.Equal(t, true, processed)
	})
	t.Run("skipped", func(t *testing.T) {
		obj := internal.TestNs
		processed, _, err := testInstance.Process(&metadata.Service{}, obj)
		assert.NoError(t, err)
		assert.Equal(t, false, processed)
	})
}
