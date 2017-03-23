package itl

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	ExitCodeOK = iota
	ExitCodeError
)

func waitForExit() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
}
