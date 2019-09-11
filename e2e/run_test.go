//
// This is the e2e package to run tests for Ignite
// Currently, we do local e2e tests only
// we have to wait until the CI setup to allow Ignite to run with sudo and in a KVM environment.
//
// How to run tests:
// sudo IGNITE_E2E_HOME=$PWD $(which go) test ./e2e/. -count 1
//

package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"
	"time"

	"gotest.tools/assert"
)

var (
	e2eHome   = os.Getenv("IGNITE_E2E_HOME")
	igniteBin = path.Join(e2eHome, "bin/ignite")
)

// stdCmd builds an *exec.Cmd hooked up to Stdout/Stderr by default
func stdCmd(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

// runWithRuntimeAndNetworkPlugin is a helper for running a vm then forcing removal
// vmName should be unique for each test
func runWithRuntimeAndNetworkPlugin(t *testing.T, vmName, runtime, networkPlugin string) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	runCmd := stdCmd(
		igniteBin,
		"--runtime="+runtime,
		"--network-plugin="+networkPlugin,
		"run", "--name="+vmName,
		"weaveworks/ignite-ubuntu",
	)
	runErr := runCmd.Run()

	defer func() {
		rmvCmd := stdCmd(
			igniteBin,
			"--runtime="+runtime,
			"--network-plugin="+networkPlugin,
			"rm", "-f", vmName,
		)
		rmvErr := rmvCmd.Run()
		assert.Check(t, rmvErr, fmt.Sprintf("vm removal should not fail: %q", rmvCmd.Args))
	}()

	assert.Check(t, runErr, fmt.Sprintf("%q should not fail to run", runCmd.Args))
}

func TestIgniteRunWithDockerAndDockerBridge(t *testing.T) {
	runWithRuntimeAndNetworkPlugin(
		t,
		"e2e_test_ignite_run_docker_and_docker_bridge",
		"docker",
		"docker-bridge",
	)
}

func TestIgniteRunWithDockerAndCNI(t *testing.T) {
	runWithRuntimeAndNetworkPlugin(
		t,
		"e2e_test_ignite_run_docker_and_cni",
		"docker",
		"cni",
	)
}

func TestIgniteRunWithContainerdAndCNI(t *testing.T) {
	runWithRuntimeAndNetworkPlugin(
		t,
		"e2e_test_ignite_run_containerd_and_cni",
		"containerd",
		"cni",
	)
}

// runCurl is a helper for testing network connectivity
// vmName should be unique for each test
func runCurl(t *testing.T, vmName, runtime, networkPlugin string, sleepDuration time.Duration) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	runCmd := stdCmd(
		igniteBin,
		"--runtime="+runtime,
		"--network-plugin="+networkPlugin,
		"run", "--name="+vmName,
		"--dns=8.8.8.8", "--dns=8.8.4.4", // override name servers
		"weaveworks/ignite-ubuntu",
		"--ssh",
	)
	runErr := runCmd.Run()

	defer func() {
		rmvCmd := stdCmd(
			igniteBin,
			"--runtime="+runtime,
			"--network-plugin="+networkPlugin,
			"rm", "-f", vmName,
		)
		rmvErr := rmvCmd.Run()
		assert.Check(t, rmvErr, fmt.Sprintf("vm removal should not fail: %q", rmvCmd.Args))
	}()

	assert.Check(t, runErr, fmt.Sprintf("%q should not fail to run", runCmd.Args))
	if runErr != nil {
		return
	}

	time.Sleep(sleepDuration)
	curlCmd := stdCmd(
		igniteBin,
		"--runtime="+runtime,
		"--network-plugin="+networkPlugin,
		"exec", vmName,
		"curl", "google.com",
	)
	curlErr := curlCmd.Run()
	assert.Check(t, curlErr, fmt.Sprintf("curl should not fail: %q", curlCmd.Args))
}

func TestCurlWithDockerAndDockerBridge(t *testing.T) {
	runCurl(
		t,
		"e2e_test_curl_docker_and_docker_bridge",
		"docker",
		"docker-bridge",
		0,
	)
}

func TestCurlWithDockerAndCNI(t *testing.T) {
	runCurl(
		t,
		"e2e_test_curl_docker_and_cni",
		"docker",
		"cni",
		0,
	)
}

func TestCurlWithContainerdAndCNI(t *testing.T) {
	runCurl(
		t,
		"e2e_test_curl_containerd_and_cni",
		"containerd",
		"cni",
		0,
	)
}

func TestCurlWithDockerAndCNISleep2(t *testing.T) {
	runCurl(
		t,
		"e2e_test_curl_docker_and_cni_sleep2",
		"docker",
		"cni",
		2 * time.Second, // TODO: why is this necessary? Can we work to eliminate this?
	)
}

func TestCurlWithContainerdAndCNISleep2(t *testing.T) {
	runCurl(
		t,
		"e2e_test_curl_containerd_and_cni_sleep2",
		"containerd",
		"cni",
		2 * time.Second,
	)
}
