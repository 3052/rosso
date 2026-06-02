package amazon

import (
   "fmt"
   "os"
   "strings"
   "testing"
)

// 3. Run this to grab an MPD after you have generated the token in Step 2.
func TestStep3_GetMPD(t *testing.T) {
   tokenFile := getTempTokenPath()
   tokenBytes, err := os.ReadFile(tokenFile)
   if err != nil || len(tokenBytes) == 0 {
      t.Fatalf("Failed to read token from %s. Please run TestStep1 and TestStep2 first.", tokenFile)
   }
   accessToken := strings.TrimSpace(string(tokenBytes))

   asin := "B075RND57T"

   opts := DefaultPlaybackOptions()
   // DOCUMENTATION: Use "HD" or "UHD" here if you have an L1 CDM.
   // opts.VideoQuality = "HD"
   opts.VideoQuality = "SD" // Swapped to SD for L3 CDMs
   opts.VideoCodec = "H264"
   opts.BitrateMode = "CVBR,CBR"

   t.Logf("Fetching manifest for ASIN %s...", asin)
   manifestResp, err := GetPlaybackResources(
      playbackEndpoint,
      accessToken,
      asin,
      marketplaceIDUS,
      defaultDevice,
      opts,
   )
   if err != nil {
      t.Fatalf("Failed to fetch playback resources: %v", err)
   }

   rawMPD, err := GetBestMPDURL(manifestResp)
   if err != nil {
      t.Fatalf("Failed to get best MPD URL: %v", err)
   }

   cleanMPD := CleanMPDURL(rawMPD)

   fmt.Println()
   fmt.Println("=====================================================")
   fmt.Println("SUCCESS! Final DASH Manifest (MPD):")
   fmt.Printf("%s\n", cleanMPD)
   fmt.Println("=====================================================")
}
