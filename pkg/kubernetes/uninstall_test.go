package kubernetes

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"

	"get.porter.sh/porter/pkg/test"
	"github.com/stretchr/testify/require"

	"gopkg.in/yaml.v2"
)

type UnInstallTest struct {
	expectedCommand string
	uninstallStep   UninstallStep
}

func TestMixin_UninstallStep(t *testing.T) {

	manifestDirectory := "/cnab/app/manifests"

	deleteCmd := "kubectl delete --ignore-not-found=true -f"

	dontWait := false

	namespace := "meditations"

	selector := "app=nginx"
	k8scontext := "context"
	kubeConfig := "$HOME/.kube/config"
	forceIt := true
	withGrace := 1

	timeout := 1

	uninstallTests := []UnInstallTest{
		{
			expectedCommand: fmt.Sprintf("%s %s --wait", deleteCmd, manifestDirectory),
			uninstallStep: UninstallStep{
				UninstallArguments: UninstallArguments{
					Step: Step{
						Description: "Hello",
					},
					Manifests: []string{manifestDirectory},
				},
			},
		},
		{
			expectedCommand: fmt.Sprintf("%s %s", deleteCmd, manifestDirectory),
			uninstallStep: UninstallStep{
				UninstallArguments: UninstallArguments{
					Step: Step{
						Description: "Hello",
					},
					Manifests: []string{manifestDirectory},
					Wait:      &dontWait,
				},
			},
		},
		{
			expectedCommand: fmt.Sprintf("%s %s -n %s", deleteCmd, manifestDirectory, namespace),
			uninstallStep: UninstallStep{
				UninstallArguments: UninstallArguments{
					Step: Step{
						Description: "Hello",
					},
					Manifests: []string{manifestDirectory},
					Namespace: namespace,
					Wait:      &dontWait,
				},
			},
		},
		{
			expectedCommand: fmt.Sprintf("%s %s --selector=%s --wait", deleteCmd, manifestDirectory, selector),
			uninstallStep: UninstallStep{
				UninstallArguments: UninstallArguments{
					Step: Step{
						Description: "Hello",
					},
					Manifests: []string{manifestDirectory},
					Selector:  selector,
				},
			},
		},
		{
			expectedCommand: fmt.Sprintf("%s %s --context=%s --wait", deleteCmd, manifestDirectory, k8scontext),
			uninstallStep: UninstallStep{
				UninstallArguments: UninstallArguments{
					Step: Step{
						Description: "Hello",
					},
					Manifests: []string{manifestDirectory},
					Context:   k8scontext,
				},
			},
		},
		{
			expectedCommand: fmt.Sprintf("%s %s --kubeconfig=%s --wait", deleteCmd, manifestDirectory, kubeConfig),
			uninstallStep: UninstallStep{
				UninstallArguments: UninstallArguments{
					Step: Step{
						Description: "Hello",
					},
					Manifests:  []string{manifestDirectory},
					KubeConfig: kubeConfig,
				},
			},
		},
		{
			expectedCommand: fmt.Sprintf("%s %s --force --grace-period=0 --wait", deleteCmd, manifestDirectory),
			uninstallStep: UninstallStep{
				UninstallArguments: UninstallArguments{
					Step: Step{
						Description: "Hello",
					},
					Manifests: []string{manifestDirectory},
					Force:     &forceIt,
				},
			},
		},
		{
			expectedCommand: fmt.Sprintf("%s %s --grace-period=%d --wait", deleteCmd, manifestDirectory, withGrace),
			uninstallStep: UninstallStep{
				UninstallArguments: UninstallArguments{
					Step: Step{
						Description: "Hello",
					},
					Manifests:   []string{manifestDirectory},
					GracePeriod: &withGrace,
				},
			},
		},
		{
			expectedCommand: fmt.Sprintf("%s %s --timeout=%ds --wait", deleteCmd, manifestDirectory, timeout),
			uninstallStep: UninstallStep{
				UninstallArguments: UninstallArguments{
					Step: Step{
						Description: "Hello",
					},
					Manifests: []string{manifestDirectory},
					Timeout:   &timeout,
				},
			},
		},
		{
			expectedCommand: fmt.Sprintf("%s %s --wait --cascade=background", deleteCmd, manifestDirectory),
			uninstallStep: UninstallStep{
				UninstallArguments: UninstallArguments{
					Step: Step{
						Description: "Hello",
					},
					Manifests: []string{manifestDirectory},
					Cascade:   "background",
				},
			},
		},
	}

	defer os.Unsetenv(test.ExpectedCommandEnv)
	for _, uninstallTest := range uninstallTests {
		t.Run(uninstallTest.expectedCommand, func(t *testing.T) {
			ctx := context.Background()
			h := NewTestMixin(t)
			h.Setenv(test.ExpectedCommandEnv, uninstallTest.expectedCommand)

			action := UninstallAction{Steps: []UninstallStep{uninstallTest.uninstallStep}}
			b, _ := yaml.Marshal(action)

			h.In = bytes.NewReader(b)

			err := h.Uninstall(ctx)

			require.NoError(t, err)
		})
	}
}
