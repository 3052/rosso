package itv

import "testing"

var watch_tests = []struct {
   category string
   watch    string
}{
   {
      category: "FILM",
      watch:    "https://itv.com/watch/mission-impossible-fallout/10a7086a0001B",
   },
   {
      category: "DRAMA_AND_SOAPS",
      watch:    "https://itv.com/watch/joan/10a3918",
   },
}

func TestWatch(t *testing.T) {
   t.Log(watch_tests)
}
