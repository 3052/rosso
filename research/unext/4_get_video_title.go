package unext

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

// GetVideoTitle fetches the title metadata and returns the episodeCode
func GetVideoTitle(client *http.Client, titleID string) (string, error) {
   baseURL := "https://cc.unext.jp/"

   params := url.Values{}
   params.Set("zxuid", "2cd3deff87a0")
   params.Set("zxemp", "29719881")
   params.Set("operationName", "cosmo_getVideoTitle")
   params.Set("variables", fmt.Sprintf(`{"code":"%s"}`, titleID))
   params.Set("extensions", `{"persistedQuery":{"version":1,"sha256Hash":"0295df1eacb9e942a2c96cb4f1e5f47c3ac96f2bc50589d167e4708b6b701bbd"}}`)

   req, err := http.NewRequest("GET", baseURL+"?"+params.Encode(), nil)
   if err != nil {
      return "", err
   }

   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   req.Header.Set("Accept", "*/*")
   req.Header.Set("Accept-Language", "en-US,en;q=0.5")
   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("Apollographql-Client-Name", "cosmo")
   req.Header.Set("Apollographql-Client-Version", "v126.0-prod-017e302")
   req.Header.Set("Origin", "https://video.unext.jp")
   req.Header.Set("Referer", fmt.Sprintf("https://video.unext.jp/title/%s", titleID))

   resp, err := client.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return "", err
   }

   var titleData VideoTitleResponse
   if err := json.Unmarshal(body, &titleData); err != nil {
      return "", err
   }

   if titleData.Data.WebfrontTitleStage.KeyEpisodes.Current.ID == "" {
      return "", fmt.Errorf("failed to extract episodeCode from response: %s", string(body))
   }

   return titleData.Data.WebfrontTitleStage.KeyEpisodes.Current.ID, nil
}

type VideoTitleResponse struct {
   Data struct {
      WebfrontTitleStage struct {
         KeyEpisodes struct {
            Current struct {
               ID string `json:"id"`
            } `json:"current"`
         } `json:"keyEpisodes"`
      } `json:"webfront_title_stage"`
   } `json:"data"`
}
