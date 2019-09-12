package run

import (
	"fmt"
	"net"
	"time"

	"github.com/weaveworks/ignite/pkg/operations"
	"github.com/weaveworks/ignite/pkg/preflight/checkers"
	"github.com/weaveworks/ignite/pkg/util"
	"k8s.io/apimachinery/pkg/util/sets"
)

type StartFlags struct {
	Interactive            bool
	Debug                  bool
	IgnoredPreflightErrors []string
}

type startOptions struct {
	*StartFlags
	*attachOptions
}

func (sf *StartFlags) NewStartOptions(vmMatch string) (*startOptions, error) {
	ao, err := NewAttachOptions(vmMatch)
	if err != nil {
		return nil, err
	}

	// Disable running check as it takes a while for ignite-spawn to update the state
	ao.checkRunning = false

	return &startOptions{sf, ao}, nil
}

func Start(so *startOptions) error {
	// Check if the given VM is already running
	if so.vm.Running() {
		return fmt.Errorf("VM %q is already running", so.vm.GetUID())
	}

	ignoredPreflightErrors := sets.NewString(util.ToLower(so.StartFlags.IgnoredPreflightErrors)...)
	if err := checkers.StartCmdChecks(so.vm, ignoredPreflightErrors); err != nil {
		return err
	}

	if err := operations.StartVM(so.vm, so.Debug); err != nil {
		return err
	}

	// When --ssh is enabled, then wait until SSH service started on port 22
	ssh := so.vm.Spec.SSH
	if ssh != nil && ssh.Generate {
		if len(so.vm.Status.IPAddresses) > 0 {
			addr := so.vm.Status.IPAddresses[0].String() + ":22"
			var err error
			for i := 0; i < 500; i++ {
				conn, dialErr := net.DialTimeout("tcp", addr, 100*time.Millisecond)
				if conn != nil {
					defer conn.Close()
					err = nil
					break
				}
				err = dialErr
				time.Sleep(100 * time.Millisecond)
			}
			if err != nil {
				if err, ok := err.(*net.OpError); ok && err.Timeout() {
					return fmt.Errorf("Tried connecting to SSH but timed out %s", err)
				}
			}
		}
	}

	// If starting interactively, attach after starting
	if so.Interactive {
		return Attach(so.attachOptions)
	}
	return nil
}
