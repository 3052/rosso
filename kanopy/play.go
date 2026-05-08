package kanopy

import (
   "encoding/json"
   "errors"
   "net/url"

   "41.neocities.org/maya"
)

func (play *PlayResponse) GetDash() (*Manifest, error) {
   for _, manifest_data := range play.Manifests {
      if manifest_data.ManifestType == "dash" {
         return &manifest_data, nil
      }
   }
   return nil, errors.New("dash manifest not found")
}

type PlayResponse struct {
   PlayId    string     `json:"playId"`
   Manifests []Manifest `json:"manifests"`
   Captions  []Caption  `json:"captions"`
}

func (m *Manifest) GetUrl() (*url.URL, error) {
   return url.Parse(m.Url)
}

type Manifest struct {
   Url            string `json:"url"`
   ManifestType   string `json:"manifestType"`
   DrmType        string `json:"drmType"`
   StorageService string `json:"storageService"`
   Cdn            string `json:"cdn"`
   DrmLicenseId   string `json:"drmLicenseID"`
}

type PlayRequest struct {
   DomainId int `json:"domainId"`
   UserId   int `json:"userId"`
   VideoId  int `json:"videoId"`
}

type StudioDrm struct {
   AuthXml      string `json:"authXml"`
   DrmLicenseId string `json:"drmLicenseId"`
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
