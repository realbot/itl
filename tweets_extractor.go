package itl

import (
	"encoding/json"
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/golang/glog"
)

const extractorVersion = "1.0.0"

type TweetsExtractor struct {
	TaskManager                                            *Tasks
	ConsumerKey, ConsumerSecret, AccessToken, AccessSecret string
	UserID                                                 string
}

type TweetPayload struct {
	UserID             string   `json:"user_id"`
	UserFriendsCount   int      `json:"user_friends_count"`
	UserFollowersCount int      `json:"user_followers_count"`
	Urls               []string `json:"urls"`
	RetweetCount       int      `json:"retweet_count"`
	FavoriteCount      int      `json:"favorite_count"`
	CreatedAt          string   `json:"created_at"`
}

func (te TweetsExtractor) Run() int {
	glog.Infof("itl extractor version %s\n", extractorVersion)

	twclient := te.createTwitterClient()
	stream := te.createTwitterStream(te.UserID, twclient)

	demux := twitter.NewSwitchDemux()
	demux.Tweet = te.processTweet
	go demux.HandleChan(stream.Messages)

	waitForExit()

	glog.Info("Stopping Stream...")
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
		glog.Fatal(err)
	}
	return
}

func (te TweetsExtractor) processTweet(tweet *twitter.Tweet) {
	if len(tweet.Entities.Urls) > 0 {
		var urls = []string{}
		for _, url := range tweet.Entities.Urls {
			if url.ExpandedURL != "" {
				urls = append(urls, url.ExpandedURL)
			}
		}
		if len(urls) > 0 {
			payload := TweetPayload{
				UserID:             te.UserID,
				UserFriendsCount:   tweet.User.FriendsCount,
				UserFollowersCount: tweet.User.FollowersCount,
				RetweetCount:       tweet.RetweetCount,
				FavoriteCount:      tweet.FavoriteCount,
				CreatedAt:          tweet.CreatedAt,
				Urls:               urls,
			}
			plb, err := json.Marshal(payload)
			if err != nil {
				glog.Warning(err)
			} else {
				te.TaskManager.EnqueueTask(string(plb))
			}
		}
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
		glog.Fatal(err)
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
		glog.Fatal(err)
	}
	return
}
