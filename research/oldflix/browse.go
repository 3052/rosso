package oldflix

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
   "strings"
)

// https://oldflix.com.br/browse/play/5d5d54a4d55dc050f8468513
func (l *Login) FetchBrowse(contentId string) (*Browse, error) {
   data := url.Values{"id": {contentId}}.Encode()
   req, err := http.NewRequest(
      "POST", BaseUrl+"/api/browse/play", strings.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer " + l.Token)
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   req.Header.Set("user-agent", "okhttp/4.12.0")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var browseResp Browse
   if err := json.NewDecoder(resp.Body).Decode(&browseResp); err != nil {
      return nil, fmt.Errorf("failed to decode browse play response: %w", err)
   }
   return &browseResp, nil
}

// GetOriginalTrack searches the available tracks for the one labeled
// "Original"
func (b *Browse) GetOriginalTrack() (*Track, error) {
   for _, track_data := range b.Movie.Tracks {
      // Using EqualFold to safely match "Original", "original", etc.
      if strings.EqualFold(track_data.Lang, "Original") {
         return &track_data, nil
      }
   }
   return nil, fmt.Errorf("track with language 'Original' not found")
}

type Track struct {
   Id   string `json:"id"`
   Lang string `json:"lang"`
   Lnk  string `json:"lnk"`
}

type Browse struct {
   Id    string `json:"id"`
   Movie struct {
      Id     string  `json:"id"`
      Tracks []Track `json:"tracks"`
   } `json:"movie"`
}
