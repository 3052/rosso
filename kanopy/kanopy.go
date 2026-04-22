package kanopy

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "io"
   "net/url"
   "path"
   "strconv"
   "strings"
)

type Video struct {
   VideoId         int    `json:"videoId"`
   Title           string `json:"title"`
   DescriptionHtml string `json:"descriptionHtml"`
   ProductionYear  int    `json:"productionYear"`
   IsKids          bool   `json:"isKids"`
   DurationSeconds int    `json:"durationSeconds"`
   IsFree          bool   `json:"isFree"`
   Alias           string `json:"alias"`
}

func GetVideo(alias string, jwt string) (*Video, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/videos/alias/" + alias,
   }

   headers := map[string]string{
      "x-version":     "!/!/!/!",
      "authorization": "Bearer " + jwt,
   }

   resp, err := maya.Get(target, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   respBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   var result struct {
      Type  string `json:"type"`
      Video Video  `json:"video"`
   }
   if err := json.Unmarshal(respBytes, &result); err != nil {
      return nil, err
   }
   return &result.Video, nil
}

func GetLicense(drmLicenseId string, body []byte, jwt string) ([]byte, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/licenses/widevine/" + drmLicenseId,
   }

   headers := map[string]string{
      "authorization": "Bearer " + jwt,
      "user-agent":    "!",
      "x-version":     "!/!/!/!",
   }

   resp, err := maya.Post(target, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   respBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   return respBytes, nil
}

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
   numeric_id, err := strconv.Atoi(identifier)
   if err != nil {
      v.Alias = identifier
   } else {
      v.VideoId = numeric_id
   }
   return v, nil
}

func (m *Manifest) GetManifest() (*url.URL, error) {
   return url.Parse(m.Url)
}
