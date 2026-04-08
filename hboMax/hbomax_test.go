package hboMax

import (
   "encoding/xml"
   "os"
   "testing"
)

var content_tests = []struct {
   location   []string
   resolution string
   url        string
}{
   {
      location:   []string{"united states"},
      resolution: "2160p",
      url:        "https://hbomax.com/movies/one-battle-after-another/bebe611d-8178-481a-a4f2-de743b5b135a",
   },
   {
      location: []string{"austria"},
      url:      "https://hbomax.com/at/en/movies/austin-powers-international-man-of-mystery/a979fb8b-f713-4de3-a625-d16ad4d37448",
   },
   {
      location: []string{
         "belgium", "brazil", "bulgaria", "chile", "colombia", "croatia",
         "czech republic", "denmark", "finland", "france", "hungary",
         "indonesia", "malaysia", "mexico", "netherlands", "norway", "peru",
         "philippines", "poland", "portugal", "romania", "singapore", "slovakia",
         "spain", "sweden", "thailand", "united states",
      },
      url: "https://hbomax.com/shows/white-lotus/14f9834d-bc23-41a8-ab61-5c8abdbea505",
   },
}

func TestContent(t *testing.T) {
   t.Log(content_tests)
}

func TestPlayReady(t *testing.T) {
   dir, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(dir + "/hboMax/userCache.xml")
   if err != nil {
      t.Fatal(err)
   }
   var cache struct {
      Login Login
   }
   err = xml.Unmarshal(data, &cache)
   if err != nil {
      t.Fatal(err)
   }
   // hbomax.com/movies/dune/e7dc7b3a-a494-4ef1-8107-f4308aa6bbf7
   play, err := cache.Login.PlayReady("06a38397-862d-4419-be84-0641939825e7")
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(
      dir+"/hboMax/PlayReady",
      []byte(play.Drm.Schemes.PlayReady.LicenseUrl), os.ModePerm,
   )
   if err != nil {
      t.Fatal(err)
   }
}
