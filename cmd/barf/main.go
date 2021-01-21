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

	app.Command("update u", "check for updates", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			err := actions.Update(map[string]interface{}{})

			if err != nil {
				fmt.Println(err)
				os.Exit(255)
			}
		}
	})

	app.Run(os.Args)
}
