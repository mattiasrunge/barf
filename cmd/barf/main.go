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
	// TODO: Check for updates and download and install updated binary
	// TODO: Don't show progress bar at all, fire and forget operation for scripts
	// TODO: JSON output for scripting
	// TODO: List history
	// TODO: Abort operation
	// TODO: Version negotioation with daemon, restart daemon if wrong version

	app.Command("list l", "list active operations", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			run.StartCLI(*width, func() error {
				return actions.List(map[string]interface{}{})
			})
		}
	})

	app.Command("monitor m", "monitors active operations", func(cmd *cli.Cmd) {
		cmd.LongDesc = "monitors active operations, if ids are given it will exit when those operations have finished"

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

	app.Command("abort a", "aborts an active operation", func(cmd *cli.Cmd) {
		cmd.Spec = "ID"
		id := cmd.StringArg("ID", "", "ID to abort")

		cmd.Action = func() {
			run.StartCLI(*width, func() error {
				return actions.Abort(map[string]interface{}{
					"id": id,
				})
			})
		}
	})

	app.Command("copy cp", "copies files or directories", func(cmd *cli.Cmd) {
		cmd.Spec = "SRC... DST"
		src := cmd.StringsArg("SRC", nil, "Source to copy")
		dst := cmd.StringArg("DST", "", "Destination where to copy to")

		cmd.Action = func() {
			run.StartCLI(*width, func() error {
				return actions.Copy(map[string]interface{}{
					"src": src,
					"dst": dst,
				})
			})
		}
	})

	app.Command("move mv", "moves files or directories", func(cmd *cli.Cmd) {
		cmd.Spec = "SRC... DST"
		src := cmd.StringsArg("SRC", nil, "Source to move")
		dst := cmd.StringArg("DST", "", "Destination where to move to")

		cmd.Action = func() {
			run.StartCLI(*width, func() error {
				return actions.Move(map[string]interface{}{
					"src": src,
					"dst": dst,
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
