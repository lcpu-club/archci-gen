package gitgen

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/gofrs/flock"
	"github.com/sirupsen/logrus"
)

func UpdateRepo(repo, dir string) error {
	// get flock on the directory
	fileLock := flock.New(dir + ".lock")
	locked, err := fileLock.TryLock()

	if err != nil {
		return fmt.Errorf("failed to get lock on directory: %v", err)
	}

	if !locked {
		// loop to sleep and try again
		for {
			locked, err = fileLock.TryLock()

			if err != nil {
				return fmt.Errorf("failed to get lock on directory: %v", err)
			}

			if locked {
				break
			}
			time.Sleep(10 * time.Second)
		}
	}

	defer fileLock.Unlock()

	// if the folder is empty, git clone --mirror repo dir
	// else, git fetch --all
	if _, err := exec.Command("ls", dir).Output(); err != nil {
		cmd := exec.Command("git", "clone", "--mirror", repo, dir)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to clone repository: %v", err)
		}
		if stderr.String() != "" {
			logrus.Infof("stderr: %s", stderr.String())
		}
		logrus.Infof("stdout: %s", stdout.String())
	} else {
		cmd := exec.Command("git", "fetch", "--all")
		cmd.Dir = dir
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to fetch repository: %v", err)
		}
		if stderr.String() != "" {
			logrus.Infof("stderr: %s", stderr.String())
		}
		logrus.Infof("stdout: %s", stdout.String())
	}
	return nil
}

// git clone --reference reference --branch branch --single-branch repo dir
func CloneBranch(repo, reference, branch, dir string) error {
	// get flock on the directory, if failed, just return
	fileLock := flock.New(dir + ".lock")
	locked, err := fileLock.TryLock()

	if err != nil {
		return fmt.Errorf("failed to get lock on directory: %v", err)
	}

	if !locked {
		return fmt.Errorf("failed to get lock on directory")
	}

	defer fileLock.Unlock()

	cmd := exec.Command("git", "clone", "--reference", reference, "--branch", branch, repo, dir)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to clone branch: %v", err)
	}
	if stderr.String() != "" {
		logrus.Infof("stderr: %s", stderr.String())
	}
	logrus.Infof("stdout: %s", stdout.String())
	return nil
}

func GetBranchChangedFile(repo, dir, branch string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-status", "origin/testing")
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to get changed files: %v", err)
	}
	if stderr.String() != "" {
		logrus.Infof("stderr: %s", stderr.String())
	}

	// Parse, split by space
	// Check https://git-scm.com/docs/git-status for more information
	allChangedFiles := strings.Split(stdout.String(), "\n")
	var changedFiles []string
	for _, file := range allChangedFiles {
		if file != "" {
			status, filename := strings.Split(file, "\t")[0], strings.Split(file, "\t")[1]
			if status == "M" || status == "A" {
				changedFiles = append(changedFiles, filename)
			}
		}
	}
	return changedFiles, nil
}
