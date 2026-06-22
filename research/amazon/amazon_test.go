package amazon

import (
   "fmt"
   "os"
   "path/filepath"
   "testing"

   "41.neocities.org/diana/playReady"
   "41.neocities.org/diana/widevine"
)

func TestPlayReadySL3000(t *testing.T) {
   runDeviceCombinations(
      t,
      "PlayReady SL3000",
      `C:\Users\Steven\AppData\Local\SL3000`,
      "PlayReady",
   )
}

func TestWidevineL3(t *testing.T) {
   runDeviceCombinations(
      t,
      "Widevine L3",
      `C:\Users\Steven\AppData\Local\L3`,
      "Widevine",
   )
}

func TestPlayReadySL2000(t *testing.T) {
   runDeviceCombinations(
      t,
      "PlayReady SL2000",
      `C:\Users\Steven\AppData\Local\SL2000`,
      "PlayReady",
   )
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
