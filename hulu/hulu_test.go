package hulu

import "testing"

var tests = []struct {
   height    int
   hulu      string
   justWatch string
}{
   {
      height:    2160,
      hulu:      "https://hulu.com/movie/stay-5742941d-4b4a-4914-8774-f5d8d57f9382",
      justWatch: "https://justwatch.com/us/movie/stay-2025",
   },
   {
      height:    1080,
      hulu:      "https://hulu.com/movie/palm-springs-f70dfd4d-dbfb-46b8-abb3-136c841bba11",
      justWatch: "https://justwatch.com/us/movie/palm-springs",
   },
   {
      hulu: "https://hulu.com/series/house-ef39603f-eb90-4248-8237-f6168d7c1be1",
   },
}

func Test(t *testing.T) {
   t.Log(tests)
}
