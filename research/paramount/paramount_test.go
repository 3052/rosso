package paramount

import (
   "fmt"
   "os/exec"
   "strings"
   "testing"
)

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
      justWatch:  "https://justwatch.com/us/tv-show/the-price-is-right",
      paramount:  "https://paramountplus.com/shows/video/ALVE01KKH4B7WREZF804N1RV4TSY4S",
      resolution: "1080p",
      cookie:     true,
   },
   {
      justWatch:  "https://justwatch.com/us/movie/zodiac",
      resolution: "2160p",
      paramount:  "https://paramountplus.com/movies/video/wjQ4RChi6BHHu4MVTncppVuCwu44uq2Q",
      cookie:     true,
   },
}

func TestStreamingUrl(t *testing.T) {
   const id = "com.cbs.tve"
   app_data, ok := Apps[id]
   if !ok {
      t.Fatal(id)
   }
   token_data, err := app_data.FetchStreamingUrl(
      "uuwl_4UT4MrVsGwmKFA_FE95RXPmbOMl", nil,
   )
   if err != nil {
      t.Fatal(err)
   }
   t.Logf("%+v", token_data)
}

func TestCookie(t *testing.T) {
   username, err := run("credential", "-h=paramountplus.com", "-k=username")
   if err != nil {
      t.Fatal(err)
   }
   password, err := run("credential", "-h=paramountplus.com", "-k=password")
   if err != nil {
      t.Fatal(err)
   }
   const id = "com.cbs.app"
   app_data, ok := Apps[id]
   if !ok {
      t.Fatal(id)
   }
   cookie, err := app_data.FetchCbsCom(username, password)
   if err != nil {
      t.Fatal(err)
   }
   t.Log(cookie)
}

func run(name string, arg ...string) (string, error) {
   var data strings.Builder
   command := exec.Command(name, arg...)
   command.Stdout = &data
   fmt.Println(command.Args)
   err := command.Run()
   if err != nil {
      return "", err
   }
   return data.String(), nil
}
