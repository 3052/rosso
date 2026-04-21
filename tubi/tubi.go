package tubi

import (
   "41.neocities.org/maya"
   "encoding/json"
   "fmt"
   "io"
   "net/url"
   "strconv"
)

func GetContent(contentId int) (*Content, error) {
   query := make(url.Values)
   query.Set("content_id", strconv.Itoa(contentId))
   query.Set("deviceId", "!")
   query.Add("limit_resolutions[]", "h264_1080p")
   query.Add("limit_resolutions[]", "h265_1080p")
   query.Set("platform", "web")
   query.Add("video_resources[]", "dash")
   query.Add("video_resources[]", "dash_widevine")
   resp, err := maya.Get(&url.URL{
      Scheme:   "https",
      Host:     "uapi.adrise.tv",
      Path:     "/cms/content",
      RawQuery: query.Encode(),
   }, nil)
   if err != nil {
      return nil, fmt.Errorf("request failed: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   var content Content
   if err := json.NewDecoder(resp.Body).Decode(&content); err != nil {
      return nil, fmt.Errorf("failed to decode JSON response: %w", err)
   }

   return &content, nil
}

func (s *LicenseServer) PostLicense(body []byte) ([]byte, error) {
   target, err := url.Parse(s.Url)
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      target, map[string]string{"content-type": "application/x-protobuf"}, body,
   )
   if err != nil {
      return nil, fmt.Errorf("request failed: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   body, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, fmt.Errorf("failed to read response body: %w", err)
   }

   return body, nil
}

func (v *VideoResource) GetManifest() (*url.URL, error) {
   return url.Parse(v.Manifest.Url)
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

// Content represents the exact structure of the provided content.json, omitting interface and null fields.
type Content struct {
   Actors               []string        `json:"actors"`
   AdLanguages          []string        `json:"ad_languages"`
   AvailabilityDuration int             `json:"availability_duration"`
   AvailabilityEnds     string          `json:"availability_ends"`
   AvailabilityStarts   string          `json:"availability_starts"`
   Backgrounds          []string        `json:"backgrounds"`
   CanonicalId          string          `json:"canonical_id"`
   ContentOrientation   string          `json:"content_orientation"`
   Country              string          `json:"country"`
   CreditCuepoints      CreditCuepoints `json:"credit_cuepoints"`
   Description          string          `json:"description"`
   DetailedType         string          `json:"detailed_type"`
   Directors            []string        `json:"directors"`
   Duration             int             `json:"duration"`
   GracenoteId          string          `json:"gracenote_id"`
   HasSubtitle          bool            `json:"has_subtitle"`
   HasTrailer           bool            `json:"has_trailer"`
   HeroImages           []string        `json:"hero_images"`
   Id                   string          `json:"id"`
   ImdbId               string          `json:"imdb_id"`
   ImportId             string          `json:"import_id"`
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
   PublisherId          string          `json:"publisher_id"`
   Ratings              []Rating        `json:"ratings"`
   Subtitles            []Subtitle      `json:"subtitles"`
   Tags                 []string        `json:"tags"`
   Thumbnails           []string        `json:"thumbnails"`
   Title                string          `json:"title"`
   Trailers             []Trailer       `json:"trailers"`
   Type                 string          `json:"type"`
   UpdatedAt            string          `json:"updated_at"`
   Url                  string          `json:"url"`
   ValidDuration        int             `json:"valid_duration"`
   Version              int             `json:"version"`
   VersionId            string          `json:"version_id"`
   VideoMetadata        []VideoMetadata `json:"video_metadata"`
   VideoPreviewUrl      string          `json:"video_preview_url"`
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
   Url             string `json:"url"`
}

type Trailer struct {
   Duration int    `json:"duration"`
   Id       string `json:"id"`
   Url      string `json:"url"`
}

type VideoMetadata struct {
   Codec      string `json:"codec"`
   Resolution string `json:"resolution"`
   Type       string `json:"type"`
}

type VideoPreview struct {
   Source string `json:"source"`
   Url    string `json:"url"`
   Uuid   string `json:"uuid"`
}

type LicenseServer struct {
   AuthHeaderKey   string `json:"auth_header_key"`
   AuthHeaderValue string `json:"auth_header_value"`
   HdcpVersion     string `json:"hdcp_version"`
   Url             string `json:"url"`
}

type Manifest struct {
   Duration int    `json:"duration"`
   Url      string `json:"url"`
}
