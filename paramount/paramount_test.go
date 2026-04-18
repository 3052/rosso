package paramount

import (
   "errors"
   "io"
   "net/http"
   "net/url"
   "testing"
   "time"
)

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

func brands(host, app_secret string) error {
   at, err := get_at(app_secret)
   if err != nil {
      return err
   }
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{}
   req.URL.Host = host
   req.URL.Path = "/apps-api/v3.0/androidphone/brands/.json"
   value := url.Values{}
   value["at"] = []string{at}
   req.URL.RawQuery = value.Encode()
   req.URL.Scheme = "https"
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   _, err = io.Copy(io.Discard, resp.Body)
   if err != nil {
      return err
   }
   if resp.StatusCode != http.StatusOK {
      return errors.New(resp.Status)
   }
   return nil
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
