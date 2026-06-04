// --- get_title_details.go ---
package amazon

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
)

// TitleDetailResponse represents the structure needed to extract the playbackEnvelope
type TitleDetailResponse struct {
   Resource struct {
      TitleActionsV2 struct {
         BuyBoxActionsView struct {
            PlaybackGroup struct {
               Items []struct {
                  ItemReference struct {
                     PlaybackExperienceMetadata struct {
                        PlaybackEnvelope string `json:"playbackEnvelope"`
                     } `json:"playbackExperienceMetadata"`
                  } `json:"itemReference"`
               } `json:"items"`
            } `json:"playbackGroup"`
         } `json:"buyBoxActionsView"`
      } `json:"titleActionsV2"`
   } `json:"resource"`
}

// GetPlaybackEnvelope fetches the title details and extracts the playback envelope required for PRS
func GetPlaybackEnvelope(titleID, deviceID, bearerToken string) (string, error) {
   baseURL := "https://abzq7aq4866p.na.api.amazonvideo.com/cdp/switchblade/android/getDataByJvmTransform/v1/dv-android/detail/vod/v1.kt"

   q := url.Values{}
   q.Add("clientName", "ATVAndroidThirdPartyClient")
   q.Add("contentType", "VOD")
   q.Add("deviceId", deviceID)
   q.Add("deviceTypeID", "A43PXU4ZN2AL1")
   q.Add("format", "json")
   q.Add("isPlaybackEnvelopeSupported", "true")
   q.Add("itemId", titleID)
   q.Add("osLocale", "en_US")
   q.Add("version", "1")
   // Additional required parameters from the capture
   q.Add("firmware", "fmw:30-app:3.0.458.357")
   q.Add("softwareVersion", "458")

   req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", baseURL, q.Encode()), nil)
   if err != nil {
      return "", err
   }

   req.Header.Set("Authorization", "Bearer "+bearerToken)
   req.Header.Set("Accept", "application/json")
   req.Header.Set("User-Agent", "Dalvik/2.1.0 (Linux; U; Android 11; sdk_gphone_x86_64 Build/RSR1.240422.006)")

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return "", fmt.Errorf("failed to get title details, status: %d", resp.StatusCode)
   }

   var detailResp TitleDetailResponse
   if err := json.NewDecoder(resp.Body).Decode(&detailResp); err != nil {
      return "", err
   }

   // Safely navigate the JSON structure to find the envelope
   items := detailResp.Resource.TitleActionsV2.BuyBoxActionsView.PlaybackGroup.Items
   if len(items) > 0 {
      envelope := items[0].ItemReference.PlaybackExperienceMetadata.PlaybackEnvelope
      if envelope != "" {
         return envelope, nil
      }
   }

   return "", fmt.Errorf("playbackEnvelope not found in title details response")
}
