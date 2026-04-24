package kanopy

import (
   "errors"
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
   var result Video
   identifier := path.Base(url_parse.Path)
   numeric_id, err := strconv.Atoi(identifier)
   if err != nil {
      result.Alias = identifier
   } else {
      result.VideoId = numeric_id
   }
   return &result, nil
}

func (m *Manifest) GetManifest() (*url.URL, error) {
   return url.Parse(m.Url)
}
