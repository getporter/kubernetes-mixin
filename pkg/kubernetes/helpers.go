package kubernetes

import (
	"testing"

	"get.porter.sh/porter/pkg/context"
)

type TestMixin struct {
	*Mixin
	TestContext *context.TestContext
}

const MockkubectlClientVersion string = "v1.15.5"

func NewTestMixin(t *testing.T) *TestMixin {
	c := context.NewTestContext(t)
	m := New()
	m.Context = c.Context
	return &TestMixin{
		Mixin:       m,
		TestContext: c,
	}
}
