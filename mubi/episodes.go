package mubi

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
)

func FetchSeason(slug string, season int) (*Season, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme: "https",
         Host: "api.mubi.com",
         Path: fmt.Sprintf(
            "/v4/series/%v/seasons/season-%v/episodes", slug, season,
         ),
      },
      Header: http.Header{},
   }
   req.Header.Set("client", client)
   req.Header.Set("client-country", ClientCountry)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &Season{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

type Season struct {
   Episodes []Film
}
