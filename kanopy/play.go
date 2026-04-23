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

type PlayResponse struct {
   PlayId    string     `json:"playId"`
   Manifests []Manifest `json:"manifests"`
}

func CreatePlay(login *LoginResponse, membershipData *Membership, video *VideoResponse) (*PlayResponse, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/plays",
   }

   payload := PlayRequest{
      DomainId: membershipData.DomainId,
      UserId:   login.UserId,
      VideoId:  video.Video.VideoId,
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

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var play PlayResponse
   if err := json.Unmarshal(respBody, &play); err != nil {
      return nil, err
   }

   return &play, nil
}
