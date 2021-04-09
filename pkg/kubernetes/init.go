package kubernetes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/mod/semver"
)

type KubectlVersion struct {
	ClientVersion struct {
		Major        string    `json:"major"`
		Minor        string    `json:"minor"`
		GitVersion   string    `json:"gitVersion"`
		GitCommit    string    `json:"gitCommit"`
		GitTreeState string    `json:"gitTreeState"`
		BuildDate    time.Time `json:"buildDate"`
		GoVersion    string    `json:"goVersion"`
		Compiler     string    `json:"compiler"`
		Platform     string    `json:"platform"`
	} `json:"clientVersion"`
	ServerVersion struct {
		Major        string    `json:"major"`
		Minor        string    `json:"minor"`
		GitVersion   string    `json:"gitVersion"`
		GitCommit    string    `json:"gitCommit"`
		GitTreeState string    `json:"gitTreeState"`
		BuildDate    time.Time `json:"buildDate"`
		GoVersion    string    `json:"goVersion"`
		Compiler     string    `json:"compiler"`
		Platform     string    `json:"platform"`
	} `json:"serverVersion"`
}

func (m *Mixin) init() error {

	serverVersion, err := getServerVersion(m)

	if err != nil {
		return err
	}

	if semver.Compare(m.KubernetesClientVersion, serverVersion) == -1 && serverVersion != "" {

		// Here we are triggering the download
		fmt.Fprintf(m.Out, "Kubectl server version (%s) does not match client version (%s); downloading a compatible client.\n",
			serverVersion, m.KubernetesClientVersion)
		// try to install the new client
		err := installClient(m, serverVersion)
		if err != nil {
			return errors.Wrap(err, "unable to install a compatible kubectl client")
		}
	}

	return err
}

func installClient(m *Mixin, version string) error {

	url := fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/amd64/kubectl", version)

	// Fetch kubectl from url
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return errors.Wrap(err, "failed to construct GET request for fetching kubectl client binary")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "failed to download kubectl client binary via url: %s", url)
	}
	defer res.Body.Close()

	// Create a temp dir
	tmpDir, err := m.FileSystem.TempDir("", "tmp")
	if err != nil {
		return errors.Wrap(err, "unable to create a temporary directory for downloading the kubectl client binary")
	}
	defer os.RemoveAll(tmpDir)

	// Create the local binary
	kubectlBinPath, err := m.FileSystem.Create(filepath.Join(tmpDir, "kubectlBin"))
	if err != nil {
		return errors.Wrap(err, "unable to create a local file for the kubectl client binary")
	}

	// Copy response body to local binary
	_, err = io.Copy(kubectlBinPath, res.Body)
	if err != nil {
		return errors.Wrap(err, "unable to copy the kubectl client binary to the local binary file")
	}

	// Move the kubectl binary into the appropriate location
	binPath := "/usr/local/bin/kubectl"
	err = m.FileSystem.Rename(fmt.Sprintf("%s", kubectlBinPath.Name()), binPath)
	if err != nil {
		return errors.Wrapf(err, "unable to install the kubectl client binary to %q", binPath)
	}
	return nil
}

func getServerVersion(m *Mixin) (string, error) {

	var stderr bytes.Buffer
	currentKubectl := KubectlVersion{}

	cmd := exec.Command("kubectl", "version", "-o", "json")
	cmd.Stderr = &stderr
	// Execute the command and the version output
	outputBytes, err := cmd.Output()
	if err != nil {
		return "", errors.Wrapf(err, "unable to determine kubernetes server version: %s", stderr.String())
	}
	// Rebuild version json object
	json.Unmarshal(outputBytes, &currentKubectl)

	version := currentKubectl.ServerVersion.GitVersion

	return version, nil
}
