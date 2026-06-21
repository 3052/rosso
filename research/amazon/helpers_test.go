package amazon

import (
   "encoding/base64"
   "encoding/xml"
   "fmt"
   "io"
   "net/http"
   "os"
   "path/filepath"
   "strings"
   "testing"

   "41.neocities.org/diana/playReady"
   "41.neocities.org/diana/widevine"
)

type mpdXML struct {
   Periods []periodXML `xml:"Period"`
}

type periodXML struct {
   AdaptationSets []adaptationSetXML `xml:"AdaptationSet"`
}

type adaptationSetXML struct {
   ContentType        string              `xml:"contentType,attr"`
   MimeType           string              `xml:"mimeType,attr"`
   ContentProtections []contentProtXML    `xml:"ContentProtection"`
   Representations    []representationXML `xml:"Representation"`
}

type representationXML struct {
   ID                 string           `xml:"id,attr"`
   Bandwidth          int              `xml:"bandwidth,attr"`
   ContentProtections []contentProtXML `xml:"ContentProtection"`
}

type contentProtXML struct {
   SchemeIdUri string `xml:"schemeIdUri,attr"`
   Pssh        string `xml:"pssh"` // Widevine urn:mpeg:cenc:2013
   Pro         string `xml:"pro"`  // PlayReady urn:microsoft:playready
}

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

   keySchemes := []string{"", "SingleKey", "DualKey"}
   hdcpLevels := []string{"1.4", "2.2", "2.3"}
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

// generateCDMChallenge generates the license challenge using the local diana DRM packages.
func generateCDMChallenge(drmType string, keyDir string, initData []byte) ([]byte, error) {
   if drmType == "Widevine" {
      pssh, err := widevine.DecodePsshData(initData)
      if err != nil {
         return nil, fmt.Errorf("failed to decode widevine pssh: %w", err)
      }

      clientIDPath := filepath.Join(keyDir, "device_client_id_blob")
      clientID, err := os.ReadFile(clientIDPath)
      if err != nil {
         return nil, fmt.Errorf("failed to read %s: %w", clientIDPath, err)
      }

      privKeyPath := filepath.Join(keyDir, "device_private_key")
      privKeyBytes, err := os.ReadFile(privKeyPath)
      if err != nil {
         return nil, fmt.Errorf("failed to read %s: %w", privKeyPath, err)
      }

      privKey, err := widevine.DecodePrivateKey(privKeyBytes)
      if err != nil {
         return nil, fmt.Errorf("failed to decode private key: %w", err)
      }

      reqData, err := pssh.EncodeLicenseRequest(clientID)
      if err != nil {
         return nil, fmt.Errorf("failed to encode license request: %w", err)
      }

      challenge, err := widevine.EncodeSignedMessage(reqData, privKey)
      if err != nil {
         return nil, fmt.Errorf("failed to sign message: %w", err)
      }

      return challenge, nil

   } else if drmType == "PlayReady" {
      wrm, err := playReady.ParsePro(initData)
      if err != nil {
         return nil, fmt.Errorf("failed to parse playready PRO: %w", err)
      }

      bcertPath := filepath.Join(keyDir, "bdevcert.dat")
      chainBytes, err := os.ReadFile(bcertPath)
      if err != nil {
         return nil, fmt.Errorf("failed to read %s: %w", bcertPath, err)
      }

      chain, err := playReady.ParseChain(chainBytes)
      if err != nil {
         return nil, fmt.Errorf("failed to parse chain: %w", err)
      }

      privKeyPath := filepath.Join(keyDir, "zprivsig.dat")
      privKeyBytes, err := os.ReadFile(privKeyPath)
      if err != nil {
         return nil, fmt.Errorf("failed to read %s: %w", privKeyPath, err)
      }

      signingKey, err := playReady.ParseRawPrivateKey(privKeyBytes)
      if err != nil {
         return nil, fmt.Errorf("failed to parse private key: %w", err)
      }

      kid := []byte(wrm.Data.Kid)

      var contentID string
      if wrm.Data.CustomAttributes != nil {
         contentID = wrm.Data.CustomAttributes.ContentId
      }

      challenge, err := chain.LicenseRequestBytes(signingKey, kid, contentID)
      if err != nil {
         return nil, fmt.Errorf("failed to generate PR license request: %w", err)
      }

      return challenge, nil
   }

   return nil, fmt.Errorf("unsupported DRM type: %s", drmType)
}
