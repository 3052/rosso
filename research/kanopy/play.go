package kanopy

import (
   "encoding/json"
   "io"
   "net/url"

   "41.neocities.org/maya"
)

type PlaysResponse struct {
   PlayID    string     `json:"playId"`
   Manifests []Manifest `json:"manifests"`
   Captions  []Caption  `json:"captions"`
   DVA       DVA        `json:"dva"`
}

type Manifest struct {
   ManifestType   string    `json:"manifestType"`
   URL            string    `json:"url"`
   DRMType        string    `json:"drmType"`
   StudioDRM      StudioDRM `json:"studioDrm"`
   StorageService string    `json:"storageService"`
   CDN            string    `json:"cdn"`
   DRMLicenseID   string    `json:"drmLicenseID"`
}

type StudioDRM struct {
   AuthXML      string `json:"authXml"`
   DRMLicenseID string `json:"drmLicenseId"`
}

type Caption struct {
   Language string        `json:"language"`
   Files    []CaptionFile `json:"files"`
   Label    string        `json:"label"`
}

type CaptionFile struct {
   Type string `json:"type"`
   URL  string `json:"url"`
}

type DVA struct {
   U int `json:"u"`
}

type createPlayRequest struct {
   DomainID int64 `json:"domainId"`
   UserID   int64 `json:"userId"`
   VideoID  int   `json:"videoId"`
}

func CreatePlay(jwt string, membership *Membership, videoID int) (*PlaysResponse, error) {
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

   reqBody := createPlayRequest{
      DomainID: membership.DomainID,
      UserID:   membership.UserID,
      VideoID:  videoID,
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

   var playsResp PlaysResponse
   if err := json.Unmarshal(respBytes, &playsResp); err != nil {
      return nil, err
   }

   return &playsResp, nil
}
