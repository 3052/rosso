package amc

import "testing"

var tests = []struct {
   content string
   url     string
}{
   {
      content: "movie",
      url:     "https://amcplus.com/movies/nocebo--1061554",
   },
   {
      content: "season",
      url:     "https://amcplus.com/shows/orphan-black/episodes--1010638",
   },
   {
      content: "series",
      url:     "https://amcplus.com/shows/orphan-black--1010578",
   },
}

func Test(t *testing.T) {
   t.Log(tests)
}
