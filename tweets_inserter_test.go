package itl

import "testing"

func TestValidateURL(t *testing.T) {
	ti := &TweetsInserter{}
	if ti.validateURL("https://twitter.com/foo/bar") != false {
		t.Error("Must skip twitter urls")
	}
	if ti.validateURL("https://facebook.com/foo/bar") != false {
		t.Error("Must skip facebook urls")
	}
	if ti.validateURL("https://plus.google.com/foo/bar") != false {
		t.Error("Must skip plus urls")
	}
}

func TestValidateTweet(t *testing.T) {
	ti := &TweetsInserter{}
	tweet := &TweetPayload{}

	tweet.UserFriendsCount = 1
	tweet.UserFollowersCount = 100
	if !ti.validateTweet(tweet) {
		t.Error("Must keep >= 0,01 ratio")
	}

	tweet.UserFriendsCount = 1
	tweet.UserFollowersCount = 101
	if ti.validateTweet(tweet) {
		t.Error("Must remove < 0,01 ratio")
	}

	tweet.UserFriendsCount = 100
	tweet.UserFollowersCount = 1
	if !ti.validateTweet(tweet) {
		t.Error("Must keep <= 100 ratio")
	}

	tweet.UserFriendsCount = 101
	tweet.UserFollowersCount = 1
	if ti.validateTweet(tweet) {
		t.Error("Must remove > 100 ratio")
	}
}
