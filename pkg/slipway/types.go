package slipway

import (
	composeTypes "github.com/docker/cli/cli/compose/types"
)

type Project struct {
	Name      string `yaml:",omitempty"`
	Container composeTypes.ServiceConfig
}
