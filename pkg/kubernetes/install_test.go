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

type InstallTest struct {
	expectedCommand string
	installStep     InstallStep
}

func TestMain(m *testing.M) {
	test.TestMainWithMockedCommandHandlers(m)
}

func TestMixin_InstallStep(t *testing.T) {

	manifestDirectory := "/cnab/app/manifests"

	installCmd := "kubectl apply -f"

	dontWait := false
	recordIt := true
	validateIt := false
	serverSide := true
	forceConflicts := true

	namespace := "meditations"

	selector := "app=nginx"
	k8scontext := "context"
	kubeConfig := "$HOME/.kube/config"
	installTests := []InstallTest{
		{
			expectedCommand: fmt.Sprintf("%s %s --wait --server-side=true", installCmd, manifestDirectory),
			installStep: InstallStep{

				InstallArguments: InstallArguments{
					Step: Step{
						Description: "Hello",
					},
					Manifests:  []string{manifestDirectory},
					ServerSide: &serverSide,
				},
			},
		},
		{
			expectedCommand: fmt.Sprintf("%s %s --wait --server-side=true --force-conflicts=true", installCmd, manifestDirectory),
			installStep: InstallStep{

				InstallArguments: InstallArguments{
					Step: Step{
						Description: "Hello",
					},
					Manifests:      []string{manifestDirectory},
					ServerSide:     &serverSide,
					ForceConflicts: &forceConflicts,
				},
			},
		},
		{
			expectedCommand: fmt.Sprintf("%s %s --wait --force-conflicts=true", installCmd, manifestDirectory),
			installStep: InstallStep{

				InstallArguments: InstallArguments{
					Step: Step{
						Description: "Hello",
					},
					Manifests:      []string{manifestDirectory},
					ForceConflicts: &forceConflicts,
				},
			},
		},
		{
			expectedCommand: fmt.Sprintf("%s %s --wait", installCmd, manifestDirectory),
			installStep: InstallStep{

				InstallArguments: InstallArguments{
					Step: Step{
						Description: "Hello",
					},
					Manifests: []string{manifestDirectory},
				},
			},
		},
		{
			expectedCommand: fmt.Sprintf("%s %s", installCmd, manifestDirectory),
			installStep: InstallStep{
				InstallArguments: InstallArguments{
					Step: Step{
						Description: "Hello",
					},
					Manifests: []string{manifestDirectory},
					Wait:      &dontWait,
				},
			},
		},
		{
			expectedCommand: fmt.Sprintf("%s %s -n %s", installCmd, manifestDirectory, namespace),
			installStep: InstallStep{
				InstallArguments: InstallArguments{
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
			expectedCommand: fmt.Sprintf("%s %s -n %s --validate=false", installCmd, manifestDirectory, namespace),
			installStep: InstallStep{
				InstallArguments: InstallArguments{
					Step: Step{
						Description: "Hello",
					},
					Manifests: []string{manifestDirectory},
					Namespace: namespace,
					Validate:  &validateIt,
					Wait:      &dontWait,
				},
			},
		},
		{
			expectedCommand: fmt.Sprintf("%s %s -n %s --record=true", installCmd, manifestDirectory, namespace),
			installStep: InstallStep{
				InstallArguments: InstallArguments{
					Step: Step{
						Description: "Hello",
					},
					Manifests: []string{manifestDirectory},
					Namespace: namespace,
					Record:    &recordIt,
					Wait:      &dontWait,
				},
			},
		},
		{
			expectedCommand: fmt.Sprintf("%s %s --selector=%s --wait", installCmd, manifestDirectory, selector),
			installStep: InstallStep{
				InstallArguments: InstallArguments{
					Step: Step{
						Description: "Hello",
					},
					Manifests: []string{manifestDirectory},
					Selector:  selector,
				},
			},
		},
		{
			expectedCommand: fmt.Sprintf("%s %s --context=%s --wait", installCmd, manifestDirectory, k8scontext),
			installStep: InstallStep{
				InstallArguments: InstallArguments{
					Step: Step{
						Description: "Hello",
					},
					Manifests: []string{manifestDirectory},
					Context:   k8scontext,
				},
			},
		},
		{
			expectedCommand: fmt.Sprintf("%s %s --kubeconfig=%s --wait", installCmd, manifestDirectory, kubeConfig),
			installStep: InstallStep{
				InstallArguments: InstallArguments{
					Step: Step{
						Description: "Hello",
					},
					Manifests:  []string{manifestDirectory},
					KubeConfig: kubeConfig,
				},
			},
		},
	}

	defer os.Unsetenv(test.ExpectedCommandEnv)
	for _, installTest := range installTests {
		t.Run(installTest.expectedCommand, func(t *testing.T) {
			ctx := context.Background()
			os.Setenv(test.ExpectedCommandEnv, installTest.expectedCommand)

			action := InstallAction{Steps: []InstallStep{installTest.installStep}}
			b, _ := yaml.Marshal(action)

			h := NewTestMixin(t)
			h.In = bytes.NewReader(b)

			err := h.Install(ctx)

			require.NoError(t, err)
		})
	}
}
