// --- get_playback_envelope.go ---
// Fetches title details and extracts the playbackEnvelope required for the PRS request.
package amazon

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
)

// Updated struct to correctly map resource -> header -> titleActionsV2
type TitleDetailResponse struct {
   Resource struct {
      Header struct {
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
      } `json:"header"`
   } `json:"resource"`
}

func GetPlaybackEnvelope(titleID, deviceID, bearerToken string) (string, error) {
   baseURL := "https://abzq7aq4866p.na.api.amazonvideo.com/cdp/switchblade/android/getDataByJvmTransform/v1/dv-android/detail/vod/v1.kt"

   q := url.Values{}
   q.Add("capabilities", "")
   q.Add("clientName", "ATVAndroidThirdPartyClient")
   q.Add("contentType", "VOD")
   q.Add("deviceId", deviceID)
   q.Add("deviceTypeID", "A43PXU4ZN2AL1")
   q.Add("featureScheme", "mobile-android-features-v11.1")
   q.Add("fetchType", "FETCH_FROM_BUTTON_CLICK")
   q.Add("firmware", "fmw:30-app:3.0.458.357")
   q.Add("format", "json")
   q.Add("isChannelPurchasingEnabled", "false")
   q.Add("isGatedVamEnabled", "true")
   q.Add("isGeneratedRequest", "false")
   q.Add("isPlaybackEnvelopeSupported", "true")
   q.Add("isPrimePurchasingEnabled", "false")
   q.Add("isPurchaseWorkflowV2Enabled", "true")
   q.Add("isReactionEnabled", "true")
   q.Add("isSwift2p7Capable", "false")
   q.Add("isTVODPurchasingEnabled", "false")
   q.Add("isWatchModalEnabled", "true")
   q.Add("itemId", titleID)
   q.Add("osLocale", "en_US")
   q.Add("overridePCON", "false")
   q.Add("priorityLevel", "2")
   q.Add("screenDensity", "DEFAULT")
   q.Add("screenWidth", "sw360dp")
   q.Add("softwareVersion", "458")
   q.Add("supportsAnyXbdLinkAction", "true")
   q.Add("supportsConsentRedirection", "true")
   q.Add("supportsDaiTimeShifting", "true")
   q.Add("supportsDvrRecording", "DISABLED")
   q.Add("supportsMAPSLiveBadging", "false")
   q.Add("supportsMseEventLevelOffers", "true")
   q.Add("supportsMultiSourcedEvents", "true")
   q.Add("supportsMultiSourcedEventsDynamicGating", "T2")
   q.Add("supportsPaymentStatus", "true")
   q.Add("supportsPKMZ", "false")
   q.Add("supportsPreorderModalMessaging", "true")
   q.Add("supportsRapidRecap", "true")
   q.Add("supportsSeeMoreIngressNodes", "true")
   q.Add("supportsStreamSelectorModal", "true")
   q.Add("supportsTitleLifeCycleComingSoonEpisodes", "true")
   q.Add("supportsTitleLifeCycleTrailer", "true")
   q.Add("supportsTitleMetadataBadging", "true")
   q.Add("swiftPriorityLevel", "critical")
   q.Add("timeZoneId", "America/Chicago")
   q.Add("useMessagePresentationV2", "true")
   q.Add("uxLocale", "en_US")
   q.Add("version", "1")
   q.Add("widgetSchemeVersion", "1")

   req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", baseURL, q.Encode()), nil)
   if err != nil {
      return "", err
   }

   req.Header.Set("x-gasc-enabled", "true")
   req.Header.Set("x-request-priority", "CRITICAL")
   req.Header.Set("x-atv-page-id", titleID)
   req.Header.Set("x-atv-page-type", "ATVDetail")
   req.Header.Set("Accept", "application/json")
   req.Header.Set("Accept-Language", "en_US")
   req.Header.Set("User-Agent", "Dalvik/2.1.0 (Linux; U; Android 11; sdk_gphone_x86_64 Build/RSR1.240422.006)")
   req.Header.Set("x-retry-count", "0")
   req.Header.Set("Authorization", "Bearer "+bearerToken)
   req.Header.Set("Connection", "Keep-Alive")
   req.Header.Set("Accept-Encoding", "identity")

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

   // Updated path to include .Header.
   items := detailResp.Resource.Header.TitleActionsV2.BuyBoxActionsView.PlaybackGroup.Items
   if len(items) > 0 {
      envelope := items[0].ItemReference.PlaybackExperienceMetadata.PlaybackEnvelope
      if envelope != "" {
         return envelope, nil
      }
   }

   return "", fmt.Errorf("playbackEnvelope not found in title details response")
}
