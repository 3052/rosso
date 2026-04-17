package kanopy

import (
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "path"
   "strconv"
   "strings"
)

// Supports URLs such as:
// - https://kanopy.com/video/6440418
// - https://kanopy.com/video/genius-party
// - https://kanopy.com/en/video/genius-party
// - https://kanopy.com/en/product/genius-party
func ParseVideo(urlData string) (*Video, error) {
   url_parse, err := url.Parse(urlData)
   if err != nil {
      return nil, err
   }
   if !strings.Contains(url_parse.Host, "kanopy.com") {
      return nil, errors.New("invalid domain")
   }
   // Get the directory of the path (removes the final identifier).
   // e.g., "/en/product/genius-party" -> "/en/product"
   dir := path.Dir(url_parse.Path)
   // Check if the directory ends with "/video" OR "/product".
   // This supports:
   // - /video/{id}
   // - /en/video/{id}
   // - /en/product/{id}
   if !strings.HasSuffix(dir, "/video") && !strings.HasSuffix(dir, "/product") {
      return nil, errors.New("invalid path structure")
   }
   v := &Video{}
   identifier := path.Base(url_parse.Path)
   numericId, err := strconv.Atoi(identifier)
   if err != nil {
      v.Alias = identifier
   } else {
      v.VideoId = numericId
   }
   return v, nil
}

// Video represents the flattened video metadata, omitting truncated nested objects.
type Video struct {
   Alias                      string `json:"alias"`
   AncestorVideoIds           []int  `json:"ancestorVideoIds"`
   DescriptionHtml            string `json:"descriptionHtml"`
   DurationSeconds            int    `json:"durationSeconds"`
   FeedId                     int    `json:"feedId"`
   HasBurntInCaptions         bool   `json:"hasBurntInCaptions"`
   HasCaptions                bool   `json:"hasCaptions"`
   HasPublicPerformanceRights bool   `json:"hasPublicPerformanceRights"`
   IsFree                     bool   `json:"isFree"`
   IsKids                     bool   `json:"isKids"`
   IsRequestable              bool   `json:"isRequestable"`
   IsSilent                   bool   `json:"isSilent"`
   NeedsTitleTreatment        bool   `json:"needsTitleTreatment"`
   ProductionYear             int    `json:"productionYear"`
   Title                      string `json:"title"`
   VideoId                    int    `json:"videoId"`
}

// GetVideo fetches video metadata and strips the outer wrapper to return the Video object.
func (s *Session) GetVideo(alias string) (*Video, error) {
   url := fmt.Sprintf("%s/kapi/videos/alias/%s", BaseUrl, alias)

   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("X-Version", Xversion)
   req.Header.Set("Authorization", "Bearer "+s.Jwt)

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

   var wrapper struct {
      Video Video `json:"video"`
   }

   if err := json.Unmarshal(respBody, &wrapper); err != nil {
      return nil, err
   }

   return &wrapper.Video, nil
}
