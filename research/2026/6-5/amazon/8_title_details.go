package amazon

import (
   "encoding/json"
   "fmt"
   "net/http"
)

func GetPlaybackEnvelope(actorAccessToken, titleId, deviceId string) (string, error) {
   url := fmt.Sprintf("https://abzq7aq4866p.na.api.amazonvideo.com/cdp/switchblade/android/getDataByJvmTransform/v1/dv-android/detail/vod/v1.kt?itemId=%s&deviceId=%s&deviceTypeID=A43PXU4ZN2AL1&format=json&version=1", titleId, deviceId)

   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      return "", err
   }

   req.Header.Set("Authorization", "Bearer "+actorAccessToken)
   req.Header.Set("User-Agent", "Dalvik/2.1.0 (Linux; U; Android 11; sdk_gphone_x86_64 Build/RSR1.240422.006)")

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   var data map[string]interface{}
   if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
      return "", err
   }

   // Amazon's response structure is deeply nested and dynamic; recursive search is the safest way to find the envelope
   envelope := findKeyRecursively(data, "playbackEnvelope")
   if envelope == "" {
      return "", fmt.Errorf("playbackEnvelope not found in title details response")
   }

   return envelope, nil
}

func findKeyRecursively(data interface{}, targetKey string) string {
   switch v := data.(type) {
   case map[string]interface{}:
      for key, val := range v {
         if key == targetKey {
            if str, ok := val.(string); ok {
               return str
            }
         }
         if res := findKeyRecursively(val, targetKey); res != "" {
            return res
         }
      }
   case []interface{}:
      for _, item := range v {
         if res := findKeyRecursively(item, targetKey); res != "" {
            return res
         }
      }
   }
   return ""
}
