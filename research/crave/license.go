package crave

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "strconv"
)

// PlaybackSession holds the necessary IDs to make subsequent requests (like licensing)
type PlaybackSession struct {
   ContentId        string
   ContentPackageId int
   DestinationId    int
}

// GetWidevineLicense issues the DRM license request using the provided payload and the session details
func (t *TokenResponse) GetWidevineLicense(session *PlaybackSession, payload string) ([]byte, error) {
   // The API expects the contentId as an integer
   contentIdInt, err := strconv.Atoi(session.ContentId)
   if err != nil {
      return nil, fmt.Errorf("failed to parse content ID to int: %w", err)
   }

   reqBody := WidevineRequest{
      Payload: payload,
      PlaybackContext: PlaybackContext{
         ContentId:        contentIdInt,
         ContentPackageId: session.ContentPackageId,
         PlatformId:       1, // Hardcoded to 1 for Web
         DestinationId:    session.DestinationId,
         GL:               "0",
         JWT: t.AccessToken,
      },
   }

   bodyBytes, err := json.Marshal(reqBody)
   if err != nil {
      return nil, err
   }

   req, err := http.NewRequest(http.MethodPost, licenseURL, bytes.NewBuffer(bodyBytes))
   if err != nil {
      return nil, err
   }

   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("Origin", "https://www.crave.ca")
   req.Header.Set("Referer", "https://www.crave.ca/")

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("license request failed with status: %d", resp.StatusCode)
   }

   // The response is usually a binary widevine license
   return io.ReadAll(resp.Body)
}

// WidevineRequest represents the JSON body needed for the DRM license request
type WidevineRequest struct {
   Payload         string          `json:"payload"`
   PlaybackContext PlaybackContext `json:"playbackContext"`
}

type PlaybackContext struct {
   ContentId        int    `json:"contentId"`
   ContentPackageId int    `json:"contentpackageId"` // Note: lower-case 'p' as per their API
   PlatformId       int    `json:"platformId"`
   DestinationId    int    `json:"destinationId"`
   GL               string `json:"gl"`
   JWT              string `json:"jwt"`
}

const licenseURL  = "https://license.9c9media.com/widevine"

