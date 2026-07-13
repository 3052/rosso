package main

import (
   "net/http"
   "net/url"
   "os"
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

func main() {
   client := &http.Client{}
   reqURL := &url.URL{
      Scheme: "https",
      Host:   "cc.unext.jp",
      Path:   "/",
   }
   q := url.Values{}
   q.Add("operationName", "cosmo_getPlaylistUrl")
   q.Add("query", cosmoGetPlaylistURLQuery)
   q.Add("variables", "{\"code\":\"ED00092859\",\"playMode\":\"caption\",\"bitrateLow\":192,\"bitrateHigh\":null,\"validationOnly\":false}")
   reqURL.RawQuery = q.Encode()
   req, err := http.NewRequest("GET", reqURL.String(), nil)
   if err != nil {
      panic(err)
   }
   req.Header.Add("sec-fetch-dest", "empty")
   req.Header.Add("sec-fetch-mode", "cors")
   req.Header.Add("sec-fetch-site", "same-site")
   req.Header.Add("te", "trailers")
   req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   req.Header.Add("accept", "*/*")
   req.Header.Add("accept-encoding", "identity")
   req.Header.Add("accept-language", "en-US,en;q=0.5")
   req.Header.Add("apollographql-client-name", "cosmo")
   req.Header.Add("apollographql-client-version", "v126.0-prod-017e302")
   req.Header.Add("content-type", "application/json")
   req.Header.Add("origin", "https://video.unext.jp")
   req.Header.Add("priority", "u=0")
   req.AddCookie(&http.Cookie{Name: "__td_signed", Value: "true"})
   req.AddCookie(&http.Cookie{Name: "_at", Value: "eyJhbGciOiJSUzI1NiIsImtpZCI6Ijl1dHRiYWpuZmkiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOlsidW5leHQiXSwiZXhwIjoxNzgzOTE5Mzk5LCJpYXQiOjE3ODM4Nzk3OTgsImlzcyI6Imh0dHBzOi8vb2F1dGgudW5leHQuanAiLCJqdGkiOiJmOWU3OTAxYS1mMGYyLTRkY2UtOWU3YS04ZGJjMjM0MjljZTMiLCJzY3AiOlsidW5leHQiLCJvZmZsaW5lIl0sInN1YiI6IlBNMDY0Njc0MDY0IiwidW5leHQiOiJleUpoYkdjaU9pSlNVMEV0VDBGRlVDSXNJbVZ1WXlJNklrRXlOVFpIUTAwaUxDSnJhV1FpT2lJNGRIUmtjV0YyWjNkaUluMC42VTdFQ2Rob01DVDJKX2pfeVVMYmhPRTZuTTFncUpzVlBPMDByd3V0MjRMZmRaYkltZFZ5QUZIeGdwZEZwNFVSQ09BVVVTa21hM3c5RDNVS2pYdjQ0SWo1bTd2alpYVzJLWk1Sbk9lNTYyOUNWSm9mY05YakU2cWFUMC1CZWZJV3QwdjNWQ1pjTmpCWVN5dGE4aWlYeHNwckttYWhNZ1lMSjN1S0k3Ti1xbXVQMlBuLUpETFZFZ0NndmIxbWJtZUVBbGUxX1NULVMzTGNOQWpzSGhES1BzNHBaSHNISlFja1pYZ3MwelNxU2dQeUtEaDlJRUg2YUFfTWdDemFnVTlWU3phcGxxdHlqMlg2Vm4xNEo4Y2lHc1pRV2FXdlkwUk5xSVhNbTEzeE9fRHUxZVl0djNvX2ZXV2gtZzZfZVEyajh5bkNaaVphMTF2SVl4QmlHZjdxZ1EuMGdIMXRiUXg4NTNWdmZBei5UTUZ2NlBpaG5VOEFXUDBhajJPX3RoMmczcXo2UnZBSWE2bG9oNHdMcXRBV0pOVElyd25HakR0dEZRVmdETXFHZHNnV28zTHlRNEhvbHVpNktsN1RqVzRiTWU5cGRTN282aGlja3hZWGRFNmd3YXJwWUNsSWtGaG1peHkyZWtrWjZaS1pnbjc1cko1Qkl2bEpPVmxfYnBkY2tLTVFSSnd0RzFOeWxsZ0tXRzFacTcxdVlmQ0VJNmFSZDk0bVhJdjFhcFZ5YU1DMnA1eHdSNDBfSFRwWVhKdUJRTzNKOFVIMHJQSHllWGlBbHdqZGtKTnVhT013eVE3R1dwTG91WVVwaFlHQXZidGljZF9Wdi1FckpVWmlzRFVtR3hvYnMxeko4VmktY3JFQmxHZEpka2gySkR3cWN3VUF6eHV0RmhQUTNJUEdRYUVyZkhXUEhhMl8yWVRHaWlBS3cxUmtaY2ZOOTNXSnRqWm5rRmQwOGVVeHZHV0hnNHB2ZGl0VHZqdy5GSHpwMFNCRjg5OElwd1JKTkNEX25BIiwidmVyc2lvbiI6InYxLjAuMiJ9.CitU0PB_sqt6jzUi9rfHZAi3Enf7ZOj-BhXeNaLsjgGsVdoqV74_DG_9_LF_M4CPWcPkciaMLSberKpWCubX4E0JAm-qjHw0SHemLQv_2ZoWFpsuolETbo4dFmhtOYs6NHUZbLcK_MdN8mwGrun6I9RAlk3liUq6s0iMQ7zcoID0byiRBGeNsDGlIIWHKEYa0pqbDYi0jPJXXJHxilhRc983q2Pyhzd1J9Dk7LcyhFAs4_t0aniuKXbQreYxuTzz4OqrXAGhQJEc9n6DtCbDSN01U_cDHvBLtXxbw8l4SNokFmZJo8DYMk3lkQl25b3YI3vR5k59vMD3W6-UI8GBBA"})
   resp, err := client.Do(req)
   if err != nil {
      panic(err)
   }
   if err := resp.Write(os.Stdout); err != nil {
      panic(err)
   }
}
