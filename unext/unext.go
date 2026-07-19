// unext.go
package unext

import (
   "bytes"
   _ "embed"
   "encoding/json"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/url"
   "strings"
)

//go:embed mad_all_episodes.graphql
var allEpisodesQuery string

//go:embed mad_playlist.graphql
var playlistQuery string

// Step6GetLicense POSTs a Widevine license challenge to the U-NEXT license
// proxy and returns the raw license response bytes.
//
// challenge is the binary SignedMessage (protobuf) produced by a Widevine CDM.
// The play_token must match the one used to fetch the MPD.
func Step6GetLicense(licenseURL *url.URL, playToken string, challenge []byte) ([]byte, error) {
   query := licenseURL.Query()
   query.Set("play_token", playToken)
   licenseURL.RawQuery = query.Encode()
   req, err := http.NewRequest("POST", licenseURL.String(), bytes.NewReader(challenge))
   if err != nil {
      return nil, fmt.Errorf("step6: creating request: %w", err)
   }
   resp, err := clientDo(req)
   if err != nil {
      return nil, fmt.Errorf("step6: sending request: %w", err)
   }
   defer resp.Body.Close()

   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, fmt.Errorf("step6: reading response body: %w", err)
   }

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("step6: expected 200, got %d: %s", resp.StatusCode, string(body))
   }

   return body, nil
}

func clientDo(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
}

func clientDoNoRedirect(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   client := &http.Client{
      CheckRedirect: func(*http.Request, []*http.Request) error {
         return http.ErrUseLastResponse
      },
   }
   return client.Do(req)
}

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
   }
   body := map[string]any{
      "query": playlistQuery,
      "variables": map[string]any{
         "code":               code,
         "playMode":           playMode,
         "bitrateLow":         192,
         "codec":              []string{"H264"},
         "disableRegionCheck": true,
         "keyOnly":            false,
         "playType":           "STREAMING",
         "validationOnly":     false,
      },
   }
   bodyJSON, err := json.Marshal(body)
   if err != nil {
      return nil, fmt.Errorf("step5: marshalling body: %w", err)
   }
   req, err := http.NewRequest("POST", reqURL.String(), bytes.NewReader(bodyJSON))
   if err != nil {
      return nil, fmt.Errorf("step5: creating request: %w", err)
   }
   req.Header.Set("authorization", "Bearer "+accessToken)
   req.Header.Set("content-type", "application/json")
   resp, err := clientDo(req)
   if err != nil {
      return nil, fmt.Errorf("step5: sending request: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("step5: expected 200, got %d", resp.StatusCode)
   }

   var plResp PlaylistResponse
   if err := json.NewDecoder(resp.Body).Decode(&plResp); err != nil {
      return nil, fmt.Errorf("step5: parsing response: %w", err)
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

// WidevineLicenseURL searches the playlist for the first DASH movie profile
// with a WIDEVINE license URL and returns it (without query parameters).
// Returns an error if no such profile is found.
func (p *PlaylistUrl) WidevineLicenseURL() (*url.URL, error) {
   for _, ui := range p.UrlInfo {
      for _, mp := range ui.MovieProfile {
         if mp.Type != "DASH" {
            continue
         }
         for _, lu := range mp.LicenseUrlList {
            if lu.Type == "WIDEVINE" && lu.LicenseUrl != "" {
               u, err := url.Parse(lu.LicenseUrl)
               if err != nil {
                  return nil, fmt.Errorf("parsing license URL: %w", err)
               }
               return u, nil
            }
         }
      }
   }
   return nil, fmt.Errorf("no DASH movie profile with WIDEVINE license found")
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

// Refresh exchanges the receiver's RefreshToken for a new set of tokens
// and writes the result back into the receiver.
func (t *TokenResponse) Refresh() error {
   tokenURL := "https://oauth.unext.jp/oauth2/token"
   form := url.Values{}
   form.Set("refresh_token", t.RefreshToken)
   form.Set("grant_type", "refresh_token")
   form.Set("client_id", "unextAndroidApp")
   req, err := http.NewRequest("POST", tokenURL, strings.NewReader(form.Encode()))
   if err != nil {
      return fmt.Errorf("refresh: creating request: %w", err)
   }
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   resp, err := clientDo(req)
   if err != nil {
      return fmt.Errorf("refresh: sending request: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return fmt.Errorf("refresh: expected 200, got %d", resp.StatusCode)
   }

   var newToken TokenResponse
   if err := json.NewDecoder(resp.Body).Decode(&newToken); err != nil {
      return fmt.Errorf("refresh: parsing response: %w", err)
   }

   // Write the new tokens back into the receiver.
   *t = newToken
   return nil
}
