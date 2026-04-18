package paramount

import (
   "41.neocities.org/maya"
   "errors"
   "io"
   "net/url"
   "testing"
   "time"
)

func brands(host, app_secret string) error {
   at, err := get_at(app_secret)
   if err != nil {
      return err
   }
   resp, err := maya.Get(
      &url.URL{
         Scheme:   "https",
         Host:     host,
         Path:     "/apps-api/v3.0/androidphone/brands/.json",
         RawQuery: url.Values{"at": {at}}.Encode(),
      },
      nil,
   )
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   _, err = io.Copy(io.Discard, resp.Body)
   if err != nil {
      return err
   }
   if resp.StatusCode != 200 {
      return errors.New(resp.Status)
   }
   return nil
}

func TestDexParamount(t *testing.T) {
   results, err := ExtractDexHexBytes("base.apk")
   if err != nil {
      t.Fatal(err)
   }
   var sleep bool
   for result := range results {
      if sleep {
         time.Sleep(time.Second)
      } else {
         sleep = true
      }
      t.Log(brands("www.paramountplus.com", result), result)
   }
}

func TestDexCbs(t *testing.T) {
   results, err := ExtractDexHexBytes("base.apk")
   if err != nil {
      t.Fatal(err)
   }
   var sleep bool
   for result := range results {
      if sleep {
         time.Sleep(time.Second)
      } else {
         sleep = true
      }
      t.Log(brands("www.cbs.com", result), result)
   }
}

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
   {
      cookie:     false,
      paramount:  "https://paramountplus.com/shows/video/ALVE01KMDREQKEENRS8QS6BASFR1TA",
      resolution: "1080p",
      justWatch:  "https://justwatch.com/us/tv-show/the-bold-and-the-beautiful",
   },
}
