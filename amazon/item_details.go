package amazon

import (
   "encoding/json"
   "fmt"
   "net/http"
)

// ItemDetails contains metadata for a specific title, including the playback
// envelope
type ItemDetails struct {
   PlaybackEnvelope string `json:"playbackEnvelope"`
}

func (*ItemDetails) CachePath() string {
   return "rosso/amazon/ItemDetails"
}

// GetItemDetails uses the actor access token to get metadata for a specific title.
// It explicitly passes UI schema flags to ensure the server returns the PlaybackEnvelope.
func GetItemDetails(actorAccessToken, titleId string) (*ItemDetails, error) {
   url := "https://s0s7.api.amazonvideo.com/lrcedge/getDataByJavaTransform/v1/lr/detailsPage/detailsPageATF"

   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      return nil, err
   }

   q := req.URL.Query()
   q.Add("itemId", titleId)

   // Critical UI and Feature flags to force the V2/V3 BuyBox response with PlaybackEnvelope
   q.Add("widgetScheme", "lrc-detail-v3")
   q.Add("isReactionEnabledInDP", "true")
   q.Add("newMaturityRatings", "true")
   q.Add("dynamicBadgingEnabled", "true")
   q.Add("isMapsV2DatumEnabled", "true")
   q.Add("roles", "live-supported,startover-supported,linear-supported,prime-offer-supported,multipart-notification-supported,tvod-supported,svod-supported,supports-dai-timeshifting,supports-rapid-recap,playback-envelope-supported,linear-playback-envelope-supported,prime-benefit-activation-supported,supports-stream-selector-modal,av-liveliness-with-unavailable-message-supported,new-title-badge-supported,coming-soon-title-badge-supported,supports-consent-redirection,supports-trailer")
   q.Add("presentationScheme", "android-tv-react")
   q.Add("dynamicFeatures", "CLIENT_DECORATION_ENABLE_DAAPI,DetailsAtf,RecordingCardSupported")
   q.Add("clientFeatures", "ShowPSEWarning,EnableBuyBoxV2,EnableRecordingExperience,EnableBuyBoxActionsV2")
   q.Add("transformStore", "local")
   q.Add("transformStage", "prod")
   q.Add("uxLocale", "en_US")
   q.Add("clientId", "pv-lrc-rust")

   // Device parameters
   q.Add("deviceTypeID", "A2SNKIF736WF4T")
   q.Add("deviceID", "uuidcbb2f9705f13437e9e515622dce02106")
   q.Add("firmware", "google/sdk_gphone_x86/generic_x86_arm:11/RSR1.240422.006/12134477:userdebug/dev-keys")
   q.Add("manufacturer", "Google")
   q.Add("chipset", "goldfish_x86")
   q.Add("model", "sdk_gphone_x86")
   q.Add("operatingSystem", "Android")

   req.URL.RawQuery = q.Encode()

   req.Header.Set("User-Agent", "Android/google/sdk_gphone_x86/generic_x86_arm:11/RSR1.240422.006/12134477:userdebug/dev-keys, Ignition X/15.5.2026042820-android, Google")
   req.Header.Set("Authorization", "Bearer "+actorAccessToken)
   req.Header.Set("Accept", "application/json")
   req.Header.Set("x-client-app", "avlrc")
   req.Header.Set("x-atv-page-type", "ATVDetail")

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   // Embed our new ItemDetails struct into the anonymous decoder struct
   var result struct {
      Resource struct {
         PrimaryActions []struct {
            NavigationAction struct {
               PlaybackMetadata struct {
                  PlaybackExperienceMetadata ItemDetails `json:"playbackExperienceMetadata"`
               } `json:"playbackMetadata"`
            } `json:"navigationAction"`
         } `json:"primaryActions"`
      } `json:"resource"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   for _, action := range result.Resource.PrimaryActions {
      details := action.NavigationAction.PlaybackMetadata.PlaybackExperienceMetadata
      if details.PlaybackEnvelope != "" {
         return &details, nil
      }
   }

   return nil, fmt.Errorf("playbackEnvelope not found in primaryActions for titleId: %s", titleId)
}
