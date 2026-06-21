package amazon

import (
   "encoding/base64"
   "encoding/xml"
   "fmt"
   "io"
   "net/http"
   "os"
   "path/filepath"
   "testing"

   "41.neocities.org/diana/playReady"
   "41.neocities.org/diana/widevine"
)

// mpdXML, periodXML, etc. are used to parse the DASH manifest to find the lowest quality video PSSH
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

func TestAmazonFlow(t *testing.T) {
   // -------------------------------------------------------------------------
   // Setup: Provide actual extracted values to run the live tests
   // -------------------------------------------------------------------------
   deviceID := ""      // e.g. "uuidb43bee409bd448cfb5ba3337bd241645"
   authBearer := ""    // e.g. "Atna|EwMDICIPxLGAmnVlZgnFhnKMSRVvjHua..."
   titleID := ""       // e.g. "amzn1.dv.gti.af991753-e4cf-4d28-880d-dfca3d1e8d24"
   marketplaceID := "" // e.g. "ATVPDKIKX0DER"
   playbackEnv := ""   // e.g. "MDJ8Cm0KBHBlbnYSJGI1YWQ0MjdhLTIyY2MtN..."

   if authBearer == "" || titleID == "" {
      t.Fatal("Missing credentials or title info")
   }

   client := NewClient(&http.Client{})

   // Define the 3 devices you want to test
   tests := []struct {
      Name    string
      KeyDir  string
      Profile DeviceProfile
   }{
      {
         Name:   "Widevine L3",
         KeyDir: `C:\Users\Steven\AppData\Local\L3`,
         Profile: DeviceProfile{
            DeviceID:      deviceID,
            AuthBearer:    authBearer,
            DRMType:       "Widevine",
            HDCPLevel:     "1.4",
            MaxResolution: "1080p",
            HDRFormats:    []string{"None"},
         },
      },
      {
         Name:   "PlayReady SL2000",
         KeyDir: `C:\Users\Steven\AppData\Local\SL2000`,
         Profile: DeviceProfile{
            DeviceID:      deviceID,
            AuthBearer:    authBearer,
            DRMType:       "PlayReady",
            HDCPLevel:     "1.4",
            MaxResolution: "1080p",
            HDRFormats:    []string{"None"},
         },
      },
      {
         Name:   "PlayReady SL3000",
         KeyDir: `C:\Users\Steven\AppData\Local\SL3000`,
         Profile: DeviceProfile{
            DeviceID:      deviceID,
            AuthBearer:    authBearer,
            DRMType:       "PlayReady",
            HDCPLevel:     "2.3",
            MaxResolution: "2160p",
            HDRFormats:    []string{"HDR10", "DolbyVision"},
         },
      },
   }

   for _, tc := range tests {
      t.Run(tc.Name, func(t *testing.T) {
         fmt.Printf("\n--- Testing %s ---\n", tc.Name)

         // 1. Get the MPD (Optimized for best quality available per profile)
         mpdURL, err := client.GetManifest(tc.Profile, titleID, marketplaceID, playbackEnv)
         if err != nil {
            t.Fatalf("Failed to get manifest: %v", err)
         }
         fmt.Printf("MPD URL: %s\n", mpdURL)

         // 2. Download the MPD
         resp, err := http.Get(mpdURL)
         if err != nil {
            t.Fatalf("Failed to download MPD: %v", err)
         }
         defer resp.Body.Close()

         mpdData, err := io.ReadAll(resp.Body)
         if err != nil {
            t.Fatalf("Failed to read MPD: %v", err)
         }

         // 3. Parse the MPD to find the lowest quality video representation
         var manifest mpdXML
         if err := xml.Unmarshal(mpdData, &manifest); err != nil {
            t.Fatalf("Failed to parse MPD XML: %v", err)
         }

         var lowestRep *representationXML
         var activeAdp *adaptationSetXML

         for _, period := range manifest.Periods {
            for _, adp := range period.AdaptationSets {
               if adp.ContentType == "video" || adp.MimeType == "video/mp4" {
                  for _, rep := range adp.Representations {
                     repCopy := rep // Capture loop variable
                     if lowestRep == nil || repCopy.Bandwidth < lowestRep.Bandwidth {
                        lowestRep = &repCopy
                        activeAdp = &adp // Keep track of parent in case DRM is at AdaptationSet level
                     }
                  }
               }
            }
         }

         if lowestRep == nil {
            t.Fatalf("No video representations found in MPD")
         }
         fmt.Printf("Selected lowest quality video: ID=%s, Bandwidth=%d\n", lowestRep.ID, lowestRep.Bandwidth)

         // 4. Extract PSSH / Init Data
         var initDataB64 string

         // Check Representation level first, fallback to AdaptationSet level
         prots := lowestRep.ContentProtections
         if len(prots) == 0 && activeAdp != nil {
            prots = activeAdp.ContentProtections
         }

         for _, cp := range prots {
            if tc.Profile.DRMType == "Widevine" && cp.SchemeIdUri == "urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed" {
               initDataB64 = cp.Pssh
               break
            }
            if tc.Profile.DRMType == "PlayReady" && cp.SchemeIdUri == "urn:uuid:9a04f079-9840-4286-ab92-e65be0885f95" {
               // Amazon PlayReady MPDs usually carry the <mspr:pro> tag
               if cp.Pro != "" {
                  initDataB64 = cp.Pro
               } else {
                  initDataB64 = cp.Pssh
               }
               break
            }
         }

         if initDataB64 == "" {
            t.Fatalf("Could not find %s init data in the MPD", tc.Profile.DRMType)
         }

         initDataBytes, err := base64.StdEncoding.DecodeString(initDataB64)
         if err != nil {
            t.Fatalf("Failed to decode init data: %v", err)
         }

         // 5. Pass init data to local CDM to generate challenge
         challengeBytes, err := generateCDMChallenge(tc.Profile.DRMType, tc.KeyDir, initDataBytes)
         if err != nil {
            t.Fatalf("Failed to generate CDM challenge: %v", err)
         }

         // 6. Make License Request
         licenseB64, err := client.GetLicense(tc.Profile, playbackEnv, challengeBytes)
         if err != nil {
            t.Fatalf("Failed to get license: %v", err)
         }

         fmt.Printf("Successfully got %s license! Length: %d\n", tc.Profile.DRMType, len(licenseB64))
      })
   }
}

// generateCDMChallenge generates the license challenge using the local diana DRM packages.
func generateCDMChallenge(drmType string, keyDir string, initData []byte) ([]byte, error) {
   if drmType == "Widevine" {
      // 1. Decode PSSH
      pssh, err := widevine.DecodePsshData(initData)
      if err != nil {
         return nil, fmt.Errorf("failed to decode widevine pssh: %w", err)
      }

      // 2. Load device credentials
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

      // 3. Generate License Request
      reqData, err := pssh.EncodeLicenseRequest(clientID)
      if err != nil {
         return nil, fmt.Errorf("failed to encode license request: %w", err)
      }

      // 4. Sign Request
      challenge, err := widevine.EncodeSignedMessage(reqData, privKey)
      if err != nil {
         return nil, fmt.Errorf("failed to sign message: %w", err)
      }

      return challenge, nil

   } else if drmType == "PlayReady" {
      // 1. Parse PRO
      // Assuming ParsePro takes the base64-decoded WRMHeader/PRO data
      wrm, err := playReady.ParsePro(initData)
      if err != nil {
         return nil, fmt.Errorf("failed to parse playready PRO: %w", err)
      }

      // 2. Load device chain (bcert)
      bcertPath := filepath.Join(keyDir, "bdevcert.dat")
      chainBytes, err := os.ReadFile(bcertPath)
      if err != nil {
         return nil, fmt.Errorf("failed to read %s: %w", bcertPath, err)
      }

      chain, err := playReady.ParseChain(chainBytes)
      if err != nil {
         return nil, fmt.Errorf("failed to parse chain: %w", err)
      }

      // 3. Load device signing key
      privKeyPath := filepath.Join(keyDir, "zprivsig.dat")
      privKeyBytes, err := os.ReadFile(privKeyPath)
      if err != nil {
         return nil, fmt.Errorf("failed to read %s: %w", privKeyPath, err)
      }

      signingKey, err := playReady.ParseRawPrivateKey(privKeyBytes)
      if err != nil {
         return nil, fmt.Errorf("failed to parse private key: %w", err)
      }

      // 4. Extract KID/ContentID from WRM Header
      // NOTE: Depending on your specific xml.WrmHeader struct implementation,
      // you might need to adjust the exact fields here to get the KID bytes.
      // As a fallback to compile, we initialize an empty byte slice and string.
      var kid []byte
      var contentID string

      // If wrm exposes KID (uncomment/adjust based on your package):
      // kid = wrm.KID
      // contentID = string(kid)

      _ = wrm // To avoid unused variable error if fields are commented out

      // 5. Generate License Request Bytes (SOAP XML)
      challenge, err := chain.LicenseRequestBytes(signingKey, kid, contentID)
      if err != nil {
         return nil, fmt.Errorf("failed to generate PR license request: %w", err)
      }

      return challenge, nil
   }

   return nil, fmt.Errorf("unsupported DRM type: %s", drmType)
}
