package kanopy

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
)

type VideoResponse struct {
   Type  string `json:"type"`
   Video struct {
      VideoID         int    `json:"videoId"`
      Title           string `json:"title"`
      DescriptionHTML string `json:"descriptionHtml"`
      ProductionYear  int    `json:"productionYear"`
      DurationSeconds int    `json:"durationSeconds"`
      Supplier        struct {
         TermID int    `json:"termId"`
         Name   string `json:"name"`
      } `json:"supplier"`
   } `json:"video"`
}

// GetVideo fetches video metadata using the movie alias.
func (s *Session) GetVideo(alias string) (*VideoResponse, error) {
   url := fmt.Sprintf("%s/kapi/videos/alias/%s", BaseURL, alias)

   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("X-Version", XVersion)
   req.Header.Set("Authorization", "Bearer "+s.JWT)
   req.Header.Set("User-Agent", "Go-http-client/2.0")

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("get video failed with status: %d", resp.StatusCode)
   }

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var videoResp VideoResponse
   if err := json.Unmarshal(respBody, &videoResp); err != nil {
      return nil, err
   }

   return &videoResp, nil
}
