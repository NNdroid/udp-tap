package common

import (
	"os/exec"
	"udp-tap/pkg/log"
)

func ExecCmd(name string, arg ...string) error {
	log.Logger().Debugf("exec cmd: %s %v", name, arg)
	cmd := exec.Command(name, arg...)
	return cmd.Run()
}
