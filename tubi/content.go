// File: content.go
package tubi

import (
   "encoding/json"
   "io"
   "net/url"

   "41.neocities.org/maya"
)

type Monetization struct {
   CuePoints []int `json:"cue_points"`
}

type VideoPreview struct {
   Source string `json:"source"`
   Url    string `json:"url"`
   Uuid   string `json:"uuid"`
}

type CreditCuepoints struct {
   Postlude          int `json:"postlude"`
   Prologue          int `json:"prologue"`
   IntroStart        int `json:"intro_start"`
   IntroEnd          int `json:"intro_end"`
   RecapStart        int `json:"recap_start"`
   RecapEnd          int `json:"recap_end"`
   EarlycreditsStart int `json:"earlycredits_start"`
   EarlycreditsEnd   int `json:"earlycredits_end"`
   Prelogue          int `json:"prelogue"`
}

type Subtitle struct {
   Url             string `json:"url"`
   Lang            string `json:"lang"`
   LangAlpha3      string `json:"lang_alpha3"`
   LangTranslation string `json:"lang_translation"`
}

type Trailer struct {
   Id       string `json:"id"`
   Url      string `json:"url"`
   Duration int    `json:"duration"`
}

type VideoMetadata struct {
   Type       string `json:"type"`
   Codec      string `json:"codec"`
   Resolution string `json:"resolution"`
}

type AudioTrack struct {
   Type        string `json:"type"`
   Lang        string `json:"lang"`
   DisplayName string `json:"display_name"`
}

type Manifest struct {
   Url      string `json:"url"`
   Duration int    `json:"duration"`
}

type LicenseServer struct {
   Url             string `json:"url"`
   HdcpVersion     string `json:"hdcp_version"`
   AuthHeaderKey   string `json:"auth_header_key"`
   AuthHeaderValue string `json:"auth_header_value"`
}

type VideoResource struct {
   Type             string        `json:"type"`
   Codec            string        `json:"codec"`
   AudioTracks      []AudioTrack  `json:"audio_tracks"`
   Resolution       string        `json:"resolution"`
   Manifest         Manifest      `json:"manifest"`
   LicenseServer    LicenseServer `json:"license_server"`
   TitanVersion     string        `json:"titan_version"`
   SsaiVersion      string        `json:"ssai_version"`
   GeneratorVersion string        `json:"generator_version"`
}

type Rating struct {
   Code   string `json:"code"`
   System string `json:"system"`
   Value  string `json:"value"`
}

type ContentResponse struct {
   HeroImages           []string        `json:"hero_images"`
   UpdatedAt            string          `json:"updated_at"`
   Monetization         Monetization    `json:"monetization"`
   GracenoteId          string          `json:"gracenote_id"`
   Type                 string          `json:"type"`
   AvailabilityStarts   string          `json:"availability_starts"`
   VideoPreviews        []VideoPreview  `json:"video_previews"`
   CanonicalId          string          `json:"canonical_id"`
   CreditCuepoints      CreditCuepoints `json:"credit_cuepoints"`
   IsCdc                bool            `json:"is_cdc"`
   AvailabilityDuration int             `json:"availability_duration"`
   Description          string          `json:"description"`
   VersionId            string          `json:"version_id"`
   Subtitles            []Subtitle      `json:"subtitles"`
   Id                   string          `json:"id"`
   Trailers             []Trailer       `json:"trailers"`
   Tags                 []string        `json:"tags"`
   Country              string          `json:"country"`
   AvailabilityEnds     string          `json:"availability_ends"`
   Directors            []string        `json:"directors"`
   VideoMetadata        []VideoMetadata `json:"video_metadata"`
   Version              int             `json:"version"`
   DetailedType         string          `json:"detailed_type"`
   VideoResources       []VideoResource `json:"video_resources"`
   LoginReason          string          `json:"login_reason"`
   Posterarts           []string        `json:"posterarts"`
   Backgrounds          []string        `json:"backgrounds"`
   NeedsLogin           bool            `json:"needs_login"`
   PolicyMatch          bool            `json:"policy_match"`
   VideoPreviewUrl      string          `json:"video_preview_url"`
   Duration             int             `json:"duration"`
   Actors               []string        `json:"actors"`
   IsReplay             bool            `json:"is_replay"`
   Url                  string          `json:"url"`
   InternalTags         []string        `json:"internal_tags"`
   PlayerType           string          `json:"player_type"`
   HasTrailer           bool            `json:"has_trailer"`
   PublisherId          string          `json:"publisher_id"`
   ImportId             string          `json:"import_id"`
   Ratings              []Rating        `json:"ratings"`
   Year                 int             `json:"year"`
   ValidDuration        int             `json:"valid_duration"`
   Lang                 string          `json:"lang"`
   ImdbId               string          `json:"imdb_id"`
   Title                string          `json:"title"`
   LandscapeImages      []string        `json:"landscape_images"`
   Thumbnails           []string        `json:"thumbnails"`
   ContentOrientation   string          `json:"content_orientation"`
   HasSubtitle          bool            `json:"has_subtitle"`
}

func GetContent(content_id string) (*ContentResponse, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "uapi.adrise.tv",
      Path:   "/cms/content",
   }

   query := url.Values{}
   query.Set("content_id", content_id)
   query.Set("deviceId", "!")
   query.Add("limit_resolutions[]", "h264_1080p")
   query.Add("limit_resolutions[]", "h265_1080p")
   query.Set("platform", "web")
   query.Add("video_resources[]", "dash")
   query.Add("video_resources[]", "dash_widevine")
   target.RawQuery = query.Encode()

   response, err := maya.Get(target, nil)
   if err != nil {
      return nil, err
   }
   defer response.Body.Close()

   payload, err := io.ReadAll(response.Body)
   if err != nil {
      return nil, err
   }

   var result ContentResponse
   if err := json.Unmarshal(payload, &result); err != nil {
      return nil, err
   }

   return &result, nil
}
