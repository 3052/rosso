// play.go
package kanopy

import (
   "encoding/json"
   "errors"
   "io"
   "net/url"

   "41.neocities.org/maya"
)

type PlayRequest struct {
   DomainId int `json:"domainId"`
   UserId   int `json:"userId"`
   VideoId  int `json:"videoId"`
}

type PlayResponse struct {
   PlayId    string     `json:"playId"`
   Manifests []Manifest `json:"manifests"`
   Captions  []Caption  `json:"captions"`
}

type Caption struct {
   Language string        `json:"language"`
   Files    []CaptionFile `json:"files"`
   Label    string        `json:"label"`
}

type CaptionFile struct {
   Type string `json:"type"`
   Url  string `json:"url"`
}

type Manifest struct {
   ManifestType   string    `json:"manifestType"`
   Url            string    `json:"url"`
   DrmType        string    `json:"drmType"`
   StudioDrm      StudioDrm `json:"studioDrm"`
   StorageService string    `json:"storageService"`
   Cdn            string    `json:"cdn"`
   DrmLicenseId   string    `json:"drmLicenseID"`
}

type StudioDrm struct {
   AuthXml      string `json:"authXml"`
   DrmLicenseId string `json:"drmLicenseId"`
}

func (resp *PlayResponse) DashManifest() (*Manifest, error) {
   for _, manifest := range resp.Manifests {
      if manifest.ManifestType == "dash" {
         return &manifest, nil
      }
   }
   return nil, errors.New("dash manifest not found")
}

func CreatePlay(domainId int, userId int, videoId int, jwt string) (*PlayResponse, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/plays",
   }

   headers := map[string]string{
      "content-type":  "application/json",
      "user-agent":    "!",
      "x-version":     "!/!/!/!",
      "authorization": "Bearer " + jwt,
   }

   reqBody := PlayRequest{
      DomainId: domainId,
      UserId:   userId,
      VideoId:  videoId,
   }

   bodyBytes, err := json.Marshal(reqBody)
   if err != nil {
      return nil, err
   }

   resp, err := maya.Post(target, headers, bodyBytes)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   respBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var playResp PlayResponse
   if err := json.Unmarshal(respBytes, &playResp); err != nil {
      return nil, err
   }

   return &playResp, nil
}
