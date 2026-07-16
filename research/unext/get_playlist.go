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

// PlaylistResponse represents the GraphQL response from the playlist endpoint.
type PlaylistResponse struct {
   Data struct {
      WebfrontPlaylistUrl struct {
         SubTitle      string `json:"subTitle"`
         PlayToken     string `json:"playToken"`
         PlayTokenHash string `json:"playTokenHash"`
         BeaconSpan    int    `json:"beaconSpan"`
         Result        struct {
            ErrorCode    string `json:"errorCode"`
            ErrorMessage string `json:"errorMessage"`
         } `json:"result"`
         ResultStatus      int    `json:"resultStatus"`
         LicenseExpireDate string `json:"licenseExpireDate"`
         UrlInfo           []struct {
            Code         string `json:"code"`
            MovieProfile []struct {
               CdnId          string `json:"cdnId"`
               Type           string `json:"type"`
               PlaylistUrl    string `json:"playlistUrl"`
               MovieAudioList []struct {
                  AudioType string `json:"audioType"`
               } `json:"movieAudioList"`
               LicenseUrlList []struct {
                  Type       string `json:"type"`
                  LicenseUrl string `json:"licenseUrl"`
               } `json:"licenseUrlList"`
            } `json:"movieProfile"`
         } `json:"urlInfo"`
      } `json:"webfront_playlistUrl"`
   } `json:"data"`
   Errors []struct {
      Message string `json:"message"`
   } `json:"errors"`
}

// GetPlaylistURL calls the cosmo GraphQL endpoint with hardcoded parameters.
// Only the accessToken is required.
func GetPlaylistURL(client *http.Client, accessToken string) (*PlaylistResponse, error) {
   reqURL := &url.URL{
      Scheme: "https",
      Host:   "cc.unext.jp",
      Path:   "/",
   }

   q := url.Values{}
   q.Add("operationName", "cosmo_getPlaylistUrl")
   q.Add("query", cosmoGetPlaylistURLQuery)
   q.Add("variables", `{"code":"ED00092859","playMode":"caption","bitrateLow":192,"bitrateHigh":null,"validationOnly":false}`)
   reqURL.RawQuery = q.Encode()

   req, err := http.NewRequest("GET", reqURL.String(), nil)
   if err != nil {
      return nil, fmt.Errorf("get_playlist: creating request: %w", err)
   }

   req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   req.Header.Add("accept", "*/*")
   req.Header.Add("accept-language", "en-US,en;q=0.5")
   req.Header.Add("content-type", "application/json")
   req.Header.Add("origin", "https://video.unext.jp")
   req.Header.Add("priority", "u=0")
   req.Header.Add("sec-fetch-dest", "empty")
   req.Header.Add("sec-fetch-mode", "cors")
   req.Header.Add("sec-fetch-site", "same-site")
   req.Header.Add("te", "trailers")
   req.Header.Add("apollographql-client-name", "cosmo")
   req.Header.Add("apollographql-client-version", "v126.0-prod-017e302")
   req.Header.Add("authorization", "Bearer "+accessToken)

   resp, err := client.Do(req)
   if err != nil {
      return nil, fmt.Errorf("get_playlist: sending request: %w", err)
   }
   defer resp.Body.Close()

   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, fmt.Errorf("get_playlist: reading response: %w", err)
   }

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("get_playlist: unexpected status %d: %s", resp.StatusCode, string(body))
   }

   var playlistResp PlaylistResponse
   if err := json.Unmarshal(body, &playlistResp); err != nil {
      return nil, fmt.Errorf("get_playlist: parsing response: %w", err)
   }

   if len(playlistResp.Errors) > 0 {
      return nil, fmt.Errorf("get_playlist: graphql errors: %v", playlistResp.Errors)
   }

   if playlistResp.Data.WebfrontPlaylistUrl.ResultStatus != 200 {
      return nil, fmt.Errorf("get_playlist: resultStatus %d (expected 200, possibly geo-blocked or region-restricted)",
         playlistResp.Data.WebfrontPlaylistUrl.ResultStatus)
   }

   return &playlistResp, nil
}

// GetDASHPlaylistURL finds the DASH MPD URL in the response and appends the play_token query parameter.
func (p *PlaylistResponse) GetDASHPlaylistURL() (*url.URL, error) {
   playToken := p.Data.WebfrontPlaylistUrl.PlayToken
   if playToken == "" {
      return nil, fmt.Errorf("play token is empty")
   }

   for _, urlInfo := range p.Data.WebfrontPlaylistUrl.UrlInfo {
      for _, profile := range urlInfo.MovieProfile {
         if profile.Type == "DASH" {
            parsedURL, err := url.Parse(profile.PlaylistUrl)
            if err != nil {
               return nil, fmt.Errorf("parsing DASH playlist URL: %w", err)
            }

            queries := parsedURL.Query()
            queries.Add("play_token", playToken)
            parsedURL.RawQuery = queries.Encode()

            return parsedURL, nil
         }
      }
   }

   return nil, fmt.Errorf("no DASH MPD URL found in response")
}
