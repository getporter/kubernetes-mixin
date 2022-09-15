package main

import (
	"get.porter.sh/mixin/kubernetes/pkg/kubernetes"
	"github.com/spf13/cobra"
)

func buildBuildCommand(m *kubernetes.Mixin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Generate Dockerfile contribution for invocation image",
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Build(cmd.Context())
		},
	}
	return cmd
}
