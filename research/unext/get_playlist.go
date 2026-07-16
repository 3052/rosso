package unext

import (
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

// GetPlaylistURL calls the cosmo GraphQL endpoint with hardcoded parameters.
// Only the accessToken is required.
func GetPlaylistURL(client *http.Client, accessToken string) ([]byte, error) {
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

   return body, nil
}
