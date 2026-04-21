package kanopy

import (
   "encoding/json"
   "errors"
   "net/url"

   "41.neocities.org/maya"
)

type StudioDrm struct {
   AuthXml      string `json:"authXml"`
   DrmLicenseId string `json:"drmLicenseId"`
}

type Manifest struct {
   ManifestType   string     `json:"manifestType"`
   Url            string     `json:"url"`
   DrmType        string     `json:"drmType"`
   StudioDrm      *StudioDrm `json:"studioDrm"`
   StorageService string     `json:"storageService"`
   Cdn            string     `json:"cdn"`
   DrmLicenseId   string     `json:"drmLicenseID"`
}

type CaptionFile struct {
   Type string `json:"type"`
   Url  string `json:"url"`
}

type Caption struct {
   Language string        `json:"language"`
   Files    []CaptionFile `json:"files"`
   Label    string        `json:"label"`
}

type PlayResponse struct {
   PlayId    string     `json:"playId"`
   Manifests []Manifest `json:"manifests"`
   Captions  []Caption  `json:"captions"`
}

func (pr *PlayResponse) DashManifest() (*Manifest, error) {
   for _, m := range pr.Manifests {
      if m.ManifestType == "dash" {
         return &m, nil
      }
   }
   return nil, errors.New("dash manifest not found")
}

func CreatePlay(loginResp *LoginResponse, domainId, videoId int) (*PlayResponse, error) {
   playsUrl := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/plays",
   }

   payload := map[string]int{
      "domainId": domainId,
      "userId":   loginResp.UserId,
      "videoId":  videoId,
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "content-type":  "application/json",
      "user-agent":    "!",
      "x-version":     "!/!/!/!",
      "authorization": "Bearer " + loginResp.Jwt,
   }

   resp, err := maya.Post(playsUrl, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var playResp PlayResponse
   if err := json.NewDecoder(resp.Body).Decode(&playResp); err != nil {
      return nil, err
   }

   return &playResp, nil
}
