package gitgen

import (
	"os"
	"testing"
)

func TestCloneRepo(t *testing.T) {
	repo := "https://github.com/lcpu-club/boardcrafts"
	dir := "/tmp/boardcrafts"

	// Clean up the directory before running the test
	err := os.RemoveAll(dir)
	if err != nil {
		t.Fatalf("Failed to clean up directory: %v", err)
	}

	err = UpdateRepo(repo, dir)
	if err != nil {
		t.Errorf("Failed to clone repository: %v", err)
	}

	// Verify that the repository was cloned successfully
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		t.Errorf("Repository directory does not exist: %v", dir)
	}
}

func TestCloneBranch(t *testing.T) {
	repo := "https://github.com/lcpu-club/boardcrafts"
	reference := "/tmp/boardcrafts"
	dir := "/tmp/boardcrafts-pvkmc"
	branch := "pvkmc"

	// Clean up the directory before running the test
	err := os.RemoveAll(reference)
	if err != nil {
		t.Fatalf("Failed to clean up directory: %v", err)
	}
	err = UpdateRepo(repo, reference)
	if err != nil {
		t.Errorf("Failed to clone repository: %v", err)
	}

	err = os.RemoveAll(dir)
	if err != nil {
		t.Fatalf("Failed to clean up directory: %v", err)
	}
	err = CloneBranch(repo, reference, branch, dir)
	if err != nil {
		t.Errorf("Failed to clone branch: %v", err)
	}

	// Verify that the branch was cloned successfully
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		t.Errorf("Branch directory does not exist: %v", dir)
	}
}

func TestBranchDiff(t *testing.T) {
	repo := "https://github.com/lcpu-club/boardcrafts"
	reference := "/tmp/boardcrafts"
	dir := "/tmp/boardcrafts-test"
	branch := "test"

	// Clean up the directory before running the test
	err := os.RemoveAll(reference)
	if err != nil {
		t.Fatalf("Failed to clean up directory: %v", err)
	}
	err = UpdateRepo(repo, reference)
	if err != nil {
		t.Errorf("Failed to clone repository: %v", err)
	}

	err = os.RemoveAll(dir)
	if err != nil {
		t.Fatalf("Failed to clean up directory: %v", err)
	}
	err = CloneBranch(repo, reference, branch, dir)
	if err != nil {
		t.Errorf("Failed to clone branch: %v", err)
	}

	// Verify that the branch was cloned successfully
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		t.Errorf("Branch directory does not exist: %v", dir)
	}

	// Test the branch diff
	files, err := GetBranchChangedFile(repo, dir, branch)
	if err != nil {
		t.Errorf("Failed to get branch diff: %v", err)
	}

	// Verify that the branch diff is not empty
	if len(files) == 0 {
		t.Errorf("Branch diff is empty")
	}
}
