package main

import (
	"flag"
	"os"
	"runtime"

	itl "github.com/realbot/itl"
)

func init() {
	if cpu := runtime.NumCPU(); cpu == 1 {
		runtime.GOMAXPROCS(2)
	} else {
		runtime.GOMAXPROCS(cpu)
	}
}

func main() {
	redisURL := flag.String("redis", "localhost:6379", "Redis address")
	numConsumers := flag.Int("num-consumers", 10, "Number of consumers")
	flag.Parse()

	ti := itl.TweetsInserter{
		TaskManager:  itl.NewTasks("itl", *redisURL),
		ChartManager: itl.NewCharts(itl.NewRedisChartsStore(*redisURL)),
	}
	exitCode := ti.Run(*numConsumers)
	os.Exit(exitCode)
}
