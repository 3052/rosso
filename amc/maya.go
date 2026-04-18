package amc

import (
   "41.neocities.org/maya"
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
)

func Playback(authToken string, videoID int) (*PlaybackResult, error) {
   url := fmt.Sprintf("https://gw.cds.amcn.com/playback-id/api/v1/playback/%d", videoID)

   payload := []byte(`{"adtags":{"lat":0,"mode":"on-demand","playerHeight":0,"playerWidth":0,"ppid":0,"url":"-"}}`)

   req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
   if err != nil {
      return nil, err
   }

   req.Header.Set("authorization", "Bearer "+authToken)
   req.Header.Set("content-type", "application/json")
   req.Header.Set("x-amcn-language", "en")
   req.Header.Set("x-amcn-network", "amcplus")
   req.Header.Set("x-amcn-platform", "web")
   req.Header.Set("x-amcn-service-id", "amcplus")
   req.Header.Set("x-amcn-tenant", "amcn")
   req.Header.Set("x-amcn-device-ad-id", "-")
   req.Header.Set("x-ccpa-do-not-sell", "doNotPassData")

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("playback failed with status: %d", resp.StatusCode)
   }

   // Internal envelope to strip the first layer
   var envelope struct {
      Success bool         `json:"success"`
      Status  int          `json:"status"`
      Data    PlaybackData `json:"data"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
      return nil, err
   }

   return &PlaybackResult{
      Data:     envelope.Data,
      BcovAuth: resp.Header.Get("x-amcn-bc-jwt"),
   }, nil
}

// Login authenticates the user. It requires the guest token (access_token)
// retrieved from calling the Unauth() function.
func Login(guestToken, email, password string) (*AuthData, error) {
   // Body
   body, err := json.Marshal(map[string]string{
      "email":    email,
      "password": password,
   })
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "gw.cds.amcn.com",
         Path:   "/auth-orchestration-id/api/v1/login",
      },
      map[string]string{
         "authorization":           "Bearer " + guestToken,
         "content-type":            "application/json",
         "x-amcn-language":         "en",
         "x-amcn-network":          "amcplus",
         "x-amcn-platform":         "web",
         "x-amcn-service-group-id": "10",
         "x-amcn-tenant":           "amcn",
         "x-amcn-device-ad-id":     "-",
         "x-amcn-device-id":        "-",
         "x-amcn-service-id":       "amcplus",
         "x-ccpa-do-not-sell":      "doNotPassData",
      },
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("login failed with status: %d", resp.StatusCode)
   }
   // Internal envelope to strip the first layer
   var envelope struct {
      Success bool     `json:"success"`
      Status  int      `json:"status"`
      Data    AuthData `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
      return nil, err
   }
   return &envelope.Data, nil
}
