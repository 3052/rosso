package itv

import "testing"

var watch_tests = []struct {
   category string
   watch    string
}{
   {
      category: "FILM",
      watch:    "https://itv.com/watch/dune/10a6768a0001B",
   },
   {
      category: "DRAMA_AND_SOAPS",
      watch:    "https://itv.com/watch/joan/10a3918",
   },
}

func TestWatch(t *testing.T) {
   t.Log(watch_tests)
}
