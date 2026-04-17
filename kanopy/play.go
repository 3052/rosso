package kanopy

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

type PlayRequest struct {
   DomainId int `json:"domainId"`
   UserId   int `json:"userId"`
   VideoId  int `json:"videoId"`
}

type CaptionFile struct {
   Type string `json:"type"`
   Url  string `json:"url"`
}

type Caption struct {
   Files    []CaptionFile `json:"files"`
   Label    string        `json:"label"`
   Language string        `json:"language"`
}

type Dva struct {
   U int `json:"u"`
}

type StudioDrm struct {
   AuthXml      string `json:"authXml"`
   DrmLicenseId string `json:"drmLicenseId"`
}

func (m *Manifest) GetManifest() (*url.URL, error) {
   return url.Parse(m.Url)
}

type Manifest struct {
   Cdn            string    `json:"cdn"`
   DrmLicenseId   string    `json:"drmLicenseID"`
   DrmType        string    `json:"drmType"`
   ManifestType   string    `json:"manifestType"`
   StorageService string    `json:"storageService"`
   StudioDrm      StudioDrm `json:"studioDrm"`
   Url            string    `json:"url"`
}

type PlayResponse struct {
   Captions  []Caption   `json:"captions"`
   Dva       Dva         `json:"dva"`
   Manifests []*Manifest `json:"manifests"`
   PlayId    string      `json:"playId"`
}

// DashManifest returns the manifest with type "dash" or an error if it is not found.
func (p *PlayResponse) DashManifest() (*Manifest, error) {
   for _, manifest := range p.Manifests {
      if manifest.ManifestType == "dash" {
         return manifest, nil
      }
   }
   return nil, fmt.Errorf("dash manifest not found in play response")
}

// CreatePlay registers a playback event using the DomainId from a Membership
// and the VideoId from a Video.
func (s *Session) CreatePlay(membership *Membership, video *Video) (*PlayResponse, error) {
   if membership == nil {
      return nil, fmt.Errorf("membership context is required to create a play")
   }
   if video == nil {
      return nil, fmt.Errorf("video context is required to create a play")
   }

   payload := PlayRequest{
      DomainId: membership.DomainId,
      UserId:   s.UserId,
      VideoId:  video.VideoId,
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   req, err := http.NewRequest("POST", BaseUrl+"/kapi/plays", bytes.NewBuffer(body))
   if err != nil {
      return nil, err
   }

   req.Header.Set("X-Version", Xversion)
   req.Header.Set("Authorization", "Bearer "+s.Jwt)
   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("User-Agent", UserAgent)

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("create play failed with status: %d", resp.StatusCode)
   }

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var playResp PlayResponse
   if err := json.Unmarshal(respBody, &playResp); err != nil {
      return nil, err
   }

   return &playResp, nil
}
