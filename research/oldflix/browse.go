package oldflix

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
   "strings"
)

// GetOriginalTrack searches the available tracks for the one labeled
// "Original"
func (b *BrowsePlayResponse) GetOriginalTrack() (*Track, error) {
   for _, track := range b.Movie.Tracks {
      // Using EqualFold to safely match "Original", "original", etc.
      if strings.EqualFold(track.Lang, "Original") {
         return &track, nil
      }
   }
   return nil, fmt.Errorf("track with language 'Original' not found")
}

type Track struct {
   ID   string `json:"id"`
   Lang string `json:"lang"`
   Lnk  string `json:"lnk"`
}

type BrowsePlayResponse struct {
   ID    string `json:"id"`
   Movie struct {
      ID     string  `json:"id"`
      Tracks []Track `json:"tracks"`
   } `json:"movie"`
}

// https://oldflix.com.br/browse/play/5d5d54a4d55dc050f8468513
// BrowsePlay retrieves internal streaming parameters required to unlock the
// M3U8 payload
func (c *Client) BrowsePlay(contentID string) (*BrowsePlayResponse, error) {
   data := url.Values{}
   data.Set("id", contentID)

   req, err := http.NewRequest("POST", BaseURL+"/api/browse/play", strings.NewReader(data.Encode()))
   if err != nil {
      return nil, err
   }

   req.Header.Set("Authorization", "Bearer "+c.Token)
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   req.Header.Set("User-Agent", "okhttp/4.12.0")

   resp, err := c.HTTPClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var browseResp BrowsePlayResponse
   if err := json.NewDecoder(resp.Body).Decode(&browseResp); err != nil {
      return nil, fmt.Errorf("failed to decode browse play response: %w", err)
   }

   return &browseResp, nil
}
