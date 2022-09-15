package kubernetes

import (
	"testing"

	"get.porter.sh/porter/pkg/portercontext"
)

type TestMixin struct {
	*Mixin
	TestContext *portercontext.TestContext
}

func NewTestMixin(t *testing.T) *TestMixin {
	c := portercontext.NewTestContext(t)
	m := New()
	m.Context = c.Context
	return &TestMixin{
		Mixin:       m,
		TestContext: c,
	}
}
