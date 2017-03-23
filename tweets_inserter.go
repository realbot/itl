package itl

import (
	"fmt"
	"io"
	"log"
)

const inserterVersion = "1.0.0"

type TweetsInserter struct {
	Out, Err    io.Writer
	TaskManager *Tasks
}

func (ti TweetsInserter) Run(numConsumers int) int {
	fmt.Fprintf(ti.Out, "itl inserter version %s\n", inserterVersion)
	ti.TaskManager.StartConsumers(numConsumers, ti.processTweet)

	waitForExit()

	return ExitCodeOK
}

func (ti TweetsInserter) processTweet(payload string) error {
	log.Println(payload)
	return nil
}
