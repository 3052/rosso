package amc

import (
   "41.neocities.org/maya"
   _ "embed"
   "encoding/json"
   "fmt"
   "net/url"
)

func Unauth() (*AuthData, error) {
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "gw.cds.amcn.com",
         Path:   "/auth-orchestration-id/api/v1/unauth",
      },
      map[string]string{
         "x-amcn-network":   "amcplus",
         "x-amcn-platform":  "web",
         "x-amcn-tenant":    "amcn",
         "x-amcn-device-id": "-",
         "x-amcn-language":  "en",
      },
      nil,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("unauth failed with status: %d", resp.StatusCode)
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

func SeriesDetail(authToken string, seriesId int) (*ContentNode, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "gw.cds.amcn.com",
         Path: fmt.Sprint(
            "/content-compiler-cr/api/v1/content/amcn/amcplus/type/series-detail/id/",
            seriesId,
         ),
      },
      map[string]string{
         "authorization":   "Bearer " + authToken,
         "x-amcn-network":  "amcplus",
         "x-amcn-platform": "android",
         "x-amcn-tenant":   "amcn",
      },
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("series detail failed with status: %d", resp.StatusCode)
   }
   // Internal envelope to strip the first layer
   var envelope struct {
      Success bool        `json:"success"`
      Status  int         `json:"status"`
      Data    ContentNode `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
      return nil, err
   }
   return &envelope.Data, nil
}

//go:embed playback.json
var playback_json []byte

func Playback(authToken string, videoId int) (*PlaybackResult, error) {
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "gw.cds.amcn.com",
         Path:   fmt.Sprint("/playback-id/api/v1/playback/", videoId),
      },
      map[string]string{
         "authorization":       "Bearer " + authToken,
         "content-type":        "application/json",
         "x-amcn-language":     "en",
         "x-amcn-network":      "amcplus",
         "x-amcn-platform":     "web",
         "x-amcn-service-id":   "amcplus",
         "x-amcn-tenant":       "amcn",
         "x-amcn-device-ad-id": "-",
         "x-ccpa-do-not-sell":  "doNotPassData",
      },
      playback_json,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
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
