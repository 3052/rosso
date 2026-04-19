// play.go
package kanopy

import (
   "encoding/json"
   "io"
   "net/url"

   "41.neocities.org/maya"
)

type PlayRequest struct {
   DomainId int `json:"domainId"`
   UserId   int `json:"userId"`
   VideoId  int `json:"videoId"`
}

type Manifest struct {
   ManifestType string `json:"manifestType"`
   Url          string `json:"url"`
   DrmType      string `json:"drmType"`
   DrmLicenseID string `json:"drmLicenseID"`
}

type PlayResponse struct {
   PlayId    string     `json:"playId"`
   Manifests []Manifest `json:"manifests"`
}

func CreatePlay(req *PlayRequest, authorization string) (*PlayResponse, error) {
   reqBody, err := json.Marshal(req)
   if err != nil {
      return nil, err
   }

   targetUrl, err := url.Parse("https://www.kanopy.com/kapi/plays")
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "content-type":  "application/json",
      "authorization": authorization,
   }

   resp, err := maya.Post(targetUrl, headers, reqBody)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var playResp PlayResponse
   if err := json.Unmarshal(respBody, &playResp); err != nil {
      return nil, err
   }

   return &playResp, nil
}
