package itl

import (
	"fmt"
	"strings"
	"time"
)

const keyprefix = "itl"

type Charts struct {
	store ChartsStore
}

func NewCharts(store ChartsStore) *Charts {
	return &Charts{
		store: store,
	}
}

func (c Charts) Hit(userid, date string, url string) {
	t := c.parseDate(date)
	c.store.update(c.dayKey(userid, t), url)
	c.store.update(c.monthKey(userid, t), url)
	c.store.update(c.globalKey(userid, t), url)
}

func (c Charts) dayKey(userid string, t time.Time) string {
	return fmt.Sprintf("%s-d-%s-%s", keyprefix, userid, strings.ToLower(t.Format("02Jan2006")))
}

func (c Charts) monthKey(userid string, t time.Time) string {
	return fmt.Sprintf("%s-m-%s-%s", keyprefix, userid, strings.ToLower(t.Format("Jan2006")))
}

func (c Charts) globalKey(userid string, t time.Time) string {
	return fmt.Sprintf("%s-g-%s", keyprefix, userid)
}

func (c Charts) parseDate(date string) time.Time {
	t, err := time.Parse(time.RubyDate, date)
	if err != nil {
		t = time.Now()
	}
	return t
}
