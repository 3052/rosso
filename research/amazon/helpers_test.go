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

   keySchemes := []string{"SingleKey", "DualKey"}
   hdcpLevels := []string{"2.1", "2.3"}
   resolutions := []string{"480p", "720p", "1080p", "1440p", "2160p"}
   hdrFormats := [][]string{
      {"None"},
      {"HDR10"},
      {"DolbyVision"},
      {"HDR10", "DolbyVision"},
   }

   fmt.Printf("\n=======================================================\n")
   fmt.Printf("Testing Device: %s\n", devName)
   fmt.Printf("=======================================================\n")

   for _, scheme := range keySchemes {
      for _, hdcp := range hdcpLevels {
         for _, res := range resolutions {
            for _, hdr := range hdrFormats {

               comboName := fmt.Sprintf("Scheme:%s HDCP:%s Res:%s HDR:%v", scheme, hdcp, res, hdr)

               // Add a 1 second delay before each combination
               time.Sleep(1 * time.Second)

               profile := DeviceProfile{
                  DeviceID:      deviceID,
                  AuthBearer:    authBearer,
                  DRMType:       drmType,
                  DRMKeyScheme:  scheme,
                  HDCPLevel:     hdcp,
                  MaxResolution: res,
                  HDRFormats:    hdr,
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

               var lowestRep *representationXML
               var activeAdp *adaptationSetXML

               for _, period := range manifest.Periods {
                  for _, adp := range period.AdaptationSets {
                     if adp.ContentType == "video" || adp.MimeType == "video/mp4" {
                        for _, rep := range adp.Representations {
                           repCopy := rep
                           if lowestRep == nil || repCopy.Bandwidth < lowestRep.Bandwidth {
                              lowestRep = &repCopy
                              activeAdp = &adp
                           }
                        }
                     }
                  }
               }

               if lowestRep == nil {
                  fmt.Printf("[FAIL No Video] %s -> No video representations found\n", comboName)
                  continue
               }

               // 4. Extract PSSH / Init Data
               var initDataB64 string
               prots := lowestRep.ContentProtections
               if len(prots) == 0 && activeAdp != nil {
                  prots = activeAdp.ContentProtections
               }

               for _, cp := range prots {
                  if drmType == "Widevine" && cp.SchemeIdUri == "urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed" {
                     initDataB64 = cp.Pssh
                     break
                  }
                  if drmType == "PlayReady" && cp.SchemeIdUri == "urn:uuid:9a04f079-9840-4286-ab92-e65be0885f95" {
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
}
