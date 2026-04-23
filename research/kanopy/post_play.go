package kanopy

import (
   "encoding/json"
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
   DrmLicenseID string `json:"drmLicenseID"`
}

type Play struct {
   PlayId    string     `json:"playId"`
   Manifests []Manifest `json:"manifests"`
}

func PostPlay(loginData *Login, membershipData *Membership, videoData *Video) (*Play, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/plays",
   }

   payload := PlayRequest{
      DomainId: membershipData.DomainId,
      UserId:   loginData.UserId,
      VideoId:  videoData.VideoId,
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "authorization": "Bearer " + loginData.Jwt,
      "x-version":     "!/!/!/!",
   }

   resp, err := maya.Post(endpoint, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var playData Play
   if err := json.NewDecoder(resp.Body).Decode(&playData); err != nil {
      return nil, err
   }

   return &playData, nil
}
