package main

import (
	"archci/gitgen"
	"net/http"
	"os"
	"strings"

	"github.com/drone/drone-go/plugin/config"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// spec provides the plugin settings.
type spec struct {
	Bind             string `envconfig:"DRONE_BIND"`
	Debug            bool   `envconfig:"DRONE_DEBUG"`
	Secret           string `envconfig:"DRONE_SECRET"`
	GitRepoDirectory string `envconfig:"DRONE_GIT_REPO_DIRECTORY"`
}

func main() {
	spec := new(spec)
	err := envconfig.Process("", spec)
	if err != nil {
		logrus.Fatal(err)
	}

	if spec.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if spec.Secret == "" {
		logrus.Fatalln("missing secret key")
	}
	if spec.GitRepoDirectory == "" {
		// use a tmp directory if no directory is provided
		tempDir, err := os.MkdirTemp("", "drone-git-")
		if err != nil {
			logrus.Fatal(err)
		}
		spec.GitRepoDirectory = tempDir
	}
	if spec.Bind == "" {
		spec.Bind = ":3000"
	}

	// Tracked git repos are set in `DRONE_TRACKED_REPOS` env var, splited by `;`
	// Example: `DRONE_TRACKED_REPOS=github.com/owner/repo1;github.com/owner/repo2`
	// This is used to track the git repos that are cloned in the `DRONE_GIT_REPO_DIRECTORY`
	// and to avoid cloning the same repo multiple times.
	trackedRepos := os.Getenv("DRONE_TRACKED_REPOS")
	if trackedRepos == "" {
		logrus.Warn("no tracked repos provided")
	}

	conf := gitgen.Config{}
	conf.Repos = strings.Split(trackedRepos, ";")
	// purge empty strings
	for i := 0; i < len(conf.Repos); i++ {
		if conf.Repos[i] == "" {
			conf.Repos = append(conf.Repos[:i], conf.Repos[i+1:]...)
			i--
		}
	}

	// clone the tracked repos
	for _, repo := range conf.Repos {
		gitgen.UpdateRepo(repo, spec.GitRepoDirectory)
	}

	handler := config.Handler(
		gitgen.New(conf),
		spec.Secret,
		logrus.StandardLogger(),
	)

	logrus.Infof("server listening on address %s", spec.Bind)

	http.Handle("/", handler)
	logrus.Fatal(http.ListenAndServe(spec.Bind, nil))
}
