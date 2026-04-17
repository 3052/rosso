// request_cms_content.go
package tubi

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
)

// Content represents the exact structure of the provided content.json, omitting interface and null fields.
type Content struct {
   Actors               []string        `json:"actors"`
   AdLanguages          []string        `json:"ad_languages"`
   AvailabilityDuration int             `json:"availability_duration"`
   AvailabilityEnds     string          `json:"availability_ends"`
   AvailabilityStarts   string          `json:"availability_starts"`
   Backgrounds          []string        `json:"backgrounds"`
   CanonicalID          string          `json:"canonical_id"`
   ContentOrientation   string          `json:"content_orientation"`
   Country              string          `json:"country"`
   CreditCuepoints      CreditCuepoints `json:"credit_cuepoints"`
   Description          string          `json:"description"`
   DetailedType         string          `json:"detailed_type"`
   Directors            []string        `json:"directors"`
   Duration             int             `json:"duration"`
   GracenoteID          string          `json:"gracenote_id"`
   HasSubtitle          bool            `json:"has_subtitle"`
   HasTrailer           bool            `json:"has_trailer"`
   HeroImages           []string        `json:"hero_images"`
   ID                   string          `json:"id"`
   ImdbID               string          `json:"imdb_id"`
   ImportID             string          `json:"import_id"`
   InternalTags         []string        `json:"internal_tags"`
   IsCdc                bool            `json:"is_cdc"`
   IsReplay             bool            `json:"is_replay"`
   LandscapeImages      []string        `json:"landscape_images"`
   Lang                 string          `json:"lang"`
   LoginReason          string          `json:"login_reason"`
   Monetization         Monetization    `json:"monetization"`
   NeedsLogin           bool            `json:"needs_login"`
   PlayerType           string          `json:"player_type"`
   PolicyMatch          bool            `json:"policy_match"`
   Posterarts           []string        `json:"posterarts"`
   PublisherID          string          `json:"publisher_id"`
   Ratings              []Rating        `json:"ratings"`
   Subtitles            []Subtitle      `json:"subtitles"`
   Tags                 []string        `json:"tags"`
   Thumbnails           []string        `json:"thumbnails"`
   Title                string          `json:"title"`
   Trailers             []Trailer       `json:"trailers"`
   Type                 string          `json:"type"`
   UpdatedAt            string          `json:"updated_at"`
   URL                  string          `json:"url"`
   ValidDuration        int             `json:"valid_duration"`
   Version              int             `json:"version"`
   VersionID            string          `json:"version_id"`
   VideoMetadata        []VideoMetadata `json:"video_metadata"`
   VideoPreviewURL      string          `json:"video_preview_url"`
   VideoPreviews        []VideoPreview  `json:"video_previews"`
   VideoResources       []VideoResource `json:"video_resources"`
   Year                 int             `json:"year"`
}

type CreditCuepoints struct {
   EarlycreditsEnd   int `json:"earlycredits_end"`
   EarlycreditsStart int `json:"earlycredits_start"`
   IntroEnd          int `json:"intro_end"`
   IntroStart        int `json:"intro_start"`
   Postlude          int `json:"postlude"`
   Prelogue          int `json:"prelogue"`
   Prologue          int `json:"prologue"`
   RecapEnd          int `json:"recap_end"`
   RecapStart        int `json:"recap_start"`
}

type Monetization struct {
   CuePoints []int `json:"cue_points"`
}

type Rating struct {
   Code   string `json:"code"`
   System string `json:"system"`
   Value  string `json:"value"`
}

type Subtitle struct {
   Lang            string `json:"lang"`
   LangAlpha3      string `json:"lang_alpha3"`
   LangTranslation string `json:"lang_translation"`
   URL             string `json:"url"`
}

type Trailer struct {
   Duration int    `json:"duration"`
   ID       string `json:"id"`
   URL      string `json:"url"`
}

type VideoMetadata struct {
   Codec      string `json:"codec"`
   Resolution string `json:"resolution"`
   Type       string `json:"type"`
}

type VideoPreview struct {
   Source string `json:"source"`
   URL    string `json:"url"`
   UUID   string `json:"uuid"`
}

type VideoResource struct {
   Codec            string        `json:"codec"`
   GeneratorVersion string        `json:"generator_version"`
   LicenseServer    LicenseServer `json:"license_server"`
   Manifest         Manifest      `json:"manifest"`
   Resolution       string        `json:"resolution"`
   SsaiVersion      string        `json:"ssai_version"`
   TitanVersion     string        `json:"titan_version"`
   Type             string        `json:"type"`
}

type LicenseServer struct {
   AuthHeaderKey   string `json:"auth_header_key"`
   AuthHeaderValue string `json:"auth_header_value"`
   HdcpVersion     string `json:"hdcp_version"`
   URL             string `json:"url"`
}

type Manifest struct {
   Duration int    `json:"duration"`
   URL      string `json:"url"`
}

// GetContent constructs the request URL, fetches the CMS data,
// and returns the fully parsed Content struct.
func GetContent() (*Content, error) {
   endpoint, err := url.Parse("https://uapi.adrise.tv/cms/content")
   if err != nil {
      return nil, fmt.Errorf("failed to parse base URL: %w", err)
   }

   query := endpoint.Query()
   query.Set("content_id", "610572")
   query.Set("deviceId", "!")
   query.Add("limit_resolutions[]", "h264_1080p")
   query.Add("limit_resolutions[]", "h265_1080p")
   query.Set("platform", "web")
   query.Add("video_resources[]", "dash")
   query.Add("video_resources[]", "dash_widevine")
   endpoint.RawQuery = query.Encode()

   req, err := http.NewRequest("GET", endpoint.String(), nil)
   if err != nil {
      return nil, fmt.Errorf("failed to create request: %w", err)
   }

   req.Header.Set("user-agent", "Go-http-client/2.0")

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, fmt.Errorf("request failed: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   var content Content
   if err := json.NewDecoder(resp.Body).Decode(&content); err != nil {
      return nil, fmt.Errorf("failed to decode JSON response: %w", err)
   }

   return &content, nil
}
