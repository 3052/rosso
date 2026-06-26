package amazon

import (
   "encoding/base64"
   "encoding/xml"
   "fmt"
   "io"
   "net/http"
   "os"
   "strings"
   "testing"
   "time"
)

// runDeviceCombinations is a shared helper that iterates through all parameters for a given device.
func runDeviceCombinations(t *testing.T, devName string, keyDir string, drmType string) {
   deviceID := "uuidcbb2f9705f13437e9e515622dce02106"
   titleID := "amzn1.dv.gti.930793f1-4a5c-4998-b335-7150770e5fe0"
   marketplaceID := "ATVPDKIKX0DER"

   playbackEnvBytes, err := os.ReadFile("playback_env.txt")
   if err != nil {
      t.Fatalf("Failed to read playback_env.txt: %v", err)
   }
   playbackEnv := strings.TrimSpace(string(playbackEnvBytes))

   authBearerBytes, err := os.ReadFile("auth_bearer.txt")
   if err != nil {
      t.Fatalf("Failed to read auth_bearer.txt: %v", err)
   }
   authBearer := strings.TrimSpace(string(authBearerBytes))

   if authBearer == "" || titleID == "" || playbackEnv == "" {
      t.Fatal("Missing credentials or title info")
   }

   client := NewClient(&http.Client{})

   hdcpLevel := "2.3" // Restricted to 2.3 only as requested
   resolutions := []string{"576p", "2160p"}
   hdrFormats := []string{"None", "HDR10", "DolbyVision"}
   bitrates := []string{"CVBR", "CBR"}

   fmt.Printf("\n=======================================================\n")
   fmt.Printf("Testing Device: %s\n", devName)
   fmt.Printf("=======================================================\n")

   for _, res := range resolutions {
      for _, hdr := range hdrFormats {
         for _, bitrate := range bitrates {

            // Removed HDCP from the display output name, locked to H265, added bitrate
            comboName := fmt.Sprintf("Res:%s HDR:%s Codec:H265 Bitrate:%s", res, hdr, bitrate)

            // Add a 1 second delay before each combination
            time.Sleep(1 * time.Second)

            profile := DeviceProfile{
               DeviceID:          deviceID,
               AuthBearer:        authBearer,
               DRMType:           drmType,
               HDCPLevel:         hdcpLevel,
               MaxResolution:     res,
               HDRFormats:        hdr,
               VideoCodec:        "H265", // Passing only H265 now
               BitrateAdaptation: bitrate,
            }

            // 1. Get Manifest
            mpdURL, err := client.GetManifest(profile, titleID, marketplaceID, playbackEnv)
            if err != nil {
               fmt.Printf("[FAIL Manifest] %s -> %v\n", comboName, err)
               continue
            }

            // 2. Download MPD
            resp, err := http.Get(mpdURL)
            if err != nil {
               fmt.Printf("[FAIL Download] %s -> %v\n", comboName, err)
               continue
            }

            mpdData, err := io.ReadAll(resp.Body)
            resp.Body.Close()
            if err != nil {
               fmt.Printf("[FAIL Read MPD] %s -> %v\n", comboName, err)
               continue
            }

            // 3. Parse MPD
            var manifest mpdXML
            if err := xml.Unmarshal(mpdData, &manifest); err != nil {
               fmt.Printf("[FAIL Parse MPD] %s -> %v\n", comboName, err)
               continue
            }

            var highestRep *representationXML
            var activeAdp *adaptationSetXML

            for _, period := range manifest.Periods {
               for _, adp := range period.AdaptationSets {
                  if adp.ContentType == "video" || adp.MimeType == "video/mp4" {
                     for _, rep := range adp.Representations {
                        repCopy := rep
                        // Selecting the HIGHEST bandwidth
                        if highestRep == nil || repCopy.Bandwidth > highestRep.Bandwidth {
                           highestRep = &repCopy
                           adpCopy := adp // Safe copy to prevent Go <1.22 pointer referencing bugs
                           activeAdp = &adpCopy
                        }
                     }
                  }
               }
            }

            if highestRep == nil {
               fmt.Printf("[FAIL No Video] %s -> No video representations found\n", comboName)
               continue
            }

            // Print Width, Height, and Codecs for the selected stream
            fmt.Printf("[INFO] %s -> Selected highest representation -> Width: %d, Height: %d, Codecs: %s\n",
               comboName, highestRep.Width, highestRep.Height, highestRep.Codecs)

            if activeAdp == nil || len(activeAdp.ContentProtections) == 0 {
               fmt.Printf("[FAIL No ContentProtection] %s -> No ContentProtection found on AdaptationSet\n", comboName)
               continue
            }

            // 4. Extract PSSH / Init Data directly from AdaptationSet
            var initDataB64 string
            for _, cp := range activeAdp.ContentProtections {
               if drmType == "Widevine" && strings.EqualFold(cp.SchemeIdUri, "urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed") {
                  initDataB64 = cp.Pssh
                  break
               }
               if drmType == "PlayReady" && strings.EqualFold(cp.SchemeIdUri, "urn:uuid:9a04f079-9840-4286-ab92-e65be0885f95") {
                  if cp.Pro != "" {
                     initDataB64 = cp.Pro
                  } else {
                     initDataB64 = cp.Pssh
                  }
                  break
               }
            }

            if initDataB64 == "" {
               fmt.Printf("[FAIL No InitData] %s -> Could not find %s init data\n", comboName, drmType)
               continue
            }

            initDataBytes, err := base64.StdEncoding.DecodeString(initDataB64)
            if err != nil {
               fmt.Printf("[FAIL Decode InitData] %s -> %v\n", comboName, err)
               continue
            }

            // 5. Generate CDM Challenge
            challengeBytes, err := generateCDMChallenge(drmType, keyDir, initDataBytes)
            if err != nil {
               fmt.Printf("[FAIL Challenge] %s -> %v\n", comboName, err)
               continue
            }

            // 6. Make License Request
            licenseB64, err := client.GetLicense(profile, playbackEnv, challengeBytes)
            if err != nil {
               fmt.Printf("[FAIL License] %s -> %v\n", comboName, err)
               continue
            }

            // If we get here, it worked!
            fmt.Printf("[SUCCESS] %s -> License length: %d\n", comboName, len(licenseB64))
         }
      }
   }
}
