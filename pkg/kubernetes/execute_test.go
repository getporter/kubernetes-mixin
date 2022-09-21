package kubernetes

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"get.porter.sh/porter/pkg/test"
	"github.com/stretchr/testify/require"

	"gopkg.in/yaml.v2"
)

type ExecuteTest struct {
	expectedCommand string
	executeStep     ExecuteStep
}

func TestMixin_ExecuteStep(t *testing.T) {

	manifestDirectory := "/cnab/app/manifests"

	upgradeCmd := "kubectl apply -f"

	dontWait := false

	recordIt := true
	validateIt := false

	namespace := "meditations"

	selector := "app=nginx"

	k8scontext := "context"

	forceIt := true
	withGrace := 1

	overwriteIt := false
	pruneIt := true

	timeout := 1

	upgradeTests := []ExecuteTest{
		// These tests are largely the same as the install, just testing that the embedded
		// install gets handled correctly
		{
			expectedCommand: fmt.Sprintf("%s %s --wait", upgradeCmd, manifestDirectory),
			executeStep: ExecuteStep{
				ExecuteInstruction: ExecuteInstruction{
					InstallArguments: InstallArguments{
						Step: Step{
							Description: "Hello",
						},
						Manifests: []string{manifestDirectory},
					},
				},
			},
		},
		{
			expectedCommand: fmt.Sprintf("%s %s --wait", upgradeCmd, manifestDirectory),
			executeStep: ExecuteStep{
				ExecuteInstruction: ExecuteInstruction{
					InstallArguments: InstallArguments{
						Step: Step{
							Description: "Hello",
						},
						Manifests: []string{manifestDirectory},
					},
				},
			},
		},
		{
			expectedCommand: fmt.Sprintf("%s %s", upgradeCmd, manifestDirectory),
			executeStep: ExecuteStep{
				ExecuteInstruction: ExecuteInstruction{
					InstallArguments: InstallArguments{
						Step: Step{
							Description: "Hello",
						},
						Manifests: []string{manifestDirectory},
						Wait:      &dontWait,
					},
				},
			},
		},
		{
			expectedCommand: fmt.Sprintf("%s %s -n %s", upgradeCmd, manifestDirectory, namespace),
			executeStep: ExecuteStep{
				ExecuteInstruction: ExecuteInstruction{
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
		},
		{
			expectedCommand: fmt.Sprintf("%s %s -n %s --validate=false", upgradeCmd, manifestDirectory, namespace),
			executeStep: ExecuteStep{
				ExecuteInstruction: ExecuteInstruction{
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
		},
		{
			expectedCommand: fmt.Sprintf("%s %s -n %s --record=true", upgradeCmd, manifestDirectory, namespace),
			executeStep: ExecuteStep{
				ExecuteInstruction: ExecuteInstruction{
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
		},
		{
			expectedCommand: fmt.Sprintf("%s %s --selector=%s --wait", upgradeCmd, manifestDirectory, selector),
			executeStep: ExecuteStep{
				ExecuteInstruction: ExecuteInstruction{
					InstallArguments: InstallArguments{
						Step: Step{
							Description: "Hello",
						},
						Manifests: []string{manifestDirectory},
						Selector:  selector,
					},
				},
			},
		},
		{
			expectedCommand: fmt.Sprintf("%s %s --context=%s --wait", upgradeCmd, manifestDirectory, k8scontext),
			executeStep: ExecuteStep{
				ExecuteInstruction: ExecuteInstruction{
					InstallArguments: InstallArguments{
						Step: Step{
							Description: "Hello",
						},
						Manifests: []string{manifestDirectory},
						Context:   k8scontext,
					},
				},
			},
		},

		// These tests exercise the upgrade options
		{
			expectedCommand: fmt.Sprintf("%s %s --wait --force --grace-period=0", upgradeCmd, manifestDirectory),
			executeStep: ExecuteStep{
				ExecuteInstruction: ExecuteInstruction{
					Force: &forceIt,
					InstallArguments: InstallArguments{
						Step: Step{
							Description: "Hello",
						},
						Manifests: []string{manifestDirectory},
					},
				},
			},
		},
		{
			expectedCommand: fmt.Sprintf("%s %s --wait --grace-period=%d", upgradeCmd, manifestDirectory, withGrace),
			executeStep: ExecuteStep{
				ExecuteInstruction: ExecuteInstruction{
					GracePeriod: &withGrace,
					InstallArguments: InstallArguments{
						Step: Step{
							Description: "Hello",
						},
						Manifests: []string{manifestDirectory},
					},
				},
			},
		},
		{
			expectedCommand: fmt.Sprintf("%s %s --wait --overwrite=false", upgradeCmd, manifestDirectory),
			executeStep: ExecuteStep{
				ExecuteInstruction: ExecuteInstruction{
					Overwrite: &overwriteIt,
					InstallArguments: InstallArguments{
						Step: Step{
							Description: "upgrade",
						},
						Manifests: []string{manifestDirectory},
					},
				},
			},
		},
		{
			expectedCommand: fmt.Sprintf("%s %s --wait --prune=true", upgradeCmd, manifestDirectory),
			executeStep: ExecuteStep{
				ExecuteInstruction: ExecuteInstruction{
					Prune: &pruneIt,
					InstallArguments: InstallArguments{
						Step: Step{
							Description: "upgrade",
						},
						Manifests: []string{manifestDirectory},
					},
				},
			},
		},
		{
			expectedCommand: fmt.Sprintf("%s %s --wait --timeout=%ds", upgradeCmd, manifestDirectory, timeout),
			executeStep: ExecuteStep{
				ExecuteInstruction: ExecuteInstruction{
					Timeout: &timeout,
					InstallArguments: InstallArguments{
						Step: Step{
							Description: "upgrade",
						},
						Manifests: []string{manifestDirectory},
					},
				},
			},
		},
	}

	for _, upgradeTest := range upgradeTests {
		upgradeTest := upgradeTest
		t.Run(upgradeTest.expectedCommand, func(t *testing.T) {
			ctx := context.Background()

			action := ExecuteAction{Steps: []ExecuteStep{upgradeTest.executeStep}}
			b, _ := yaml.Marshal(action)

			h := NewTestMixin(t)
			h.Setenv(test.ExpectedCommandEnv, upgradeTest.expectedCommand)
			h.In = bytes.NewReader(b)

			err := h.Execute(ctx)

			require.NoError(t, err)
		})
	}
}
