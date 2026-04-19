package kanopy

import (
   "41.neocities.org/maya"
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

func (s *Session) GetWidevine(manifest *Manifest, challenge []byte) ([]byte, error) {
   if manifest == nil {
      return nil, fmt.Errorf("a valid stream manifest is required to request a DRM license")
   }
   if manifest.DrmLicenseId == "" {
      return nil, fmt.Errorf("manifest does not contain a DRM license ID")
   }
   const BaseUrl = "https://www.kanopy.com"
   url := fmt.Sprintf("%s/kapi/licenses/widevine/%s", BaseUrl, manifest.DrmLicenseId)
   req, err := http.NewRequest("POST", url, bytes.NewBuffer(challenge))
   if err != nil {
      return nil, err
   }

   req.Header.Set("Authorization", "Bearer "+s.Jwt)
   req.Header.Set("User-Agent", UserAgent)
   req.Header.Set("X-Version", Xversion)

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
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
