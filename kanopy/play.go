package kanopy

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
)

type PlayRequest struct {
   DomainId int `json:"domainId"`
   UserId   int `json:"userId"`
   VideoId  int `json:"videoId"`
}

type Manifest struct {
   ManifestType string `json:"manifestType"`
   Url          string `json:"url"`
   DrmType      string `json:"drmType"`
   DrmLicenseId string `json:"drmLicenseID"`
   StudioDrm    struct {
      AuthXml      string `json:"authXml"`
      DrmLicenseId string `json:"drmLicenseId"`
   } `json:"studioDrm"`
}

type PlayResponse struct {
   PlayId    string      `json:"playId"`
   Manifests []*Manifest `json:"manifests"`
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
