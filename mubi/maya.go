package mubi

import (
   "41.neocities.org/maya"
   "bytes"
   "encoding/base64"
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
)

func (s *Session) FetchWidevine(body []byte) ([]byte, error) {
   // final slash is needed
   req, err := http.NewRequest(
      "POST", "https://lic.drmtoday.com/license-proxy-widevine/cenc/",
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   data, err := json.Marshal(map[string]any{
      "merchant":  "mubi",
      "sessionId": s.Token,
      "userId":    s.User.Id,
   })
   if err != nil {
      return nil, err
   }

   req.Header.Set("dt-custom-data", base64.StdEncoding.EncodeToString(data))

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   // Check if the response is not a 200 OK
   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("unexpected HTTP error %v", resp.StatusCode)
   }
   var result struct {
      License []byte
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return result.License, nil
}

func FetchFilm(slug string) (*Film, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https", Host: "api.mubi.com", Path: "/v3/films/" + slug,
      },
      map[string]string{
         "client":         client,
         "client-country": ClientCountry,
      },
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &Film{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

func FetchEpisodes(slug string, season int) ([]Film, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "api.mubi.com",
         Path:   fmt.Sprintf("/v4/series/%v/seasons/season-%v/episodes", slug, season),
      },
      map[string]string{
         "client":         client,
         "client-country": ClientCountry,
      },
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Episodes []Film
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return result.Episodes, nil
}

func FetchLinkCode() (*LinkCode, error) {
   resp, err := maya.Get(
      &url.URL{Scheme: "https", Host: "api.mubi.com", Path: "/v3/link_code"},
      map[string]string{
         "client":         client,
         "client-country": ClientCountry,
      },
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &LinkCode{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

func (l *LinkCode) FetchSession() (*Session, error) {
   body, err := json.Marshal(map[string]string{"auth_token": l.AuthToken})
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "api.mubi.com",
         Path:   "/v3/authenticate",
      },
      map[string]string{
         "client":         client,
         "client-country": ClientCountry,
         "content-type":   "application/json",
      },
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &Session{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}
