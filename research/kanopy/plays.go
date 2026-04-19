// plays.go
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
   ManifestType   string `json:"manifestType"`
   URL            string `json:"url"`
   DrmType        string `json:"drmType"`
   StorageService string `json:"storageService"`
   CDN            string `json:"cdn"`
   DrmLicenseID   string `json:"drmLicenseID"`
}

type PlayResponse struct {
   PlayID    string     `json:"playId"`
   Manifests []Manifest `json:"manifests"`
}

func Plays(jwt string, domainId int, userId int, videoId int) (*PlayResponse, error) {
   reqData := PlayRequest{
      DomainID: domainId,
      UserID:   userId,
      VideoID:  videoId,
   }
   body, err := json.Marshal(reqData)
   if err != nil {
      return nil, err
   }

   targetUrl, err := url.Parse("https://www.kanopy.com/kapi/plays")
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "content-type":  "application/json",
      "authorization": "Bearer " + jwt,
      "x-version":     "!/!/!/!",
   }

   resp, err := maya.Post(targetUrl, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var playResp PlayResponse
   err = json.NewDecoder(resp.Body).Decode(&playResp)
   if err != nil {
      return nil, err
   }

   return &playResp, nil
}
