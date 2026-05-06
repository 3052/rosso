package tubi

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "net/url"
   "strconv"
)

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
   Resolution       string        `json:"resolution"`
   Manifest         Manifest      `json:"manifest"`
   LicenseServer    LicenseServer `json:"license_server"`
   TitanVersion     string        `json:"titan_version"`
   SsaiVersion      string        `json:"ssai_version"`
   GeneratorVersion string        `json:"generator_version"`
}

type ContentResponse struct {
   HeroImages           []string        `json:"hero_images"`
   UpdatedAt            string          `json:"updated_at"`
   GracenoteId          string          `json:"gracenote_id"`
   Type                 string          `json:"type"`
   AvailabilityStarts   string          `json:"availability_starts"`
   IsCdc                bool            `json:"is_cdc"`
   AvailabilityDuration int             `json:"availability_duration"`
   Description          string          `json:"description"`
   Id                   int             `json:"id,string"`
   Tags                 []string        `json:"tags"`
   Country              string          `json:"country"`
   AvailabilityEnds     string          `json:"availability_ends"`
   Directors            []string        `json:"directors"`
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

func GetContent(contentId int) (*ContentResponse, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "uapi.adrise.tv",
      Path:   "/cms/content",
   }
   query := url.Values{}
   query.Set("content_id", strconv.Itoa(contentId))
   query.Set("deviceId", "!")
   query.Add("limit_resolutions[]", "h264_1080p")
   query.Add("limit_resolutions[]", "h265_1080p")
   query.Set("platform", "web")
   query.Add("video_resources[]", "dash")
   query.Add("video_resources[]", "dash_widevine")
   target.RawQuery = query.Encode()

   resp, err := maya.Get(target, nil)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var content ContentResponse
   if err := json.NewDecoder(resp.Body).Decode(&content); err != nil {
      return nil, err
   }

   // Return an error if no VideoResources were found in the response
   if len(content.VideoResources) == 0 {
      return nil, errors.New("no video resources found for this content")
   }

   return &content, nil
}
