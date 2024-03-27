package gitgen

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func GetChangedFromDiff(diff []string) ([]string, error) {
	// Get the list of changed files from the diff
	changed := make([]string, 0)
	for _, line := range diff {
		// split by '/', get 0
		changed = append(changed, strings.Split(line, "/")[0])
	}
	return changed, nil
}

type BuildJob struct {
	Name       string   `yaml:"name"`
	Image      string   `yaml:"image"`
	Commands   []string `yaml:"commands"`
	Depends_on []string `yaml:"depends_on"`
}

type BuildEnv struct {
	Environment map[string]BuildSecret `yaml:"environment"`
}

type BuildSecret struct {
	FromSecret string `yaml:"from_secret"`
}

type BuildPipeline struct {
	Steps []BuildJob `yaml:"steps"`
}

type BuildFile struct {
	Kind        string     `yaml:"kind"`
	Type        string     `yaml:"type"`
	Name        string     `yaml:"name"`
	Environment BuildEnv   `yaml:"environment"`
	Steps       []BuildJob `yaml:"steps"`
}

func GenereateYaml(pipeline BuildPipeline) (string, error) {
	buildenv := BuildEnv{
		Environment: make(map[string]BuildSecret),
	}
	env_file, err := os.ReadFile(".drone-env")
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(env_file), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		buildenv.Environment[line] = BuildSecret{
			FromSecret: line,
		}
	}

	buildfile := BuildFile{
		Kind:        "pipeline",
		Type:        "docker",
		Name:        "default",
		Environment: buildenv,
		Steps:       pipeline.Steps,
	}

	data, err := yaml.Marshal(&buildfile)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
