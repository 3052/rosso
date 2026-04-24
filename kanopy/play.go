package kanopy

import (
   "encoding/json"
   "errors"
   "net/url"

   "41.neocities.org/maya"
)

type PlayRequest struct {
   DomainId int `json:"domainId"`
   UserId   int `json:"userId"`
   VideoId  int `json:"videoId"`
}

type StudioDrm struct {
   AuthXml      string `json:"authXml"`
   DrmLicenseId string `json:"drmLicenseId"`
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

type File struct {
   Type string `json:"type"`
   Url  string `json:"url"`
}

type Caption struct {
   Language string `json:"language"`
   Files    []File `json:"files"`
   Label    string `json:"label"`
}

type PlayResponse struct {
   PlayId    string     `json:"playId"`
   Manifests []Manifest `json:"manifests"`
   Captions  []Caption  `json:"captions"`
}

func CreatePlay(login *LoginResponse, membershipData *Membership, videoData *Video) (*PlayResponse, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/plays",
   }

   payload := PlayRequest{
      DomainId: membershipData.DomainId,
      UserId:   login.UserId,
      VideoId:  videoData.VideoId,
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "authorization": "Bearer " + login.Jwt,
   }

   resp, err := maya.Post(endpoint, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var play PlayResponse
   if err := json.NewDecoder(resp.Body).Decode(&play); err != nil {
      return nil, err
   }

   return &play, nil
}

func (play *PlayResponse) GetDashManifest() (*Manifest, error) {
   for _, manifestData := range play.Manifests {
      if manifestData.ManifestType == "dash" {
         return &manifestData, nil
      }
   }
   return nil, errors.New("dash manifest not found")
}
