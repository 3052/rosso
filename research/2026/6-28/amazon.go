package main

import (
   "encoding/json"
   "flag"
   "fmt"
   "net/http"
   "net/url"
   "strings"
   "time"
)

// fetchAndTest builds the request, executes it, and unmarshals into a struct to validate
func fetchAndTest(itemID string) error {
   // Wait 1 second before making the request
   time.Sleep(1 * time.Second)
   client := &http.Client{}
   reqURL := &url.URL{
      Scheme: "https",
      Host: "atv-ps.amazon.com",
      Path:   "/lrcedge/getDataByJavaTransform/v1/lr/detailsPage/detailsPageATF",
   }
   // Build Query Parameters
   q := url.Values{}
   q.Add("itemId", itemID)
   q.Add("presentationScheme", "android-tv-react")
   q.Add("deviceID", "uuidb43bee409bd448cfb5ba3337bd241645")
   q.Add("deviceTypeID", "A3NM0WFSU3DLT5") // Swapped to new ID
   q.Add("roles", "playback-envelope-supported")
   reqURL.RawQuery = q.Encode()
   req, err := http.NewRequest("GET", reqURL.String(), nil)
   if err != nil {
      return fmt.Errorf("failed to create request: %w", err)
   }
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

   // Check path: actions
   if len(data.Resource.Actions) > 0 {
      env := data.Resource.Actions[0].Metadata.PlaybackExperienceMetadata.PlaybackEnvelope
      if env != "" {
         foundEnvelope = env
      }
   }

   // Assert we found a non-empty string in the actions path
   if foundEnvelope == "" {
      return fmt.Errorf("playbackEnvelope key is missing or empty in 'actions'")
   }

   // Optional: Output a snippet of the string to confirm
   if len(foundEnvelope) > 20 {
      fmt.Printf("   -> Found playbackEnvelope in 'actions': %s...\n", foundEnvelope[:20])
   } else {
      fmt.Printf("   -> Found playbackEnvelope in 'actions': %s\n", foundEnvelope)
   }

   return nil
}

func main() {
   // Define an empty string flag so no default runs automatically
   region := flag.String("region", "", "Specify the region to test: 'us' or 'gb'")
   flag.Parse()

   if *region == "" {
      fmt.Println("❌ Error: Region flag is required. Please specify a region using -region=us or -region=gb")
      return
   }

   switch strings.ToLower(*region) {
   case "us":
      fmt.Println("Running test for US Item...")
      usItemID := "amzn1.dv.gti.af991753-e4cf-4d28-880d-dfca3d1e8d24"
      if err := fetchAndTest(usItemID); err != nil {
         fmt.Printf("❌ US Test Failed: %v\n\n", err)
      } else {
         fmt.Println("✅ US Test Passed!\n")
      }
   case "gb":
      fmt.Println("Running test for GB Item...")
      gbItemID := "amzn1.dv.gti.775a185a-8920-4711-8dbf-d3791538d5af"
      if err := fetchAndTest(gbItemID); err != nil {
         fmt.Printf("❌ GB Test Failed: %v\n\n", err)
      } else {
         fmt.Println("✅ GB Test Passed!\n")
      }
   default:
      fmt.Printf("❌ Invalid region '%s' specified. Please use -region=us or -region=gb\n", *region)
   }
}

// APIResponse represents the JSON paths we need to validate.
type APIResponse struct {
   Resource struct {
      // Path: resource.actions[0].metadata.playbackExperienceMetadata.playbackEnvelope
      Actions []struct {
         Metadata struct {
            PlaybackExperienceMetadata struct {
               PlaybackEnvelope string `json:"playbackEnvelope"`
            } `json:"playbackExperienceMetadata"`
         } `json:"metadata"`
      } `json:"actions"`
   } `json:"resource"`
}
