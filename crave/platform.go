package crave

import (
   "bytes"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

// SL2000 max 1080p
func (c *ContentPackage) LicensePlayReady(contentId int, accessToken string, payload []byte) ([]byte, error) {
   return c.fetchLicense(contentId, accessToken, payload, 48, "playready")
}

// L3 max 720p
func (c *ContentPackage) LicenseWidevine(contentId int, accessToken string, payload []byte) ([]byte, error) {
   return c.fetchLicense(contentId, accessToken, payload, 1, "widevine")
}

func (c *ContentPackage) ManifestWidevine(contentId int, accessToken string) (*Manifest, error) {
   return c.fetchManifest(contentId, accessToken, 1)
}

func (c *ContentPackage) ManifestPlayReady(contentId int, accessToken string) (*Manifest, error) {
   return c.fetchManifest(contentId, accessToken, 48)
}

type Manifest struct {
   Message  string
   Playback string
}

type ContentPackage struct {
   Id            int
   DestinationId int
}

// --- Private Helpers ---

func (c *ContentPackage) fetchManifest(contentId int, accessToken string, platformId int) (*Manifest, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme: "https",
         Host:   "stream.video.9c9media.com",
         Path: fmt.Sprintf(
            "/meta/content/%v/contentpackage/%v/destination/%v/platform/%v",
            contentId, c.Id, c.DestinationId, platformId,
         ),
         RawQuery: url.Values{
            "filter": {"ff"}, // 1080p
            "format": {"mpd"},
            "hd":     {"true"}, // 1080p
            "mcv":    {"true"}, // H.264 + HEVC
            "uhd":    {"true"}, // HEVC
         }.Encode(),
      },
      Header: http.Header{},
   }
   req.Header.Set("authorization", "Bearer "+accessToken)

   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var result Manifest
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }

   return &result, nil
}

func (c *ContentPackage) fetchLicense(contentId int, accessToken string, payload []byte, platformId int, path string) ([]byte, error) {
   data, err := json.Marshal(map[string]any{
      "payload": payload,
      "playbackContext": map[string]any{
         "contentId":        contentId,
         "contentpackageId": c.Id, // lower-case 'p' as per their API
         "platformId":       platformId,
         "destinationId":    c.DestinationId,
         "jwt":              accessToken,
      },
   })
   if err != nil {
      return nil, err
   }

   req, err := http.NewRequest(
      "POST", "https://license.9c9media.com/"+path, bytes.NewBuffer(data),
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

   return data, nil
}
