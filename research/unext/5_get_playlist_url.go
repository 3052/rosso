package unext

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

// GetPlaylistUrl fetches the MPD URL, License URL, and Play Token
func GetPlaylistUrl(client *http.Client, episodeCode string) (mpdUrl string, licenseUrl string, playToken string, err error) {
   baseURL := "https://cc.unext.jp/"

   params := url.Values{}
   params.Set("zxuid", "2cd3deff87a0")
   params.Set("zxemp", "29719881")
   params.Set("operationName", "cosmo_getPlaylistUrl")
   params.Set("variables", fmt.Sprintf(`{"code":"%s","playMode":"caption","bitrateLow":192,"bitrateHigh":null,"validationOnly":false}`, episodeCode))
   params.Set("extensions", `{"persistedQuery":{"version":1,"sha256Hash":"a2309e22a6819ff747cf9a389dd78db35fa3c386fac1d53461061ba20fa44e34"}}`)

   req, err := http.NewRequest("GET", baseURL+"?"+params.Encode(), nil)
   if err != nil {
      return "", "", "", err
   }

   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   req.Header.Set("Accept", "*/*")
   req.Header.Set("Accept-Language", "en-US,en;q=0.5")
   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("Apollographql-Client-Name", "cosmo")
   req.Header.Set("Apollographql-Client-Version", "v126.0-prod-017e302")
   req.Header.Set("Origin", "https://video.unext.jp")
   req.Header.Set("Referer", fmt.Sprintf("https://video.unext.jp/play/SID0020149/%s", episodeCode))

   resp, err := client.Do(req)
   if err != nil {
      return "", "", "", err
   }
   defer resp.Body.Close()

   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return "", "", "", err
   }

   var playlistData PlaylistResponse
   if err := json.Unmarshal(body, &playlistData); err != nil {
      return "", "", "", err
   }

   playToken = playlistData.Data.WebfrontPlaylistUrl.PlayToken
   if playToken == "" {
      return "", "", "", fmt.Errorf("failed to extract playToken from response: %s", string(body))
   }

   // Find the DASH profile and Widevine license URL
   for _, urlInfo := range playlistData.Data.WebfrontPlaylistUrl.UrlInfo {
      for _, profile := range urlInfo.MovieProfile {
         if profile.Type == "DASH" {
            mpdUrl = profile.PlaylistUrl
            for _, lic := range profile.LicenseUrlList {
               if lic.Type == "WIDEVINE" {
                  licenseUrl = lic.LicenseUrl
                  break
               }
            }
            break
         }
      }
   }

   if mpdUrl == "" || licenseUrl == "" {
      return "", "", "", fmt.Errorf("could not find DASH MPD URL or Widevine License URL in response: %s", string(body))
   }

   return mpdUrl, licenseUrl, playToken, nil
}

type PlaylistResponse struct {
   Data struct {
      WebfrontPlaylistUrl struct {
         PlayToken string `json:"playToken"`
         UrlInfo   []struct {
            MovieProfile []struct {
               Type           string `json:"type"`
               PlaylistUrl    string `json:"playlistUrl"`
               LicenseUrlList []struct {
                  Type       string `json:"type"`
                  LicenseUrl string `json:"licenseUrl"`
               } `json:"licenseUrlList"`
            } `json:"movieProfile"`
         } `json:"urlInfo"`
      } `json:"webfront_playlistUrl"`
   } `json:"data"`
}
