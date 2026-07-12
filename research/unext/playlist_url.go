// playlist_url.go
package unext

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

const persistedQueryHashPlaylistURL = "a2309e22a6819ff747cf9a389dd78db35fa3c386fac1d53461061ba20fa44e34"

// Inputs that come from previous requests / session state.
type PlaylistURLInputs struct {
   EpisodeCode    string // e.g. "ED00092859" (from cosmo_getVideoTitleEpisodes / cosmo_getVideoTitle)
   PlayMode       string // e.g. "caption" (or "dub")
   BitrateLow     int    // e.g. 192
   BitrateHigh    *int   // null in capture; pass nil to send null
   ValidationOnly bool   // false in capture

   // Tracking / client identifiers (zxuid looks generated per session; zxemp appears stable).
   ZXUID string // e.g. "c1ee23a13f82"
   ZXEMP string // e.g. "29719883"

   // Auth/session cookie string (must include _at access token, _rt refresh token, _ut, _st, etc.)
   Cookie string

   // Optional overrides
   UserAgent string
   Referer   string // e.g. https://video.unext.jp/play/SID0020149/ED00092859
}

// PlaylistURLResult holds the parsed result of the cosmo_getPlaylistUrl request.
type PlaylistURLResult struct {
   PlayToken     string
   DashMPDURL    string
   HLSURL        string
   SmoothURL     string
   MovieFileCode string
   LicenseURLs   map[string]string // drm type -> license url
   Raw           []byte
   Response      *graphQLResponse
}

// GetPlaylistURL performs the cosmo_getPlaylistUrl request and returns the playlist URLs.
// It prefers the DASH MPD URL but also extracts HLS and Smooth urls.
func GetPlaylistURL(in PlaylistURLInputs) (*PlaylistURLResult, error) {
   variables := map[string]interface{}{
      "code":           in.EpisodeCode,
      "playMode":       in.PlayMode,
      "bitrateLow":     in.BitrateLow,
      "bitrateHigh":    in.BitrateHigh,
      "validationOnly": in.ValidationOnly,
   }
   variablesJSON, _ := json.Marshal(variables)

   extensions := map[string]interface{}{
      "persistedQuery": map[string]interface{}{
         "version":    1,
         "sha256Hash": persistedQueryHashPlaylistURL,
      },
   }
   extensionsJSON, _ := json.Marshal(extensions)

   q := url.Values{}
   q.Set("zxuid", in.ZXUID)
   q.Set("zxemp", in.ZXEMP)
   q.Set("operationName", "cosmo_getPlaylistUrl")
   q.Set("variables", string(variablesJSON))
   q.Set("extensions", string(extensionsJSON))

   reqURL := "https://cc.unext.jp/?" + q.Encode()

   req, err := http.NewRequest(http.MethodGet, reqURL, nil)
   if err != nil {
      return nil, err
   }

   ua := in.UserAgent
   if ua == "" {
      ua = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0"
   }
   req.Header.Set("User-Agent", ua)
   req.Header.Set("Accept", "*/*")
   req.Header.Set("Accept-Language", "en-US,en;q=0.5")
   req.Header.Set("Accept-Encoding", "identity")
   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("Referer", in.Referer)
   req.Header.Set("Origin", "https://video.unext.jp")
   req.Header.Set("apollographql-client-name", "cosmo")
   req.Header.Set("apollographql-client-version", "v126.0-prod-017e302")
   req.Header.Set("Sec-Fetch-Dest", "empty")
   req.Header.Set("Sec-Fetch-Mode", "cors")
   req.Header.Set("Sec-Fetch-Site", "same-site")
   if in.Cookie != "" {
      req.Header.Set("Cookie", in.Cookie)
   }

   httpResp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer httpResp.Body.Close()

   raw, err := io.ReadAll(httpResp.Body)
   if err != nil {
      return nil, err
   }

   resp := &graphQLResponse{}
   if err := json.Unmarshal(raw, resp); err != nil {
      return nil, fmt.Errorf("failed to parse JSON: %w; body: %s", err, string(raw))
   }
   if len(resp.Errors) > 0 {
      return &PlaylistURLResult{Raw: raw, Response: resp}, fmt.Errorf("graphql error: %s", resp.Errors[0].Message)
   }

   result := &PlaylistURLResult{
      PlayToken:   resp.Data.WebfrontPlaylistURL.PlayToken,
      LicenseURLs: make(map[string]string),
      Raw:         raw,
      Response:    resp,
   }

   for _, ui := range resp.Data.WebfrontPlaylistURL.URLInfo {
      if result.MovieFileCode == "" {
         result.MovieFileCode = ui.Code
      }
      for _, mp := range ui.MovieProfile {
         switch mp.Type {
         case "DASH":
            if result.DashMPDURL == "" {
               result.DashMPDURL = mp.PlaylistURL
            }
         case "HLS_FP", "HLS":
            if result.HLSURL == "" {
               result.HLSURL = mp.PlaylistURL
            }
         case "SMOOTH":
            if result.SmoothURL == "" {
               result.SmoothURL = mp.PlaylistURL
            }
         }
         for _, lic := range mp.LicenseURLList {
            if _, ok := result.LicenseURLs[lic.Type]; !ok {
               result.LicenseURLs[lic.Type] = lic.LicenseURL
            }
         }
      }
   }

   if result.DashMPDURL == "" {
      return result, fmt.Errorf("no DASH profile found in response; resultStatus=%d", resp.Data.WebfrontPlaylistURL.ResultStatus)
   }
   return result, nil
}

// Minimal structs to parse the GraphQL response.
type graphQLResponse struct {
   Data struct {
      WebfrontPlaylistURL struct {
         PlayToken    string `json:"playToken"`
         ResultStatus int    `json:"resultStatus"`
         URLInfo      []struct {
            Code         string `json:"code"`
            MovieProfile []struct {
               CdnID          string `json:"cdnId"`
               Type           string `json:"type"`
               PlaylistURL    string `json:"playlistUrl"`
               LicenseURLList []struct {
                  Type       string `json:"type"`
                  LicenseURL string `json:"licenseUrl"`
               } `json:"licenseUrlList"`
            } `json:"movieProfile"`
         } `json:"urlInfo"`
      } `json:"webfront_playlistUrl"`
   } `json:"data"`
   Errors []struct {
      Message string `json:"message"`
   } `json:"errors"`
}
