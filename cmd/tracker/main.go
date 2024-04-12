package main

import (
	"os"
	"os/signal"
	"syscall"

	_ "modernc.org/sqlite"

	"github.com/msmkdenis/go-task-tracker/internal/app/tracker"
)

func main() {
	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	tracker.Run(quitSignal)
}
