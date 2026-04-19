package mubi

import (
   "41.neocities.org/maya"
   "encoding/base64"
   "encoding/json"
   "errors"
   "fmt"
   "net/url"
)

func (s *Session) FetchSecureUrl(id int) (*SecureUrl, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "api.mubi.com",
         Path:   fmt.Sprintf("/v3/films/%v/viewing/secure_url", id),
      },
      map[string]string{
         "authorization":  "Bearer " + s.Token,
         "client":         client,
         "client-country": ClientCountry,
         "user-agent":     "Firefox",
      },
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result SecureUrl
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.UserMessage != "" {
      return nil, errors.New(result.UserMessage)
   }
   return &result, nil
}

// to get the MPD you have to call this or view video on the website. request
// is hard geo blocked only the first time
func (s *Session) FetchViewing(id int) error {
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "api.mubi.com",
         Path:   fmt.Sprintf("/v3/films/%v/viewing", id),
      },
      map[string]string{
         "authorization":  "Bearer " + s.Token,
         "client":         client,
         "client-country": ClientCountry,
      },
      nil,
   )
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   var result struct {
      UserMessage string `json:"user_message"`
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return err
   }
   if result.UserMessage != "" {
      return errors.New(result.UserMessage)
   }
   return nil
}

func (s *Session) FetchWidevine(body []byte) ([]byte, error) {
   data, err := json.Marshal(map[string]any{
      "merchant":  "mubi",
      "sessionId": s.Token,
      "userId":    s.User.Id,
   })
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "lic.drmtoday.com",
         Path:   "/license-proxy-widevine/cenc/", // final slash is needed
      },
      map[string]string{
         "dt-custom-data": base64.StdEncoding.EncodeToString(data),
      },
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, errors.New(resp.Status)
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
