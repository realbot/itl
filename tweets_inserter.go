package itl

import (
	"encoding/json"
	"strings"

	"github.com/golang/glog"
)

const inserterVersion = "1.0.0"

type TweetsInserter struct {
	TaskManager  *Tasks
	ChartManager *Charts
}

func (ti TweetsInserter) Run(numConsumers int) int {
	glog.Infof("itl inserter version %s\n", inserterVersion)
	ti.TaskManager.StartConsumers(numConsumers, ti.processTweet)

	waitForExit()

	return ExitCodeOK
}

func (ti TweetsInserter) validateTweet(tweet *TweetPayload) bool {
	ratio := float64(tweet.UserFriendsCount) / float64(tweet.UserFollowersCount)
	return ratio >= 0.01 && ratio <= 100
}

func (ti TweetsInserter) validateURL(url string) bool {
	if strings.HasPrefix(url, "https://twitter.com") ||
		strings.HasPrefix(url, "https://facebook.com") ||
		strings.HasPrefix(url, "https://www.facebook.com") ||
		strings.HasPrefix(url, "https://plus.google.com") ||
		strings.HasPrefix(url, "https://youtu.be") ||
		strings.HasPrefix(url, "http://www.theverge.com") ||
	    strings.HasPrefix(url, "https://www.theverge.com") {
		return false
	}
	return true
}

func (ti TweetsInserter) weigth(tweet *TweetPayload) float64 {
	if tweet.FavoriteCount > 0 || tweet.RetweetCount > 0 {
		return 1.2
	}
	return 1
}

func (ti TweetsInserter) processTweet(payload string) error {
	tp := &TweetPayload{}
	err := json.Unmarshal([]byte(payload), tp)
	if err != nil {
		glog.Warning(err)
	} else {
		for _, url := range tp.Urls {
			if ti.validateURL(url) {
				glog.Infof("%f %s", ti.weigth(tp), url)
				ti.ChartManager.Hit(tp.UserID, tp.CreatedAt, url, ti.weigth(tp))
			}
		}
	}
	return nil
}
