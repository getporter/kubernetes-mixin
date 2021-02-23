// +build mage

package main

import (
	"context"

	// mage:import
	"get.porter.sh/porter/mage/releases"
)

// We are migrating to mage, but for now keep using make as the main build script interface.

// Publish the cross-compiled binaries.
func Publish(ctx context.Context, mixin string, version string, permalink string) {
	releases.PrepareMixinForPublish(mixin, version, permalink)
	releases.PublishMixin(mixin, version, permalink)
	releases.PublishMixinFeed(ctx)
}
