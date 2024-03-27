package gitgen

import (
	"fmt"
	"testing"
)

func TestGenereateYaml(t *testing.T) {
	pipeline := BuildPipeline{
		Steps: []BuildJob{
			{
				Name:     "build",
				Image:    "golang:1.16",
				Commands: []string{"go build", "go test"},
			},
			{
				Name:       "deploy",
				Image:      "alpine:latest",
				Commands:   []string{"echo 'Deploying...'", "echo 'Deployed'"},
				Depends_on: []string{"build"},
			},
		},
	}
	yaml, err := GenereateYaml(pipeline)
	if err != nil {
		t.Errorf("Failed to generate yaml: %v", err)
	}
	if yaml == "" {
		t.Errorf("Generated yaml is empty")
	}
	fmt.Print(yaml)
}
