package amazon

import (
   "encoding/json"
   "fmt"
   "net/http"
)

// GetItemDetails uses the actor access token to get metadata for a specific title, returning the playback envelope.
func GetItemDetails(actorAccessToken, titleId string) (string, error) {
   url := "https://s0s7.api.amazonvideo.com/lrcedge/getDataByJavaTransform/v1/lr/detailsPage/detailsPageATF"

   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      return "", err
   }

   q := req.URL.Query()
   q.Add("itemId", titleId)
   q.Add("deviceTypeID", "A2SNKIF736WF4T")
   q.Add("deviceID", "uuidcbb2f9705f13437e9e515622dce02106")
   q.Add("firmware", "google/sdk_gphone_x86/generic_x86_arm:11/RSR1.240422.006/12134477:userdebug/dev-keys")
   q.Add("manufacturer", "Google")
   q.Add("chipset", "goldfish_x86")
   q.Add("model", "sdk_gphone_x86")
   q.Add("operatingSystem", "Android")
   q.Add("clientId", "pv-lrc-rust")
   req.URL.RawQuery = q.Encode()

   req.Header.Set("User-Agent", "Android/google/sdk_gphone_x86/generic_x86_arm:11/RSR1.240422.006/12134477:userdebug/dev-keys, Ignition X/15.5.2026042820-android, Google")
   req.Header.Set("Authorization", "Bearer "+actorAccessToken)
   req.Header.Set("Accept", "application/json")
   req.Header.Set("x-client-app", "avlrc")

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   var result struct {
      Resource struct {
         SucceededItems map[string]struct {
            PlaybackEnvelope string `json:"playbackEnvelope"`
         } `json:"succeededItems"`
      } `json:"resource"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return "", err
   }

   item, exists := result.Resource.SucceededItems[titleId]
   if !exists || item.PlaybackEnvelope == "" {
      return "", fmt.Errorf("playbackEnvelope not found for titleId: %s", titleId)
   }

   return item.PlaybackEnvelope, nil
}
