package crave

import (
   _ "embed"
   "encoding/json"
   "errors"
   "fmt"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

func (m *Manifest) GetManifest() (*url.URL, error) {
   return url.Parse(m.Playback)
}

type Manifest struct {
   Message  string
   Playback string
}

type Account struct {
   AccessToken  string `json:"access_token"`
   AccountId    string `json:"account_id"`
   RefreshToken string `json:"refresh_token"`
}

type Media struct {
   FirstContent struct {
      Id int `json:"id,string"`
   }
   Id int `json:"id,string"`
}

func (p *Profile) String() string {
   var data strings.Builder
   data.WriteString("nickname = ")
   data.WriteString(p.Nickname)
   if p.HasPin {
      data.WriteString("\nhas pin = true")
   } else {
      data.WriteString("\nhas pin = false")
   }
   if p.Master {
      data.WriteString("\nmaster = true")
   } else {
      data.WriteString("\nmaster = false")
   }
   data.WriteString("\nmaturity = ")
   data.WriteString(p.Maturity)
   data.WriteString("\nid = ")
   data.WriteString(p.Id)
   return data.String()
}

type Profile struct {
   Nickname string `json:"nickname"`
   HasPin   bool   `json:"hasPin"`
   Master   bool
   Maturity string
   Id       string `json:"id"`
}

func (s *Subscription) String() string {
   var data strings.Builder
   data.WriteString("display name = ")
   data.WriteString(s.Experience.DisplayName)
   data.WriteString("\nexpiration date = ")
   data.WriteString(s.ExpirationDate)
   return data.String()
}

type Subscription struct {
   Experience struct {
      DisplayName string
   }
   ExpirationDate string
}

var Language = "EN"

//go:embed GetShowpage.gql
var get_showpage string

/*
https://crave.ca/en/movie/anaconda-2025-59881
https://crave.ca/en/play/anaconda-2025-3300246
https://crave.ca/movie/anaconda-2025-59881
https://crave.ca/play/anaconda-2025-3300246
*/
func ParseMedia(rawUrl string) (*Media, error) {
   parsedUrl, err := url.Parse(rawUrl)
   if err != nil {
      return nil, err
   }
   // Split the path directly.
   // e.g., "/en/movie/anaconda-2025-59881" -> ["", "en", "movie", "anaconda-2025-59881"]
   parts := strings.Split(parsedUrl.Path, "/")
   // We need at least 3 parts: the empty string (before the first "/"), the type, and the slug
   if len(parts) < 3 {
      return nil, errors.New("invalid URL path format")
   }
   // Safely grab the last two segments
   lastPart := parts[len(parts)-1] // e.g., "anaconda-2025-59881"
   typePart := parts[len(parts)-2] // e.g., "movie" or "play"
   // Find the last dash to extract the ID
   dashIdx := strings.LastIndex(lastPart, "-")
   if dashIdx == -1 || dashIdx == len(lastPart)-1 {
      return nil, errors.New("no ID found at the end of the URL")
   }
   idStr := lastPart[dashIdx+1:]
   // Convert extracted string to integer
   id, err := strconv.Atoi(idStr)
   if err != nil {
      return nil, fmt.Errorf("invalid ID format: %w", err)
   }
   // Populate struct based on the type
   media := &Media{}
   switch typePart {
   case "movie":
      media.Id = id
   case "play":
      media.FirstContent.Id = id
   default:
      return nil, fmt.Errorf("unknown media type: %s", typePart)
   }
   return media, nil
}

type ContentPackage struct {
   DestinationId int
   Id            int
}

// SL2000 max 2160p
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
            "filter": {"ff"}, // 2160p HEVC
            "format": {"mpd"},
            "hd":     {"true"}, // 1080p H.264
            "mcv":    {"true"}, // H.264 + HEVC
            "uhd":    {"true"}, // 2160p HEVC
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
