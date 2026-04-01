package crave

import (
   "bytes"
   _ "embed"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "strconv"
)

// GetWidevineLicense issues the DRM license request using the provided payload
// and the session details
func (a *Account) GetWidevineLicense(session *PlaybackSession, payload string) ([]byte, error) {
   // The API expects the contentId as an integer
   contentIdInt, err := strconv.Atoi(session.ContentId)
   if err != nil {
      return nil, fmt.Errorf("failed to parse content ID to int: %w", err)
   }
   data, err := json.Marshal(WidevineRequest{
      Payload: payload,
      PlaybackContext: PlaybackContext{
         ContentId:        contentIdInt,
         ContentPackageId: session.ContentPackageId,
         PlatformId:       1, // Hardcoded to 1 for Web
         DestinationId:    session.DestinationId,
         Jwt:              a.AccessToken,
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      http.MethodPost, "https://license.9c9media.com/widevine",
      bytes.NewBuffer(data),
   )
   if err != nil {
      return nil, err
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      var result struct {
         Message string
      }
      err = json.Unmarshal(data, &result)
      if err != nil {
         return nil, err
      }
      return nil, errors.New(result.Message)
   }
   // The response is usually a binary widevine license
   return data, nil
}
