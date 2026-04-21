package kanopy

import (
   "encoding/json"
   "errors"
   "net/url"

   "41.neocities.org/maya"
)

func (m *Manifest) GetManifest() (*url.URL, error) {
   return url.Parse(m.URL)
}

type StudioDRM struct {
   AuthXML      string `json:"authXml"`
   DRMLicenseID string `json:"drmLicenseId"`
}

type Manifest struct {
   ManifestType   string     `json:"manifestType"`
   URL            string     `json:"url"`
   DRMType        string     `json:"drmType"`
   StudioDRM      *StudioDRM `json:"studioDrm"`
   StorageService string     `json:"storageService"`
   CDN            string     `json:"cdn"`
   DRMLicenseID   string     `json:"drmLicenseID"`
}

type PlayResponse struct {
   PlayID    string     `json:"playId"`
   Manifests []Manifest `json:"manifests"`
}

func (pr *PlayResponse) DashManifest() (*Manifest, error) {
   for _, m := range pr.Manifests {
      if m.ManifestType == "dash" {
         return &m, nil
      }
   }
   return nil, errors.New("dash manifest not found")
}

func CreatePlay(loginResp *LoginResponse, domainID, videoID int) (*PlayResponse, error) {
   playsURL := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/plays",
   }

   payload := map[string]int{
      "domainId": domainID,
      "userId":   loginResp.UserID,
      "videoId":  videoID,
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "content-type":  "application/json",
      "user-agent":    "!",
      "x-version":     "!/!/!/!",
      "authorization": "Bearer " + loginResp.JWT,
   }

   resp, err := maya.Post(playsURL, headers, body)
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
