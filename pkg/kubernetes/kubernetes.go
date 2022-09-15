package kubernetes

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"get.porter.sh/porter/pkg/runtime"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonschema"
)

const (
	defaultKubernetesClientVersion string = "v1.15.5"
)

type Mixin struct {
	runtime.RuntimeConfig
	KubernetesClientVersion string
}

func New() *Mixin {
	return &Mixin{
		RuntimeConfig:           runtime.NewConfig(),
		KubernetesClientVersion: defaultKubernetesClientVersion,
	}
}

func (m *Mixin) getCommandFile(commandFile string, w io.Writer) ([]byte, error) {
	if commandFile == "" {
		reader := bufio.NewReader(m.In)
		return ioutil.ReadAll(reader)
	}
	return ioutil.ReadFile(commandFile)
}

func (m *Mixin) getPayloadData() ([]byte, error) {
	reader := bufio.NewReader(m.In)
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		errors.Wrap(err, "could not read payload from STDIN")
	}
	return data, nil
}

func (m *Mixin) ValidatePayload(b []byte) error {
	// Load the step as a go dump
	s := make(map[string]interface{})
	err := yaml.Unmarshal(b, &s)
	if err != nil {
		return errors.Wrap(err, "could not marshal payload as yaml")
	}
	manifestLoader := gojsonschema.NewGoLoader(s)

	// Load the step schema
	schemaLoader := gojsonschema.NewStringLoader(schema)

	validator, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return errors.Wrap(err, "unable to compile the mixin step schema")
	}

	// Validate the manifest against the schema
	result, err := validator.Validate(manifestLoader)
	if err != nil {
		return errors.Wrap(err, "unable to validate the mixin step schema")
	}
	if !result.Valid() {
		errs := make([]string, 0, len(result.Errors()))
		for _, err := range result.Errors() {
			errs = append(errs, err.String())
		}
		return errors.New(strings.Join(errs, "\n\t* "))
	}

	return nil
}

func (m *Mixin) getOutput(ctx context.Context, resourceType, resourceName, namespace, jsonPath string) ([]byte, error) {
	args := []string{"get", resourceType, resourceName}
	args = append(args, fmt.Sprintf("-o=jsonpath=%s", jsonPath))
	if namespace != "" {
		args = append(args, fmt.Sprintf("--namespace=%s", namespace))
	}
	cmd := m.NewCommand(ctx, "kubectl", args...)
	cmd.Stderr = m.Err

	prettyCmd := fmt.Sprintf("%s%s", cmd.Dir, strings.Join(cmd.Args, " "))
	if m.DebugMode {
		fmt.Fprintln(m.Err, prettyCmd)
	}
	out, err := cmd.Output()

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("couldn't run command %s", prettyCmd))
	}
	return out, nil
}

func (m *Mixin) handleOutputs(ctx context.Context, outputs []KubernetesOutput) error {
	//Now get the outputs
	for _, output := range outputs {
		bytes, err := m.getOutput(ctx,
			output.ResourceType,
			output.ResourceName,
			output.Namespace,
			output.JSONPath,
		)
		if err != nil {
			return err
		}
		err = m.Context.WriteMixinOutputToFile(output.Name, bytes)
		if err != nil {
			return err
		}
	}
	return nil
}
