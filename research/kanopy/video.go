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

type Term struct {
   TermID int    `json:"termId"`
   Name   string `json:"name"`
}

type ImageSizes struct {
   Small  string `json:"small"`
   Medium string `json:"medium"`
   Large  string `json:"large"`
}

type AgeRating struct {
   AgeRatingCountry string `json:"ageRatingCountry"`
   AgeRating        string `json:"ageRating"`
   AgeRatingIconSvg string `json:"ageRatingIconSvg"`
   AgeRatingIconPng string `json:"ageRatingIconPng"`
   AgeRatingIconJpg string `json:"ageRatingIconJpg"`
   AgeRatingAria    string `json:"ageRatingAria"`
}

type TrackLabel struct {
   Language string `json:"language"`
   Label    string `json:"label"`
   Type     string `json:"type"`
}

// Video represents the core video metadata object.
type Video struct {
   VideoId         int    `json:"videoId"`
   Title           string `json:"title"`
   DescriptionHTML string `json:"descriptionHtml"`
   Images          struct {
      Landscapes ImageSizes `json:"landscapes"`
      Posters    ImageSizes `json:"posters"`
   } `json:"images"`
   HasBurntInCaptions bool     `json:"hasBurntInCaptions"`
   HasCaptions        bool     `json:"hasCaptions"`
   CaptionLanguages   []string `json:"captionLanguages"`
   ProductionYear     int      `json:"productionYear"`
   Taxonomies         struct {
      Subjects   []Term `json:"subjects"`
      Tags       []Term `json:"tags"`
      Filmmakers []Term `json:"filmmakers"`
      Cast       []Term `json:"cast"`
      Languages  []Term `json:"languages"`
      Supplier   Term   `json:"supplier"`
   } `json:"taxonomies"`
   StarRating struct {
      Count   int `json:"count"`
      Average int `json:"average"`
   } `json:"starRating"`
   IsKids                     bool                 `json:"isKids"`
   DurationSeconds            int                  `json:"durationSeconds"`
   HasPublicPerformanceRights bool                 `json:"hasPublicPerformanceRights"`
   AgeClassificationByCountry map[string]AgeRating `json:"ageClassificationByCountry"`
   IsFree                     bool                 `json:"isFree"`
   IsRequestable              bool                 `json:"isRequestable"`
   AncestorVideoIds           []int                `json:"ancestorVideoIds"`
   Alias                      string               `json:"alias"`
   FeedID                     int                  `json:"feedId"`
   NeedsTitleTreatment        bool                 `json:"needsTitleTreatment"`
   IsSilent                   bool                 `json:"isSilent"`
   CaptionLabels              []TrackLabel         `json:"captionLabels"`
   AudioTrackLabels           []TrackLabel         `json:"audioTrackLabels"`
   SubtextLabels              []TrackLabel         `json:"subtextLabels"`
}

// GetVideo fetches comprehensive video metadata and returns the inner Video object directly.
func (s *Session) GetVideo(alias string) (*Video, error) {
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

   // Use an anonymous struct to strip the root wrapper
   var wrapper struct {
      Video Video `json:"video"`
   }

   if err := json.Unmarshal(respBody, &wrapper); err != nil {
      return nil, err
   }

   return &wrapper.Video, nil
}
