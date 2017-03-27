package itl

import (
	"encoding/json"

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

func (ti TweetsInserter) processTweet(payload string) error {
	tp := &TweetPayload{}
	err := json.Unmarshal([]byte(payload), tp)
	if err != nil {
		glog.Warning(err)
	} else {
		ti.ChartManager.Hit(tp.UserID, tp.CreatedAt, tp.Urls)
	}
	return nil
}
