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

type MockKubectlDownloader struct {
	GetServerVersion func(m *Mixin) (string, error)
	InstallClient    func(m *Mixin, version string) error
}

func (kd MockKubectlDownloader) getServerVersion(m *Mixin) (string, error) {
	return kd.GetServerVersion(m)
}

func (kd MockKubectlDownloader) installClient(m *Mixin, version string) error {
	return kd.InstallClient(m, version)
}

func NewMockKubectlDownloader() MockKubectlDownloader {

	return MockKubectlDownloader{
		GetServerVersion: func(m *Mixin) (string, error) {
			return MockkubectlClientVersion, nil
		},
		InstallClient: func(m *Mixin, version string) error {
			return nil
		},
	}
}

func NewTestMixin(t *testing.T) *TestMixin {
	c := context.NewTestContext(t)
	m := New()
	m.Context = c.Context
	m.KubectlDownloader = NewMockKubectlDownloader()
	return &TestMixin{
		Mixin:       m,
		TestContext: c,
	}
}
