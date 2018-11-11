package slipway

import (
	composeTypes "github.com/docker/cli/cli/compose/types"
)

/* slipway.yaml
name: slipway
container:
  build:
    target: development
  volumes:
    - .:/go/src/github.com/timjones/slipway
*/

var SlipProject = &Project{
	Name: "slipway",
	Container: composeTypes.ServiceConfig{
		Build: composeTypes.BuildConfig{
			Target: "development",
		},
		Volumes: []composeTypes.ServiceVolumeConfig{
			{
				Source: ".",
				Target: "/go/src/github.com/timjones/slipway",
			},
		},
	},
}
