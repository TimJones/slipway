package slipway

import (
	"fmt"
	"io"
	"os"

	"github.com/docker/cli/cli/command/image/build"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/pkg/streamformatter"
	"golang.org/x/net/context"
)

func (project *Project) Build() error {
	projectBuild := project.Container.Build

	oldImageID, err := getImageID(fmt.Sprintf("%s:slipway", project.Name))
	if err != nil {
		return err
	}

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
		Tags:       []string{fmt.Sprintf("%s:slipway", project.Name)},
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

	newImageID, err := getImageID(fmt.Sprintf("%s:slipway", project.Name))
	if err != nil {
		return err
	}

	if oldImageID != "" && oldImageID != newImageID {
		_, err = docker.ImageRemove(ctx, oldImageID, types.ImageRemoveOptions{
			Force:         false,
			PruneChildren: true,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func getImageID(name string) (string, error) {
	ctx := context.Background()
	docker, err := client.NewClientWithOpts(client.WithVersion("1.34"), client.FromEnv)
	if err != nil {
		return "", err
	}

	images, err := docker.ImageList(ctx, types.ImageListOptions{
		All: false,
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "reference",
			Value: name,
		}),
	})

	if err != nil {
		return "", err
	}

	if len(images) == 0 {
		return "", nil
	} else if len(images) > 1 {
		return "", fmt.Errorf("got %d images named %s", len(images), name)
	}

	return images[0].ID, nil
}
