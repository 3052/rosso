package roku

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "io"
   "net/url"
)

// /api/v3/playback
func FetchPlayback(authToken, rokuId string) (*Playback, error) {
   body, err := json.Marshal(map[string]string{
      "mediaFormat": "DASH",
      "providerId":  "rokuavod",
      "rokuId":      rokuId,
   })
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "googletv.web.roku.com",
         Path:   "/api/v3/playback",
      },
      map[string]string{
         "content-type":         "application/json",
         "user-agent":           user_agent,
         "x-roku-content-token": authToken,
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
   result := &Playback{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

func FetchWidevine(licenseServer string, body []byte) ([]byte, error) {
   target, err := url.Parse(licenseServer)
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      target, map[string]string{"content-type": "application/x-protobuf"}, body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

// codeToken can be empty
//
// /api/v1/account/token
func FetchToken(codeToken string) (*Token, error) {
   header := map[string]string{"user-agent": user_agent}
   if codeToken != "" {
      header["x-roku-content-token"] = codeToken
   }
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "googletv.web.roku.com",
         Path:   "/api/v1/account/token",
      },
      header,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &Token{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

// /api/v1/account/activation/code
func FetchCode(authToken, activationCode string) (*Code, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "googletv.web.roku.com",
         Path:   "/api/v1/account/activation/" + activationCode,
      },
      map[string]string{
         "user-agent":           user_agent,
         "x-roku-content-token": authToken,
      },
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &Code{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

// /api/v1/account/activation
func FetchActivation(authToken string) (*Activation, error) {
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "googletv.web.roku.com",
         Path:   "/api/v1/account/activation",
      },
      map[string]string{
         "content-type":         "application/json",
         "user-agent":           user_agent,
         "x-roku-content-token": authToken,
      },
      []byte(`{"platform": "googletv"}`),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &Activation{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}
