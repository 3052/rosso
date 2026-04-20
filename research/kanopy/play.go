// File: create_play.go
package kanopy

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type PlayRequest struct {
   DomainID int `json:"domainId"`
   UserID   int `json:"userId"`
   VideoID  int `json:"videoId"`
}

type Manifest struct {
   ManifestType string `json:"manifestType"`
   URL          string `json:"url"`
   DrmType      string `json:"drmType"`
   DrmLicenseID string `json:"drmLicenseID"`
   StudioDrm    struct {
      AuthXML      string `json:"authXml"`
      DrmLicenseID string `json:"drmLicenseId"`
   } `json:"studioDrm"`
}

type PlayResponse struct {
   PlayID    string     `json:"playId"`
   Manifests []Manifest `json:"manifests"`
}

func CreatePlay(playReq *PlayRequest, token string) (*PlayResponse, error) {
   reqURL, err := url.Parse("https://www.kanopy.com/kapi/plays")
   if err != nil {
      return nil, err
   }

   payload, err := json.Marshal(playReq)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "content-type":  "application/json",
      "authorization": "Bearer " + token,
      "user-agent":    "!",
      "x-version":     "!/!/!/!",
   }

   resp, err := maya.Post(reqURL, headers, payload)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var result PlayResponse
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   return &result, nil
}
