package main

import (
	"fmt"
	"os"

	"barf/cmd/barf/run"
	"barf/internal/actions"
	"barf/internal/config"

	cli "github.com/jawher/mow.cli"
)

func main() {
	run.CheckDaemon()

	app := cli.App(config.Name, config.Description)
	app.Version("v version", fmt.Sprintf("%s\n%s\n%s", config.Version, config.BuildChecksum, config.BuildTime))

	width := app.IntOpt("width, w", 0, "terminal width to use, if not set (or zero) it will be auto detected and failing that set to 132")

	app.Action = func() {
		run.StartCLI(*width, func() error {
			return actions.Monitor(map[string]interface{}{})
		})
	}

	// TODO: --log, -l flag to output raw rsync stdout/stderr with the progressbar on bottom... how to handle multi monitor... log only a few lines for each or only support on single monitor?

	app.Command("monitor m", "monitors active operations", func(cmd *cli.Cmd) {
		cmd.Spec = "[IDS...]"
		ids := cmd.StringsArg("IDS", nil, "IDs to monitor")

		cmd.Action = func() {
			run.StartCLI(*width, func() error {
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
			run.StartCLI(*width, func() error {
				return actions.Copy(map[string]interface{}{
					"from": src,
					"to":   dst,
				})
			})
		}
	})

	if !config.IsProduction() {
		app.Command("dummy", "starts dummy operations", func(cmd *cli.Cmd) {
			cmd.Spec = "[ITER]"
			i := cmd.StringArg("ITER", "10", "Iterations to run")

			cmd.Action = func() {
				run.StartCLI(*width, func() error {
					return actions.Dummy(map[string]interface{}{
						"iterations": i,
					})
				})
			}
		})
	}

	app.Run(os.Args)
}
