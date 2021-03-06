package container

import (
	"fmt"
	"strings"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/client"
	"github.com/docker/docker/cli"
	"github.com/spf13/cobra"
)

type waitOptions struct {
	containers []string
}

// NewWaitCommand creates a new cobra.Command for `docker wait`
func NewWaitCommand(dockerCli *client.DockerCli) *cobra.Command {
	var opts waitOptions

	cmd := &cobra.Command{
		Use:   "wait CONTAINER [CONTAINER...]",
		Short: "阻塞直到一个或多个容器停止运行，并打印它们的容器退出码",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.containers = args
			return runWait(dockerCli, &opts)
		},
	}

	return cmd
}

func runWait(dockerCli *client.DockerCli, opts *waitOptions) error {
	ctx := context.Background()

	var errs []string
	for _, container := range opts.containers {
		status, err := dockerCli.Client().ContainerWait(ctx, container)
		if err != nil {
			errs = append(errs, err.Error())
		} else {
			fmt.Fprintf(dockerCli.Out(), "%d\n", status)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, "\n"))
	}
	return nil
}
