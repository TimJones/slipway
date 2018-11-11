package slipway

import (
	"io"
	"os"

	"github.com/docker/cli/cli/command/image/build"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/pkg/streamformatter"
	"golang.org/x/net/context"
)

func (*Project) Build() error {
	projectBuild := SlipProject.Container.Build

	contextDir, relDockerfile, err := build.GetContextFromLocalDir(projectBuild.Context, projectBuild.Dockerfile)
	if err != nil {
		return err
	}

	excludes, err := build.ReadDockerignore(contextDir)
	if err != nil {
		return err
	}

	relDockerfile = archive.CanonicalTarNameForPath(relDockerfile)

	excludes = build.TrimBuildFilesFromExcludes(excludes, relDockerfile, false)
	buildCtx, err := archive.TarWithOptions(contextDir, &archive.TarOptions{
		ExcludePatterns: excludes,
		ChownOpts:       &idtools.Identity{UID: 0, GID: 0},
	})
	if err != nil {
		return err
	}

	imageBuildOpts := types.ImageBuildOptions{
		BuildArgs:  projectBuild.Args,
		Labels:     projectBuild.Labels,
		CacheFrom:  projectBuild.CacheFrom,
		Target:     projectBuild.Target,
		Dockerfile: relDockerfile,
	}

	progressOutput := streamformatter.NewProgressOutput(os.Stdout)
	body := progress.NewProgressReader(buildCtx, progressOutput, 0, "", "Sending build context to Docker daemon")

	ctx := context.Background()
	docker, err := client.NewClientWithOpts(client.WithVersion("1.34"), client.FromEnv)
	if err != nil {
		return err
	}

	res, err := docker.ImageBuild(ctx, body, imageBuildOpts)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if _, err := io.Copy(os.Stdout, res.Body); err != nil {
		return err
	}

	return nil
}
