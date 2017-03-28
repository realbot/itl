package itl

import (
	"reflect"
	"testing"
)

func TestAddURL(t *testing.T) {
	fs := &FakeStore{}
	charts := NewCharts(fs)

	charts.Hit("someuserid", "Fri Mar 24 13:08:45 +0000 2017", "someurl", 1)

	expectedKeys := []string{"itl-d-someuserid-24mar2017", "itl-m-someuserid-mar2017", "itl-g-someuserid"}
	expectedURLs := []string{"someurl", "someurl", "someurl"}
	if !reflect.DeepEqual(expectedKeys, fs.LastKeys) {
		t.Errorf("Wrong keys expected %v, got %v", expectedKeys, fs.LastKeys)
	}
	if !reflect.DeepEqual(expectedURLs, fs.LastURLs) {
		t.Errorf("Wrong urls expected %v, got %v", expectedURLs, fs.LastURLs)
	}
}

type FakeStore struct {
	LastKeys, LastURLs []string
}

func (fs *FakeStore) update(key, url string, weight float64) {
	fs.LastKeys = append(fs.LastKeys, key)
	fs.LastURLs = append(fs.LastURLs, url)
}
