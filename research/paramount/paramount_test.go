package paramount

import "testing"

var videos = []struct {
   justWatch    string
   paramount    string
   resolution   string
   subscription string
}{
   {
      paramount: "https://cbs.com/shows/video/uuwl_4UT4MrVsGwmKFA_FE95RXPmbOMl",
   },
   {
      paramount:    "https://paramountplus.com/shows/video/8PO2sBBr6lFb7J4nklXuzNZRhUR_V9dd",
      justWatch:    "https://justwatch.com/us/tv-show/cia",
      subscription: "FREE",
   },
   {
      paramount:    "https://paramountplus.com/shows/video/ALVE01KKH4B7WREZF804N1RV4TSY4S",
      justWatch:    "https://justwatch.com/us/tv-show/the-price-is-right",
      subscription: "PAID",
   },
   {
      paramount:    "https://paramountplus.com/movies/video/wjQ4RChi6BHHu4MVTncppVuCwu44uq2Q",
      justWatch:    "https://justwatch.com/us/movie/zodiac",
      resolution:   "2160p",
      subscription: "PAID",
   },
}

func Test(t *testing.T) {
   t.Log(videos)
}
