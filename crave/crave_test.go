package crave

import "testing"

var tests = []struct {
   resolution string
   url        string
}{
   {
      resolution: "1080p",
      url:        "https://crave.ca/movie/anaconda-2025-59881",
   },
   {
      resolution: "1080p",
      url:        "https://crave.ca/play/anaconda-2025-3300246",
   },
   {
      resolution: "2160p",
      url:        "https://crave.ca/play/heated-rivalry/ill-believe-in-anything-s1e5-3233873",
   },
}

func Test(t *testing.T) {
   t.Log(tests)
}
