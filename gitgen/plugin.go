package gitgen

import (
	"context"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/config"
	"github.com/go-git/go-git/plumbing/transport"
)

// New returns a new config plugin.
func New(config Config) config.Plugin {
	return newArchPlugin(config)
}

func newArchPlugin(config Config) *archPlugin {
	plugin := &archPlugin{}
	return plugin
}

type archPlugin struct {
	auth  transport.AuthMethod
	index map[string]string
}

func (r *archPlugin) Find(ctx context.Context, req *config.Request) (*drone.Config, error) {
	// load
	return nil, nil
}
