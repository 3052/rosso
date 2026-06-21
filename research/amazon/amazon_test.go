package amazon

import (
   "encoding/base64"
   "encoding/xml"
   "fmt"
   "io"
   "net/http"
   "testing"
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
      t.Skip("Skipping Amazon API test: missing credentials")
   }

   client := NewClient(&http.Client{})

   // Define the 3 devices you want to test
   tests := []struct {
      Name    string
      Profile DeviceProfile
   }{
      {
         Name: "Widevine L3",
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
         Name: "PlayReady SL2000",
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
         Name: "PlayReady SL3000",
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
         challengeBytes, err := generateCDMChallenge(tc.Profile.DRMType, initDataBytes)
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

// generateCDMChallenge is the hook for your local CDM devices.
func generateCDMChallenge(drmType string, initData []byte) ([]byte, error) {
   // TODO: Initialize your local Widevine/PlayReady CDM instance here.
   // Feed 'initData' into the CDM session to produce the license challenge payload.
   //
   // For Widevine: return the raw challenge bytes.
   // For PlayReady: return the generated SOAP XML string cast to []byte.

   return []byte("mock_challenge"), nil
}
