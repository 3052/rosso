package kanopy

import (
   "errors"
   "fmt"
   "net/url"
   "path"
   "strconv"
   "strings"
)

// Session represents an authenticated user context.
type Session struct {
   Jwt       string `json:"jwt"`
   UserId    int    `json:"userId"`
   UserRole  string `json:"userRole"`
   VisitorId string `json:"visitorId"`
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

type CaptionFile struct {
   Type string `json:"type"`
   Url  string `json:"url"`
}

type Caption struct {
   Files    []CaptionFile `json:"files"`
   Label    string        `json:"label"`
   Language string        `json:"language"`
}

type Dva struct {
   U int `json:"u"`
}

type StudioDrm struct {
   AuthXml      string `json:"authXml"`
   DrmLicenseId string `json:"drmLicenseId"`
}

func (m *Manifest) GetManifest() (*url.URL, error) {
   return url.Parse(m.Url)
}

type PlayResponse struct {
   Captions  []Caption   `json:"captions"`
   Dva       Dva         `json:"dva"`
   Manifests []*Manifest `json:"manifests"`
   PlayId    string      `json:"playId"`
}

// DashManifest returns the manifest with type "dash" or an error if it is not found.
func (p *PlayResponse) DashManifest() (*Manifest, error) {
   for _, manifest := range p.Manifests {
      if manifest.ManifestType == "dash" {
         return manifest, nil
      }
   }
   return nil, fmt.Errorf("dash manifest not found in play response")
}

type Membership struct {
   IdentityId         int    `json:"identityId"`
   DomainId           int    `json:"domainId"`
   UserId             int    `json:"userId"`
   Status             string `json:"status"`
   IsDefault          bool   `json:"isDefault"`
   Sitename           string `json:"sitename"`
   Subdomain          string `json:"subdomain"`
   TicketsAvailable   int    `json:"ticketsAvailable"`
   MaxTicketsPerMonth int    `json:"maxTicketsPerMonth"`
}

const (
   host      = "www.kanopy.com"
   UserAgent = "!"
   Xversion  = "!/!/!/!"
)
