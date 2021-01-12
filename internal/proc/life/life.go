package life

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type closeFunction func()

var shutdownHooks = []closeFunction{}

func Start() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func(c chan os.Signal) {
		sig := <-c

		fmt.Printf("Caught signal %s: Running %d hooks...\n", sig, len(shutdownHooks))

		RunShutdownHooks()

		fmt.Printf("Shutdown complete, goodbye!\n")

		os.Exit(0)
	}(c)
}

func RunShutdownHooks() {
	for _, fn := range shutdownHooks {
		fn()
	}
}

func AddShutdownHook(fn closeFunction) {
	shutdownHooks = append(shutdownHooks, fn)
}
