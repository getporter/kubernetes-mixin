package main

import (
	"github.com/deislabs/porter-kubernetes/pkg/kubernetes"
	"github.com/spf13/cobra"
)

func buildUninstallCommand(m *kubernetes.Mixin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Execute the uninstall functionality of this mixin",
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Execute()
		},
	}
	return cmd
}
