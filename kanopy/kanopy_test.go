package kanopy

import "testing"

var tests = []struct {
   genre string
   url   string
}{
   {
      genre: "Movies",
      url:   "https://kanopy.com/video/13808102",
   },
   {
      genre: "TV Series",
      url:   "https://kanopy.com/video/14098194",
   },
}

func TestLog(t *testing.T) {
   t.Log(tests)
}
