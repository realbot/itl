package main

import (
	"flag"
	"log"
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
	consumerKey := flag.String("consumer-key", "", "Twitter Consumer Key")
	consumerSecret := flag.String("consumer-secret", "", "Twitter Consumer Secret")
	accessToken := flag.String("access-token", "", "Twitter Access Token")
	accessSecret := flag.String("access-secret", "", "Twitter Access Secret")
	userID := flag.String("userid", "", "Twitter User ID")
	flag.Parse()

	if *userID == "" {
		log.Fatal("User ID required")
	}

	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
	}

	te := itl.TweetsExtractor{
		TaskManager: itl.NewTasks("itl", *redisURL),
		ConsumerKey: *consumerKey, ConsumerSecret: *consumerSecret,
		AccessToken: *accessToken, AccessSecret: *accessSecret,
		UserID: *userID,
	}
	exitCode := te.Run()
	os.Exit(exitCode)
}
