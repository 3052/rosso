package paramount

import "testing"

func TestVideos(t *testing.T) {
   t.Log(videos)
}

var videos = []struct {
   justWatch  string
   paramount  string
   resolution string
   cookie     bool
}{
   {
      justWatch:  "https://justwatch.com/us/tv-show/cia",
      paramount:  "https://paramountplus.com/shows/video/8PO2sBBr6lFb7J4nklXuzNZRhUR_V9dd",
      resolution: "1080p",
      cookie:     false,
   },
   {
      justWatch:  "https://justwatch.com/us/movie/zodiac",
      resolution: "2160p",
      paramount:  "https://paramountplus.com/movies/video/wjQ4RChi6BHHu4MVTncppVuCwu44uq2Q",
      cookie:     true,
   },
   {
      justWatch:  "https://justwatch.com/us/tv-show/the-price-is-right",
      paramount:  "https://paramountplus.com/shows/video/ALVE01KKH4B7WREZF804N1RV4TSY4S",
      resolution: "1080p",
      cookie:     true,
   },
   {
      justWatch:  "https://justwatch.com/us/tv-show/60-minutes",
      paramount:  "https://cbs.com/shows/video/uuwl_4UT4MrVsGwmKFA_FE95RXPmbOMl",
      resolution: "1080p",
      cookie:     false,
   },
}
