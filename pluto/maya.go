package pluto

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "io"
   "net/url"
   "strings"
)

// pluto.tv/on-demand/movies/64946365c5ae350013623630
// pluto.tv/on-demand/movies/disobedience-ca-2018-1-1
func FetchSeries(movieShow string) (*Series, error) {
   data := url.Values{}
   data.Set("appName", app_name)
   data.Set("appVersion", "9")
   data.Set("clientID", "9")
   data.Set("clientModelNumber", "9")
   data.Set("deviceMake", "9")
   data.Set("deviceModel", "9")
   data.Set("deviceVersion", "9")
   data.Set("drmCapabilities", drm_capabilities)
   if strings.Contains(movieShow, "-") {
      data.Set("episodeSlugs", movieShow)
   } else {
      data.Set("seriesIDs", movieShow)
   }
   resp, err := maya.Get(
      &url.URL{
         Scheme:   "https",
         Host:     "boot.pluto.tv",
         Path:     "/v4/start",
         RawQuery: data.Encode(),
      },
      nil,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Series
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if strings.Contains(movieShow, "-") {
      if result.Vod[0].Slug != movieShow {
         return nil, errors.New("slug mismatch")
      }
   } else if result.Vod[0].Id != movieShow {
      return nil, errors.New("id mismatch")
   }
   return &result, nil
}

func FetchWidevine(body []byte) ([]byte, error) {
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "service-concierge.clusters.pluto.tv",
         Path:   "/v1/wv/alt",
      },
      map[string]string{"content-type": "application/x-protobuf"},
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}
