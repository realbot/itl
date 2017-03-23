package itl

import (
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

const extractorVersion = "1.0.0"

type TweetsExtractor struct {
	Out, Err                                               io.Writer
	TaskManager                                            *Tasks
	ConsumerKey, ConsumerSecret, AccessToken, AccessSecret string
}

func (te TweetsExtractor) Run(userid string) int {
	fmt.Fprintf(te.Out, "itl extractor version %s\n", extractorVersion)

	twclient := te.createTwitterClient()
	stream := te.createTwitterStream(userid, twclient)

	demux := twitter.NewSwitchDemux()
	demux.Tweet = te.processTweet
	go demux.HandleChan(stream.Messages)

	waitForExit()

	fmt.Fprintln(te.Out, "Stopping Stream...")
	stream.Stop()

	return ExitCodeOK
}

func (te TweetsExtractor) createTwitterStream(userid string, twclient *twitter.Client) (stream *twitter.Stream) {
	friends := te.friendsOf(userid, twclient)
	filterParams := &twitter.StreamFilterParams{
		Follow:        friends,
		StallWarnings: twitter.Bool(true),
	}
	stream, err := twclient.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (te TweetsExtractor) processTweet(tweet *twitter.Tweet) {
	fmt.Fprintln(te.Out, tweet.Text)
	if len(tweet.Entities.Urls) > 0 {
		te.TaskManager.EnqueueTask(tweet.Entities.Urls[0].ExpandedURL)
		fmt.Fprintln(te.Out, tweet.Entities.Urls[0].ExpandedURL)
		fmt.Fprintln(te.Out, "r:"+strconv.Itoa(tweet.RetweetCount)+" f:"+strconv.Itoa(tweet.FavoriteCount))
	}
}

func (te TweetsExtractor) createTwitterClient() (twclient *twitter.Client) {
	config := oauth1.NewConfig(te.ConsumerKey, te.ConsumerSecret)
	token := oauth1.NewToken(te.AccessToken, te.AccessSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	twclient = twitter.NewClient(httpClient)
	return
}

func (te TweetsExtractor) friendsOf(userID string, twclient *twitter.Client) (friends []string) {
	uid, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	params := &twitter.FriendIDParams{
		UserID: uid,
		Count:  500,
		Cursor: -1,
	}
	friendsIDs, _, err := twclient.Friends.IDs(params)
	if err == nil {
		for _, f := range friendsIDs.IDs {
			friends = append(friends, strconv.FormatInt(f, 10))
		}
	} else {
		log.Fatal(err)
	}
	return
}
