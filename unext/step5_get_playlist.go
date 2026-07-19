// step5_get_playlist.go
package unext

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

// AudioTrack describes an audio track language.
type AudioTrack struct {
   Lang     string `json:"lang"`
   IsNative bool   `json:"isNative"`
}

// DownloadTitleMeta is metadata about the downloaded title.
type DownloadTitleMeta struct {
   TitleInKatakana string   `json:"titleInKatakana"`
   Keywords        []string `json:"keywords"`
}

// GraphQLError represents a GraphQL error.
type GraphQLError struct {
   Message string `json:"message"`
}

// LicenseUrl describes a DRM license endpoint.
type LicenseUrl struct {
   Type       string `json:"type"`
   LicenseUrl string `json:"licenseUrl"`
}

// MovieAudio describes the audio codec type.
type MovieAudio struct {
   AudioType string `json:"audioType"`
}

// MoviePartsPosition describes a part of the movie (e.g. ENDING).
type MoviePartsPosition struct {
   Type             string  `json:"type"`
   FromSeconds      float64 `json:"fromSeconds"`
   EndSeconds       float64 `json:"endSeconds"`
   HasRemainingPart bool    `json:"hasRemainingPart"`
}

// MovieProfile describes a streaming profile (DASH, HLS, SMOOTH, etc.).
type MovieProfile struct {
   CdnId          string       `json:"cdnId"`
   Type           string       `json:"type"`
   PlaylistUrl    string       `json:"playlistUrl"`
   MovieAudioList []MovieAudio `json:"movieAudioList"`
   LicenseUrlList []LicenseUrl `json:"licenseUrlList"`
   AudioTrackList []AudioTrack `json:"audioTrackList"`
}

// PlaylistResponse is the JSON envelope returned by the GraphQL endpoint.
type PlaylistResponse struct {
   Data struct {
      WebfrontPlaylistUrl *PlaylistUrl `json:"webfront_playlistUrl"`
   } `json:"data"`
   Errors []GraphQLError `json:"errors"`
}

// PlaylistResult contains error info from the playlist request.
type PlaylistResult struct {
   ErrorCode    string `json:"errorCode"`
   ErrorMessage string `json:"errorMessage"`
}

// PlaylistUrl maps to the PlayList fragment on the PlaylistUrl type.
type PlaylistUrl struct {
   SubTitle          string             `json:"subTitle"`
   PlayToken         string             `json:"playToken"`
   PlayTokenHash     string             `json:"playTokenHash"`
   BeaconSpan        int                `json:"beaconSpan"`
   ResultStatus      int                `json:"resultStatus"`
   LicenseExpireDate string             `json:"licenseExpireDate"`
   IsKids            bool               `json:"isKids"`
   DownloadTitleMeta *DownloadTitleMeta `json:"downloadTitleMeta"`
   UrlInfo           []UrlInfo          `json:"urlInfo"`
   Result            *PlaylistResult    `json:"result"`
}

// Step5GetPlaylist fetches the playlist using the access token obtained in step 4.
// playMode must be either "dub" or "caption".
func Step5GetPlaylist(accessToken, code, playMode string) (*PlaylistUrl, error) {
   reqURL := &url.URL{
      Scheme: "https",
      Host:   "cc.unext.jp",
      Path:   "/",
   }

   variables := map[string]any{
      "code":               code,
      "playMode":           playMode,
      "bitrateLow":         192,
      "validationOnly":     false,
      "codec":              []string{"H264"},
      "playType":           "STREAMING",
      "keyOnly":            false,
      "mediaType":          "NORMAL",
      "disableRegionCheck": false,
   }

   body := map[string]any{
      "operationName": "Mad_Playlist",
      "variables":     variables,
      "query":         minPlaylistQuery,
   }

   bodyJSON, err := json.Marshal(body)
   if err != nil {
      return nil, fmt.Errorf("step5: marshalling body: %w", err)
   }

   req, err := http.NewRequest("POST", reqURL.String(), bytes.NewReader(bodyJSON))
   if err != nil {
      return nil, fmt.Errorf("step5: creating request: %w", err)
   }

   req.Header.Set("accept", "multipart/mixed;deferSpec=20220824, application/graphql-response+json, application/json")
   req.Header.Set("content-type", "application/json")
   req.Header.Set("apollo-require-preflight", "true")
   req.Header.Set("apollographql-client-name", "mad_for_mobile_jp.unext.mediaplayer")
   req.Header.Set("apollographql-client-version", "5.73.1")
   req.Header.Set("filmratingcode", "")
   req.Header.Set("u-device-type", "920")
   req.Header.Set("user-agent", "U-NEXT Phone App Android12 5.73.1 sdk_gphone64_x86_64")
   req.Header.Set("x-apollo-operation-name", "Mad_Playlist")
   req.Header.Set("x-forwarded-for", "159.26.119.122")
   req.Header.Set("authorization", "Bearer "+accessToken)

   resp, err := clientDo(req)
   if err != nil {
      return nil, fmt.Errorf("step5: sending request: %w", err)
   }
   defer resp.Body.Close()

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, fmt.Errorf("step5: reading response body: %w", err)
   }

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("step5: expected 200, got %d: %s", resp.StatusCode, string(respBody))
   }

   var plResp PlaylistResponse
   if err := json.Unmarshal(respBody, &plResp); err != nil {
      return nil, fmt.Errorf("step5: parsing response: %w (body starts with: %q)", err, string(respBody[:min(len(respBody), 50)]))
   }

   if len(plResp.Errors) > 0 {
      return nil, fmt.Errorf("step5: GraphQL error: %s", plResp.Errors[0].Message)
   }

   if plResp.Data.WebfrontPlaylistUrl == nil {
      return nil, fmt.Errorf("step5: webfront_playlistUrl was null")
   }

   return plResp.Data.WebfrontPlaylistUrl, nil
}

func (*PlaylistUrl) CachePath() string {
   return "rosso/unext/PlaylistUrl"
}

// MPDURL searches the playlist for the first DASH movie profile and returns
// its playlistUrl as a *url.URL with the play_token query parameter appended.
// Returns an error if no DASH profile is found or the URL cannot be parsed.
func (p *PlaylistUrl) MPDURL() (*url.URL, error) {
   for _, ui := range p.UrlInfo {
      for _, mp := range ui.MovieProfile {
         if mp.Type == "DASH" && mp.PlaylistUrl != "" {
            u, err := url.Parse(mp.PlaylistUrl)
            if err != nil {
               return nil, fmt.Errorf("parsing MPD URL: %w", err)
            }
            q := u.Query()
            q.Set("play_token", p.PlayToken)
            u.RawQuery = q.Encode()
            return u, nil
         }
      }
   }
   return nil, fmt.Errorf("no DASH movie profile found")
}

// SceneSearchList contains IMS (image search) URLs.
type SceneSearchList struct {
   IMS_AD1 string `json:"ims_ad1"`
   IMS_L   string `json:"ims_l"`
   IMS_M   string `json:"ims_m"`
   IMS_S   string `json:"ims_s"`
}

// UrlInfo represents one entry in the urlInfo array.
type UrlInfo struct {
   Code                   string               `json:"code"`
   StartPoint             float64              `json:"startPoint"`
   EndPoint               float64              `json:"endPoint"`
   ResumePoint            float64              `json:"resumePoint"`
   EndrollStartPosition   float64              `json:"endrollStartPosition"`
   CommodityCode          string               `json:"commodityCode"`
   SaleTypeCode           string               `json:"saleTypeCode"`
   CaptionFlg             bool                 `json:"captionFlg"`
   DubFlg                 bool                 `json:"dubFlg"`
   SceneSearchList        *SceneSearchList     `json:"sceneSearchList"`
   MovieProfile           []MovieProfile       `json:"movieProfile"`
   MoviePartsPositionList []MoviePartsPosition `json:"moviePartsPositionList"`
}
