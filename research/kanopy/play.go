// play.go
package kanopy

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type PlayManifest struct {
   ManifestType string `json:"manifestType"`
   Url          string `json:"url"`
   DrmType      string `json:"drmType"`
   DrmLicenseID string `json:"drmLicenseID"`
}

type PlayResponse struct {
   PlayId    string         `json:"playId"`
   Manifests []PlayManifest `json:"manifests"`
}

func CreatePlay(session *Session, VideoId int) (*PlayResponse, error) {
   payload := map[string]int{
      "domainId": session.DomainId,
      "userId":   session.UserId,
      "videoId":  VideoId,
   }

   bodyBytes, marshalError := json.Marshal(payload)
   if marshalError != nil {
      return nil, marshalError
   }

   targetUrl, parseError := url.Parse("https://www.kanopy.com/kapi/plays")
   if parseError != nil {
      return nil, parseError
   }

   headers := map[string]string{
      "content-type":  "application/json",
      "authorization": "Bearer " + session.Authorization,
      "x-version":     "!/!/!/!",
      "user-agent":    "!",
   }

   resp, requestError := maya.Post(targetUrl, headers, bodyBytes)
   if requestError != nil {
      return nil, requestError
   }
   defer resp.Body.Close()

   var playResp PlayResponse
   decodeError := json.NewDecoder(resp.Body).Decode(&playResp)
   if decodeError != nil {
      return nil, decodeError
   }
   return &playResp, nil
}
