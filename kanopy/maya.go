package kanopy

import (
   "41.neocities.org/maya"
   "encoding/json"
   "fmt"
   "io"
   "net/url"
)

type Manifest struct {
   Cdn            string    `json:"cdn"`
   DrmLicenseId   string    `json:"drmLicenseID"`
   DrmType        string    `json:"drmType"`
   ManifestType   string    `json:"manifestType"`
   StorageService string    `json:"storageService"`
   StudioDrm      StudioDrm `json:"studioDrm"`
   Url            string    `json:"url"`
}

func (s *Session) GetWidevine(manifest *Manifest, body []byte) ([]byte, error) {
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   host,
         Path:   "/kapi/licenses/widevine/" + manifest.DrmLicenseId,
      },
      map[string]string{
         "authorization": "Bearer " + s.Jwt,
         "user-agent":    UserAgent,
         "x-version":     Xversion,
      },
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("widevine license request failed with status: %d", resp.StatusCode)
   }

   return io.ReadAll(resp.Body)
}

// Login authenticates the user and returns an initialized Session.
func Login(email, password string) (*Session, error) {
   body, err := json.Marshal(map[string]any{
      "credentialType": "email",
      "emailUser": map[string]string{
         "email":    email,
         "password": password,
      },
   })
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   host,
         Path:   "/kapi/login",
      },
      map[string]string{
         "content-type": "application/json",
         "user-agent":   UserAgent,
      },
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("login failed with status: %d", resp.StatusCode)
   }

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var session Session
   if err := json.Unmarshal(respBody, &session); err != nil {
      return nil, err
   }

   return &session, nil
}

// GetMemberships fetches the library memberships associated with the session's UserId
// and returns the list of memberships directly.
func (s *Session) GetMemberships() ([]Membership, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme:   "https",
         Host:     host,
         Path:     "/kapi/memberships",
         RawQuery: fmt.Sprint("userId=", s.UserId),
      },
      map[string]string{
         "authorization": "Bearer " + s.Jwt,
         "user-agent":    UserAgent,
         "x-version":     Xversion,
      },
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("get memberships failed with status: %d", resp.StatusCode)
   }

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var wrapper struct {
      List []Membership `json:"list"`
   }

   if err := json.Unmarshal(respBody, &wrapper); err != nil {
      return nil, err
   }

   return wrapper.List, nil
}

// CreatePlay registers a playback event using the DomainId from a Membership
// and the VideoId from a Video.
func (s *Session) CreatePlay(membership *Membership, video *Video) (*PlayResponse, error) {
   body, err := json.Marshal(map[string]int{
      "domainId": membership.DomainId,
      "userId":   s.UserId,
      "videoId":  video.VideoId,
   })
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   host,
         Path:   "/kapi/plays",
      },
      map[string]string{
         "authorization": "Bearer " + s.Jwt,
         "content-type":  "application/json",
         "user-agent":    UserAgent,
         "x-version":     Xversion,
      },
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != 200 {
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

// GetVideo fetches video metadata and strips the outer wrapper to return the Video object.
func (s *Session) GetVideo(alias string) (*Video, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   host,
         Path:   "/kapi/videos/alias/" + alias,
      },
      map[string]string{
         "authorization": "Bearer " + s.Jwt,
         "x-version":     Xversion,
      },
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("get video failed with status: %d", resp.StatusCode)
   }

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var wrapper struct {
      Video Video `json:"video"`
   }

   if err := json.Unmarshal(respBody, &wrapper); err != nil {
      return nil, err
   }

   return &wrapper.Video, nil
}
