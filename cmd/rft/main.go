package main

import (
	"fmt"
	"os"

	"rft/cmd/rft/run"
	"rft/internal/actions"
	"rft/internal/config"

	cli "github.com/jawher/mow.cli"
)

func main() {
	run.CheckDaemon()

	app := cli.App(config.Name, config.Description)
	app.Version("v version", fmt.Sprintf("%s\n%s\n%s", config.Version, config.BuildChecksum, config.BuildTime))

	app.Action = func() {
		run.StartCLI(func() error {
			return actions.Monitor(map[string]interface{}{})
		})
	}

	// TODO: --log, -l flag to output raw rsync stdout/stderr with the progressbar on bottom... how to handle multi monitor... log only a few lines for each or only support on single monitor?

	app.Command("monitor m", "monitors active operations", func(cmd *cli.Cmd) {
		cmd.Spec = "[IDS...]"
		ids := cmd.StringsArg("IDS", nil, "IDs to monitor")

		cmd.Action = func() {
			run.StartCLI(func() error {
				return actions.Monitor(map[string]interface{}{
					"ids": ids,
				})
			})
		}
	})

	app.Command("copy cp", "copies files or folders", func(cmd *cli.Cmd) {
		cmd.Spec = "SRC... DST"
		src := cmd.StringsArg("SRC", nil, "Source files to copy")
		dst := cmd.StringArg("DST", "", "Destination where to copy files to")

		cmd.Action = func() {
			run.StartCLI(func() error {
				return actions.Copy(map[string]interface{}{
					"from": src,
					"to":   dst,
				})
			})
		}
	})

	app.Command("dummy", "starts dummy operations", func(cmd *cli.Cmd) {
		cmd.Spec = "[ITER]"
		i := cmd.StringArg("ITER", "10", "Iterations to run")

		cmd.Action = func() {
			run.StartCLI(func() error {
				return actions.Dummy(map[string]interface{}{
					"iterations": i,
				})
			})
		}
	})

	app.Run(os.Args)
}
