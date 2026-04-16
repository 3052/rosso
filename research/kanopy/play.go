package kanopy

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
)

type PlayRequest struct {
   DomainID int `json:"domainId"`
   UserID   int `json:"userId"`
   VideoID  int `json:"videoId"`
}

type Manifest struct {
   ManifestType string `json:"manifestType"`
   URL          string `json:"url"`
   DrmType      string `json:"drmType"`
   DrmLicenseID string `json:"drmLicenseID"`
   StudioDrm    struct {
      AuthXML      string `json:"authXml"`
      DrmLicenseID string `json:"drmLicenseId"`
   } `json:"studioDrm"`
}

type PlayResponse struct {
   PlayID    string     `json:"playId"`
   Manifests []Manifest `json:"manifests"`
}

// CreatePlay registers a playback event using the DomainID from a Membership
// and the VideoID from a VideoResponse.
func (s *Session) CreatePlay(membership *Membership, video *VideoResponse) (*PlayResponse, error) {
   if membership == nil {
      return nil, fmt.Errorf("membership context is required to create a play")
   }
   if video == nil {
      return nil, fmt.Errorf("video context is required to create a play")
   }

   payload := PlayRequest{
      DomainID: membership.DomainID,
      UserID:   s.UserID,
      VideoID:  video.Video.VideoID,
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   req, err := http.NewRequest("POST", BaseURL+"/kapi/plays", bytes.NewBuffer(body))
   if err != nil {
      return nil, err
   }

   req.Header.Set("X-Version", XVersion)
   req.Header.Set("Authorization", "Bearer "+s.JWT)
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
