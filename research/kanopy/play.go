// file: plays.go
package kanopy

import (
   "encoding/json"
   "io"
   "net/url"

   "41.neocities.org/maya"
)

type PlayRequest struct {
   DomainID int `json:"domainId"`
   UserID   int `json:"userId"`
   VideoID  int `json:"videoId"`
}

type Manifest struct {
   ManifestType   string `json:"manifestType"`
   URL            string `json:"url"`
   DRMType        string `json:"drmType"`
   StorageService string `json:"storageService"`
   CDN            string `json:"cdn"`
   DRMLicenseID   string `json:"drmLicenseID"`
}

type PlayResponse struct {
   PlayID    string     `json:"playId"`
   Manifests []Manifest `json:"manifests"`
}

func (m *Membership) CreatePlay(jwt string, videoId int) (*PlayResponse, error) {
   targetUrl := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/plays",
   }

   reqBody := PlayRequest{
      DomainID: m.DomainID,
      UserID:   m.UserID,
      VideoID:  videoId,
   }

   jsonData, err := json.Marshal(reqBody)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "content-type":  "application/json",
      "user-agent":    "!",
      "x-version":     "!/!/!/!",
      "authorization": "Bearer " + jwt,
   }

   resp, err := maya.Post(targetUrl, headers, jsonData)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   bodyBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var playResp PlayResponse
   if err := json.Unmarshal(bodyBytes, &playResp); err != nil {
      return nil, err
   }

   return &playResp, nil
}
