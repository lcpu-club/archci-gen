package gitgen

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

func GetRebuildList(apps []string) []string {
	// python genrebuild.py apps
	cmdline := "python genrebuild.py "
	for _, app := range apps {
		cmdline += app + " "
	}

	var stdout, stderr bytes.Buffer
	cmd := exec.Command("sh", "-c", cmdline)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		logrus.Errorf("failed to get rebuild list: %v", err)
		return nil
	}

	rebuild_list := strings.Split(strings.Split(stdout.String(), "\n")[0], " ")
	return rebuild_list
}
