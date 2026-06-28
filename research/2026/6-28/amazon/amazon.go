package main

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
   "time"
)

// APIResponse represents the JSON paths we need to validate.
// It accounts for both `actions` and `primaryActions`.
type APIResponse struct {
   Resource struct {
      // Path 1: resource.actions[0].metadata.playbackExperienceMetadata.playbackEnvelope
      Actions []struct {
         Metadata struct {
            PlaybackExperienceMetadata struct {
               PlaybackEnvelope string `json:"playbackEnvelope"`
            } `json:"playbackExperienceMetadata"`
         } `json:"metadata"`
      } `json:"actions"`

      // Path 2: resource.primaryActions[0].navigationAction.playbackMetadata.playbackExperienceMetadata.playbackEnvelope
      PrimaryActions []struct {
         NavigationAction struct {
            PlaybackMetadata struct {
               PlaybackExperienceMetadata struct {
                  PlaybackEnvelope string `json:"playbackEnvelope"`
               } `json:"playbackExperienceMetadata"`
            } `json:"playbackMetadata"`
         } `json:"navigationAction"`
      } `json:"primaryActions"`
   } `json:"resource"`
}

func main() {
   // 1. Call for the US item
   fmt.Println("Running test for US Item...")
   usItemID := "amzn1.dv.gti.af991753-e4cf-4d28-880d-dfca3d1e8d24"
   if err := fetchAndTest("US", usItemID); err != nil {
      fmt.Printf("❌ US Test Failed: %v\n\n", err)
   } else {
      fmt.Println("✅ US Test Passed!\n")
   }

   // 2. Call for the GB item
   fmt.Println("Running test for GB Item...")
   gbItemID := "amzn1.dv.gti.775a185a-8920-4711-8dbf-d3791538d5af"
   if err := fetchAndTest("GB", gbItemID); err != nil {
      fmt.Printf("❌ GB Test Failed: %v\n\n", err)
   } else {
      fmt.Println("✅ GB Test Passed!\n")
   }
}

// fetchAndTest builds the request, executes it, and unmarshals into a struct to validate
func fetchAndTest(geoLocation, itemID string) error {
   // Wait 1 second before making the request
   time.Sleep(1 * time.Second)

   client := &http.Client{}
   reqURL := &url.URL{
      Scheme: "https",
      Host:   "s0s7.api.amazonvideo.com",
      Path:   "/lrcedge/getDataByJavaTransform/v1/lr/detailsPage/detailsPageATF",
   }

   // Build Query Parameters (Restored geoLocation)
   q := url.Values{}
   q.Add("dctrHR", "ZAZ")
   q.Add("gascEnabled", "true")
   q.Add("geoLocation", geoLocation)
   q.Add("dynamicFeatures", "CLIENT_DECORATION_ENABLE_DAAPI,DetailsAtf,RecordingCardSupported")
   q.Add("firmware", "google/sdk_gphone_x86/generic_x86_arm:11/RSR1.240422.006/12134477:userdebug/dev-keys")
   q.Add("model", "sdk_gphone_x86")
   q.Add("transformStore", "local")
   q.Add("uxLocale", "en_US")
   q.Add("supportedLocales", "de_DE,en_US,es_ES,fr_FR,it_IT,nl_NL,pl_PL,pt_BR,pt_PT,ms_MY,fil_PH,hi_IN,ta_IN,te_IN,nb_NO,sv_SE,da_DK,zh_CN,zh_TW,ko_KR,th_TH,fi_FI,tr_TR,id_ID,ru_RU,ja_JP,es_US,el_GR,ro_RO,cs_CZ,hu_HU,es_419,ar_XD,en_XD")
   q.Add("manufacturer", "Google")
   q.Add("transformStage", "prod")
   q.Add("deviceTypeID", "A3NM0WFSU3DLT5") // Swapped to new ID
   q.Add("deviceID", "uuidb43bee409bd448cfb5ba3337bd241645")
   q.Add("widgetScheme", "lrc-detail-v3")
   q.Add("isReactionEnabledInDP", "true")
   q.Add("isMapsV2DatumEnabled", "true")
   q.Add("clientId", "pv-lrc-rust")
   q.Add("javaTransformTimeout", "5000")
   q.Add("timeZoneId", "America/Chicago")
   q.Add("osLocale", "en_US")
   q.Add("tid", "ab8mt4dd97et")
   q.Add("dctrCV", "aHImZmImYmM9MTAwJmRzJjA=")
   q.Add("newMaturityRatings", "true")
   q.Add("roles", "live-supported,startover-supported,linear-supported,prime-offer-supported,multipart-notification-supported,tvod-supported,svod-supported,supports-dai-timeshifting,supports-rapid-recap,playback-envelope-supported,linear-playback-envelope-supported,prime-benefit-activation-supported,supports-stream-selector-modal,av-liveliness-with-unavailable-message-supported,new-title-badge-supported,coming-soon-title-badge-supported,supports-trailer,supports-stream-selector-modal,sponsored-promotion-slot-supported")
   q.Add("operatingSystem", "Android")
   q.Add("nerid", "Z7f6/GH3C7S138TVRFYQQb00")
   q.Add("dynamicBadgingEnabled", "true")
   q.Add("presentationScheme", "android-tv-react")
   q.Add("journeyIngressContext", "8|EgRzdm9k")
   q.Add("clientFeatures", "ShowPSEWarning,EnableBuyBoxV2,EnableRecordingExperience,EnableDetailsRecordButton,EnableBuyBoxActionsV2")
   q.Add("chipset", "goldfish_x86")
   q.Add("dctrEI", "0")
   q.Add("itemId", itemID)

   reqURL.RawQuery = q.Encode()
   req, err := http.NewRequest("GET", reqURL.String(), nil)
   if err != nil {
      return fmt.Errorf("failed to create request: %w", err)
   }

   // Add Headers
   req.Header.Add("Accept-Encoding", "identity")
   req.Header.Add("Host", "s0s7.api.amazonvideo.com")
   req.Header.Add("User-Agent", "Android/google/sdk_gphone_x86/generic_x86_arm:11/RSR1.240422.006/12134477:userdebug/dev-keys, Ignition X/15.5.2026042820-android, Google")
   req.Header.Add("accept", "application/json")
   req.Header.Add("content-type", "application/json")
   req.Header.Add("x-atv-page-type", "ATVDetail")
   req.Header.Add("x-client-app", "avlrc")
   req.Header.Add("x-client-version", "unknown-version")
   req.Header.Add("x-request-priority", "CRITICAL")
   
   // Added Authorization header
   req.Header.Set("Authorization", "Bearer Atna|EwMDIA-tequdXsxFciQc9tBeNpN93XGtmYJKFcKNAOkWchz6ce12IL3Z14B6wRoKfwRW6serUz6PVMJJnL9v0hg6rMQI8VLYMDVMt6US0i-rJ5CW9svjDv_xxajreYn51oWeiNg2gOCfKynqPyXs5OS32HdkUx7roEyBpvYSqQO_c4ECKqYaVRhE-3txwStc3Qe7O4t4OOOvWPdB0i6Mx2CtNA6ubzYqC5_bzlVIaoDA9zngh3jMMoWWk8ICkctoyFw4TgNSE07UiQCXn5ZKuaT_2svdHLgMc4U1YkaLSaHH6JHfDiKlSqWN0RZbAsdSUayl5P3Y_bkOzTf3aY0P1QJvgbqsaSKHQozcTIxYSXblFlCI_Z_dDdaAxHBN1QGDEINUQ4sJMaay_-lS0sBTLy2H0YR54ayGFwcRSSfzV2is56ODNA")

   // Execute Request
   resp, err := client.Do(req)
   if err != nil {
      return fmt.Errorf("request failed: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   // Unmarshal into the struct
   var data APIResponse
   if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
      return fmt.Errorf("failed to decode JSON response: %w", err)
   }

   var foundEnvelope string

   // Check path 1: actions
   if len(data.Resource.Actions) > 0 {
      env := data.Resource.Actions[0].Metadata.PlaybackExperienceMetadata.PlaybackEnvelope
      if env != "" {
         foundEnvelope = env
      }
   }

   // Check path 2 (only if we didn't find it in path 1): primaryActions
   if foundEnvelope == "" && len(data.Resource.PrimaryActions) > 0 {
      env := data.Resource.PrimaryActions[0].NavigationAction.PlaybackMetadata.PlaybackExperienceMetadata.PlaybackEnvelope
      if env != "" {
         foundEnvelope = env
      }
   }

   // Assert we found a non-empty string in at least one of the paths
   if foundEnvelope == "" {
      return fmt.Errorf("playbackEnvelope key is missing or empty in both 'actions' and 'primaryActions'")
   }

   // Optional: Output a snippet of the string to confirm
   if len(foundEnvelope) > 20 {
      fmt.Printf("   -> Found playbackEnvelope: %s...\n", foundEnvelope[:20])
   } else {
      fmt.Printf("   -> Found playbackEnvelope: %s\n", foundEnvelope)
   }

   return nil
}
