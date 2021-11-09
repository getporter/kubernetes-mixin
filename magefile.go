// +build mage

package main

import (
	// mage:import
	"get.porter.sh/porter/mage/releases"
)

// We are migrating to mage, but for now keep using make as the main build script interface.

// Publish the cross-compiled binaries.
func Publish(mixin string) {
	releases.PrepareMixinForPublish(mixin)
	releases.PublishMixin(mixin)
	releases.PublishMixinFeed(mixin)
}
