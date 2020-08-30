package kubernetes

import (
	"get.porter.sh/porter/pkg/mixin"
	"get.porter.sh/porter/pkg/porter/version"
	"github.com/deislabs/porter-kubernetes/pkg"
)

func (m *Mixin) PrintVersion(opts version.Options) error {
	metadata := mixin.Metadata{
		Name: "kubernetes",
		VersionInfo: mixin.VersionInfo{
			Version: pkg.Version,
			Commit:  pkg.Commit,
			Author:  "deislabs",
		},
	}
	return version.PrintVersion(m.Context, opts, metadata)
}
