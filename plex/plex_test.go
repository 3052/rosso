package plex

import "testing"

var watch_tests = []struct {
   drm bool
   url string
}{
   {
      url: "https://watch.plex.tv/movie/limitless",
   },
   {
      url: "https://watch.plex.tv/show/broadchurch/season/3/episode/5",
   },
   {
      drm: true,
      url: "https://watch.plex.tv/movie/ghost-in-the-shell",
   },
}

func Test(t *testing.T) {
   t.Log(watch_tests)
}
