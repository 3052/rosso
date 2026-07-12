// playlist_url.go
package unext

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

const cosmoGetPlaylistURLQuery = `query cosmo_getPlaylistUrl($code: String, $playMode: String, $bitrateLow: Int, $bitrateHigh: Int, $validationOnly: Boolean) {
  webfront_playlistUrl(
    code: $code
    playMode: $playMode
    bitrateLow: $bitrateLow
    bitrateHigh: $bitrateHigh
    validationOnly: $validationOnly
  ) {
    subTitle
    playToken
    playTokenHash
    beaconSpan
    result {
      errorCode
      errorMessage
      __typename
    }
    resultStatus
    licenseExpireDate
    urlInfo {
      code
      startPoint
      resumePoint
      endPoint
      endrollStartPosition
      holderId
      saleTypeCode
      sceneSearchList {
        IMS_AD1
        IMS_L
        IMS_M
        IMS_S
        __typename
      }
      movieProfile {
        cdnId
        type
        playlistUrl
        movieAudioList {
          audioType
          __typename
        }
        licenseUrlList {
          type
          licenseUrl
          __typename
        }
        __typename
      }
      umcContentId
      movieSecurityLevelCode
      captionFlg
      dubFlg
      commodityCode
      movieAudioList {
        audioType
        __typename
      }
      moviePartsPositionList {
        type
        fromSeconds
        endSeconds
        hasRemainingPart
        __typename
      }
      __typename
    }
    __typename
  }
}
`

// PlaylistURLInputs holds only the values that come from previous requests or session state.
type PlaylistURLInputs struct {
   EpisodeCode string // from cosmo_getVideoTitleEpisodes / cosmo_getVideoTitle
   ZXUID       string // session-generated tracking ID
   ZXEMP       string // session tracking ID
   Cookie      string // auth session cookies (_at, _rt, _ut, _st, etc.)
}

// PlaylistURLResult holds the parsed result of the cosmo_getPlaylistUrl request.
type PlaylistURLResult struct {
   PlayToken     string
   DashMPDURL    string
   HLSURL        string
   SmoothURL     string
   MovieFileCode string
   LicenseURLs   map[string]string
   Raw           []byte
   Response      *graphQLResponse
}

// GetPlaylistURL performs the cosmo_getPlaylistUrl request and returns the playlist URLs.
func GetPlaylistURL(in *PlaylistURLInputs) (*PlaylistURLResult, error) {
   if in == nil {
      return nil, fmt.Errorf("inputs must not be nil")
   }

   variables := map[string]interface{}{
      "code":           in.EpisodeCode,
      "playMode":       "caption",
      "bitrateLow":     192,
      "bitrateHigh":    nil,
      "validationOnly": false,
   }
   variablesJSON, _ := json.Marshal(variables)

   q := url.Values{}
   q.Set("zxuid", in.ZXUID)
   q.Set("zxemp", in.ZXEMP)
   q.Set("operationName", "cosmo_getPlaylistUrl")
   q.Set("variables", string(variablesJSON))
   q.Set("query", cosmoGetPlaylistURLQuery)

   reqURL := "https://cc.unext.jp/?" + q.Encode()

   req, err := http.NewRequest(http.MethodGet, reqURL, nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   req.Header.Set("Accept", "*/*")
   req.Header.Set("Accept-Language", "en-US,en;q=0.5")
   req.Header.Set("Accept-Encoding", "identity")
   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("Referer", "https://video.unext.jp/")
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
