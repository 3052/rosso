package crave

import (
   "bytes"
   "encoding/base64"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "strings"
   "time"
)

// InitPlayback orchestrates the entire flow from a public URL to getting a playback session.
func (c *Client) InitPlayback(publicURL string) (*PlaybackSession, error) {
   mediaID, err := extractMediaID(publicURL)
   if err != nil {
      return nil, fmt.Errorf("failed to extract media ID: %w", err)
   }

   contentID, err := c.GetContentID(mediaID)
   if err != nil {
      return nil, fmt.Errorf("failed to get content ID: %w", err)
   }

   pkgID, destID, err := c.GetPlaybackDetails(contentID)
   if err != nil {
      return nil, fmt.Errorf("failed to get playback details: %w", err)
   }

   manifest, err := c.GetManifest(contentID, pkgID, destID)
   if err != nil {
      return nil, fmt.Errorf("failed to get manifest: %w", err)
   }

   return &PlaybackSession{
      ContentID:        contentID,
      ContentPackageID: pkgID,
      DestinationID:    destID,
      ManifestURL:      manifest,
   }, nil
}

// GetWidevineLicense issues the DRM license request using the provided payload and the session details
func (c *Client) GetWidevineLicense(session *PlaybackSession, payload string) ([]byte, error) {
   // The API expects the contentId as an integer
   contentIDInt, err := strconv.Atoi(session.ContentID)
   if err != nil {
      return nil, fmt.Errorf("failed to parse content ID to int: %w", err)
   }

   reqBody := WidevineRequest{
      Payload: payload,
      PlaybackContext: PlaybackContext{
         ContentID:        contentIDInt,
         ContentPackageID: session.ContentPackageID,
         PlatformID:       1, // Hardcoded to 1 for Web
         DestinationID:    session.DestinationID,
         GL:               "0",
         JWT:              c.jwtToken,
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

   resp, err := c.httpClient.Do(req)
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

// PlaybackSession holds the necessary IDs to make subsequent requests (like licensing)
type PlaybackSession struct {
   ContentID        string
   ContentPackageID int
   DestinationID    int
   ManifestURL      string
}

// WidevineRequest represents the JSON body needed for the DRM license request
type WidevineRequest struct {
   Payload         string          `json:"payload"`
   PlaybackContext PlaybackContext `json:"playbackContext"`
}

type PlaybackContext struct {
   ContentID        int    `json:"contentId"`
   ContentPackageID int    `json:"contentpackageId"` // Note: lower-case 'p' as per their API
   PlatformID       int    `json:"platformId"`
   DestinationID    int    `json:"destinationId"`
   GL               string `json:"gl"`
   JWT              string `json:"jwt"`
}

const licenseURL  = "https://license.9c9media.com/widevine"

