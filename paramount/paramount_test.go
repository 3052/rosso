package paramount

import "testing"

var videos = []struct {
   justWatch string
   paramount string
   cookie    bool
   height    int
}{
   {
      cookie:    false,
      height:    2160,
      justWatch: "https://justwatch.com/us/tv-show/dexter-original-sin",
      paramount: "https://paramountplus.com/shows/video/6z6Nn5LjJqCraw_6dSSL1dY_j__kv0EH",
   },
   {
      cookie:    true,
      height:    2160,
      justWatch: "https://justwatch.com/us/movie/zodiac",
      paramount: "https://paramountplus.com/movies/video/wjQ4RChi6BHHu4MVTncppVuCwu44uq2Q",
   },
   {
      cookie:    true,
      height:    1080,
      justWatch: "https://justwatch.com/us/tv-show/the-price-is-right",
      paramount: "https://paramountplus.com/shows/video/ALVE01KKH4B7WREZF804N1RV4TSY4S",
   },
   {
      cookie:    false,
      height:    1080,
      justWatch: "https://justwatch.com/us/tv-show/60-minutes",
      paramount: "https://cbs.com/shows/video/uuwl_4UT4MrVsGwmKFA_FE95RXPmbOMl",
   },
}

func TestVideos(t *testing.T) {
   t.Log(videos)
}
