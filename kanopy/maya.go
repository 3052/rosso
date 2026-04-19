package kanopy

import (
   "41.neocities.org/maya"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

// GetMemberships fetches the library memberships associated with the session's UserId
// and returns the list of memberships directly.
func (s *Session) GetMemberships() ([]Membership, error) {
   const BaseUrl = "https://www.kanopy.com"
   url := fmt.Sprintf("%s/kapi/memberships?userId=%d", BaseUrl, s.UserId)

   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("User-Agent", UserAgent)
   req.Header.Set("X-Version", Xversion)
   req.Header.Set("Authorization", "Bearer "+s.Jwt)

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
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
