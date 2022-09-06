package job

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/arttor/helmify/pkg/metadata"

	"github.com/arttor/helmify/internal"
)

const (
	strDepl = `apiVersion: batch/v1
kind: Job
metadata:
  name: admission-init
  namespace: system
  labels:
    app: admission-init
spec:
  backoffLimit: 3
  template:
    spec:
      imagePullSecrets:
        - name: ""
      serviceAccountName: admission-manager
      priorityClassName: system-cluster-critical
      restartPolicy: Never
      containers:
        - image: admission:latest
          imagePullPolicy: Always
          name: admission
          command: ["./gen-admission-secret.sh", "--service", "admission-service", "--namespace",
                    "system", "--secret", "admission-secret"]
`
)

func Test_job_Process(t *testing.T) {
	var testInstance job

	t.Run("processed", func(t *testing.T) {
		obj := internal.GenerateObj(strDepl)
		processed, _, err := testInstance.Process(&metadata.Service{}, obj)
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
