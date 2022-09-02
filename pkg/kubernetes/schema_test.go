package kubernetes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMixin_PrintSchema(t *testing.T) {
	m := NewTestMixin(t)

	err := m.PrintSchema()
	require.NoError(t, err)

	gotSchema := m.TestContext.GetOutput()

	wantSchema, err := ioutil.ReadFile("./schema/schema.json")
	require.NoError(t, err)

	assert.Equal(t, string(wantSchema), gotSchema)
}

func TestMixin_CheckSchema(t *testing.T) {
	// Long term it would be great to have a helper function in Porter that a mixin can use to check that it meets certain interfaces
	// check that certain characteristics of the schema that Porter expects are present
	// Once we have a mixin library, that would be a good place to package up this type of helper function
	var schemaMap map[string]interface{}
	err := json.Unmarshal([]byte(schema), &schemaMap)
	require.NoError(t, err, "could not unmarshal the schema into a map")

	t.Run("mixin configuration", func(t *testing.T) {
		// Check that mixin config is defined, and has all the supported fields
		configSchema, err := jsonpath.Get("$.definitions.config", schemaMap)
		require.NoError(t, err, "could not find the mixin config schema declaration")
		_, err = jsonpath.Get("$.properties.kubernetes.properties.clientVersion", configSchema)
		require.NoError(t, err, "client version was not included in the mixin config schema")
	})

	// Check that schema are defined for each action
	actions := []string{"install", "upgrade", "invoke", "uninstall"}
	for _, action := range actions {
		t.Run("supports "+action, func(t *testing.T) {
			actionPath := fmt.Sprintf("$.definitions.%sStep", action)
			_, err := jsonpath.Get(actionPath, schemaMap)
			require.NoErrorf(t, err, "could not find the %sStep declaration", action)
		})
	}

	// Check that the invoke action is registered
	additionalSchema, err := jsonpath.Get("$.additionalProperties.items", schemaMap)
	require.NoError(t, err, "the invoke action was not registered in the schema")
	require.Contains(t, additionalSchema, "$ref")
	invokeRef := additionalSchema.(map[string]interface{})["$ref"]
	require.Equal(t, "#/definitions/invokeStep", invokeRef, "the invoke action was not registered correctly")
}

func TestMixin_ValidatePayload(t *testing.T) {
	testcases := []struct {
		name  string
		step  string
		pass  bool
		error string
	}{
		{"install", "testdata/install-input.yaml", true, ""},
		{"install-with-kubeconfig", "testdata/install-input-with-kubeconfig.yaml", true, ""},
		{"upgrade", "testdata/upgrade-input.yaml", true, ""},
		{"invoke", "testdata/invoke-input.yaml", true, ""},
		{"uninstall", "testdata/uninstall-input.yaml", true, ""},
		{"install-bad-wait-flag", "testdata/install-input-bad-wait-flag.yaml", false, "install.0.kubernetes.wait: Invalid type. Expected: boolean, given: string"},
		{"install-no-manifests", "testdata/install-input-no-manifests.yaml", false, "install.0.kubernetes: manifests is required"},
		{"install-bad-outputs", "testdata/install-input-bad-outputs.yaml", false, "install.0.kubernetes.outputs.0: resourceType is required\n\t* install.0.kubernetes.outputs.0: resourceName is required\n\t* install.0.kubernetes.outputs.0: jsonPath is required"},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewTestMixin(t)
			b, err := ioutil.ReadFile(tc.step)
			require.NoError(t, err)

			err = m.ValidatePayload(b)
			if tc.pass {
				require.NoError(t, err)
			} else {
				require.Contains(t, err.Error(), tc.error)
			}
		})
	}
}

// TODO: Currently can't test this with pkg/test because it runs multiple commands
// func TestOuputSchema(t *testing.T) {
// 	manifestDirectory := "/cnab/app/manifests"

// 	installCmd := "kubectl apply -f"

// 	installTests := []InstallTest{
// 		{
// 			expectedCommand: fmt.Sprintf("%s %s --wait", installCmd, manifestDirectory),
// 			installStep: InstallStep{

// 				InstallArguments: InstallArguments{
// 					Step: Step{
// 						Description: "Hello",
// 						Outputs: []KubernetesOutput{
// 							KubernetesOutput{
// 								Name:         "test",
// 								Namespace:    "Default",
// 								ResourceType: "service",
// 								ResourceName: "aservice",
// 								JSONPath:     "a.path.",
// 							},
// 						},
// 					},
// 					Manifests: []string{manifestDirectory},
// 				},
// 			},
// 		},
// 		{
// 			expectedCommand: fmt.Sprintf("%s %s --wait", installCmd, manifestDirectory),
// 			installStep: InstallStep{

// 				InstallArguments: InstallArguments{
// 					Step: Step{
// 						Description: "Hello",
// 						Outputs: []KubernetesOutput{
// 							KubernetesOutput{
// 								Name:         "test",
// 								ResourceType: "service",
// 								ResourceName: "aservice",
// 								JSONPath:     "a.path.",
// 							},
// 						},
// 					},
// 					Manifests: []string{manifestDirectory},
// 				},
// 			},
// 		},
// 	}

// 	defer os.Unsetenv(test.ExpectedCommandEnv)
// 	for _, installTest := range installTests {
// 		t.Run(installTest.expectedCommand, func(t *testing.T) {
// 			os.Setenv(test.ExpectedCommandEnv, installTest.expectedCommand)

// 			action := InstallAction{Steps: []InstallStep{installTest.installStep}}
// 			b, _ := yaml.Marshal(action)

// 			h := NewTestMixin(t)
// 			h.In = bytes.NewReader(b)

// 			err := h.Install()

// 			require.NoError(t, err)
// 		})
// 	}
// }
