package main

import (
   "bytes"
   "encoding/json"
   "io"
   "log"
   "net/http"
   "net/url"
   "strings"
)

func main() {
   client := &http.Client{}
   reqURL := &url.URL{
      Scheme: "https",
      Host:   "ab8mt4dd97et.na.api.amazonvideo.com",
      Path:   "/playback/prs/GetVodPlaybackResources",
   }
   q := url.Values{}
   q.Add("deviceID", "uuid1aec7cc26f4440b1a313c1c7dfecd881")
   q.Add("deviceTypeID", "A2HYAJ0FEWP6N3")
   reqURL.RawQuery = q.Encode()
   value := map[string]any{
      "globalParameters": map[string]any{
         "playbackEnvelope":       "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",
         "deviceCapabilityFamily": "AndroidPlayer",
      },
      "vodPlaylistedPlaybackUrlsRequest": map[string]any{
         "device": map[string]any{
            "displayHeight":                  2160,
            "displayWidth":                   3840,
            "hdcpLevel":                      "2.3",
            "maxVideoResolution":             "2160p",
            "supportedStreamingTechnologies": []string{"DASH"},
            "streamingTechnologies": map[string]any{
               "DASH": map[string]any{
                  "bitrateAdaptations":  []string{"CVBR"},
                  "codecs":              []string{"H265"},
                  "drmStrength":         "L40",
                  "drmType":             "PLAYREADY",
                  "dynamicRangeFormats": []string{"HDR10", "None"},
               },
            },
         },
         "playbackSettingsRequest": map[string]any{
            "firmware": "MTC/MStar-T22/sadang:11/RP1A.200622.001/59334:user/release-keys",
            "titleId":  "0NUYZGCVPFY6LG1Z3TUTYQ1AZK",
         },
      },
   }
   reqData, err := json.Marshal(value)
   if err != nil {
      panic(err)
   }
   req, err := http.NewRequest("POST", reqURL.String(), bytes.NewBuffer(reqData))
   if err != nil {
      panic(err)
   }
   //go
   req.Header.Set("Authorization", "Bearer Atna|EwMDICwPZSUUUHzHlSTUj3J-zTJugQL871-Gh-juO8menKhfPCef8ALrLp5AoLOf2h4MrdCbi48lNoUhN9Y-6JIZjv32j0Acbv8Hhdx9VBwtIW_Wt8HdWIYA2SND4ca4WlkzI0W_Dn8Ge4f5GCf3fsQMpMFKwrTWrQFEVB9rPD9oAEEOqJGYwNbzjxG2ZImHn8tukFisV-bfOEuDfukubd7OX4lnDibIb-lc4ND5vGG8LGOR91i3dVMdm2jVloASXjLMoReGHv_pCrVJ4u2UfOuN5IMehbtjgitmx2LRunCIn2OL4hRlgucE8n_NHzwmdIlzLlWRazd6w7mfYgFZeLbvGcNHutbj5o6ougx-L6L_-G7AU4Y1RUvxA2SiIeqmoPWpyvGXHHES6-s0EVhY8Kc76-TLUHWdKc_NbucWJGQIbcqtFOOJNGy3Fyzak3xVL4YDc8k")
   resp, err := client.Do(req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()

   var result struct {
      VodPlaylistedPlaybackUrls struct {
         Result struct {
            PlaybackUrls struct {
               IntraTitlePlaylist []struct {
                  Urls []struct {
                     Url string
                  }
               }
            }
         }
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      panic(err)
   }
   dirty := result.VodPlaylistedPlaybackUrls.Result.PlaybackUrls.IntraTitlePlaylist[0].Urls[0].Url
   clean, err := DoClean(dirty)
   if err != nil {
      panic(err)
   }
   data, err := get(clean.String())
   if err != nil {
      panic(err)
   }
   if strings.Contains(data, `height="2160"`) {
      log.Print("pass")
   } else {
      log.Print("fail")
   }
}

func DoClean(address string) (*url.URL, error) {
   parsedUrl, err := url.Parse(address)
   if err != nil {
      return nil, err
   }
   parts := strings.Split(parsedUrl.Path, "/")
   // parts[0] is "" (leading slash)
   // parts[1] is "dm"
   // parts[2] is "3$..."
   // parts[3] is "iad_2"
   // parts[4:] is the raw 4K path
   parsedUrl.Path = "/" + strings.Join(parts[4:], "/")
   return parsedUrl, nil
}

func get(address string) (string, error) {
   resp, err := http.Get(address)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   var data strings.Builder
   _, err = io.Copy(&data, resp.Body)
   if err != nil {
      return "", err
   }
   return data.String(), nil
}
