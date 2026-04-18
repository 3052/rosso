package oldflix

import (
   "encoding/json"
   "errors"
   "fmt"
   "net/http"
   "net/url"
   "strings"
)

func (b *Browse) FetchWatch(trackId, token string) (*Watch, error) {
   data := url.Values{
      "id": {b.Id},
      "m":  {b.Movie.Id},
      "tk": {trackId}, // tk is the audio/language track id
   }.Encode()
   req, err := http.NewRequest(
      "POST", BaseUrl+"/api/watch/play", strings.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Authorization", "Bearer "+token)
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var result Watch
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, fmt.Errorf("failed to decode watch play response: %w", err)
   }
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   return &result, nil
}
