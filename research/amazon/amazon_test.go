package amazon

import (
   "encoding/json"
   "fmt"
   "net/http"
   "os"
   "path/filepath"
   "strings"
   "testing"
)

type AuthState struct {
   CodePair *CodePairResponse `json:"code_pair"`
   Device   map[string]string `json:"device"`
}

const (
   codePairEndpoint = "https://api.amazon.com/auth/create/codepair"
   registerEndpoint = "https://api.amazon.com/auth/register"
   playbackEndpoint = "https://atv-ps.amazon.com/cdp/catalog/GetPlaybackResources"
   marketplaceIDUS  = "ATVPDKIKX0DER"
)

// Define the device identity we are pretending to be
var defaultDevice = map[string]string{
   "domain":        "Device",
   "app_name":      "com.amazon.amazonvideo.livingroom",
   "app_version":   "1.1",
   "device_model":  "LG-Tv",
   "os_version":    "6.0.1",
   "device_type":   "A71I8788P1ZV8",
   "device_name":   "My Go Device",
   "device_serial": "a906a7f9bfd6a7ab",
}

// Helper functions for temp files
func getTempStatePath() string {
   return filepath.Join(os.TempDir(), "amazon_auth_state.json")
}

func getTempTokenPath() string {
   return filepath.Join(os.TempDir(), "amazon_token.txt")
}

// 1. Run this to get the code
func TestStep1_StartProcess(t *testing.T) {
   client := &http.Client{}

   t.Log("Fetching Code Pair...")
   codePair, err := GetCodePair(client, codePairEndpoint, defaultDevice)
   if err != nil {
      t.Fatalf("Failed to get Code Pair: %v", err)
   }

   state := AuthState{
      CodePair: codePair,
      Device:   defaultDevice,
   }

   data, _ := json.MarshalIndent(state, "", "  ")
   stateFile := getTempStatePath()
   if err := os.WriteFile(stateFile, data, 0644); err != nil {
      t.Fatalf("Failed to write state file to temp dir: %v", err)
   }

   fmt.Println()
   fmt.Println("========================================================================")
   fmt.Println("1. Go to: https://amazon.com/mytv")
   fmt.Println("2. Sign in to your Amazon account normally using your email & password.")
   fmt.Println("   (⚠️  IMPORTANT: Do NOT click the 'Sign in with a code' button at the")
   fmt.Println("   bottom of the login screen. That is a different Amazon feature.)")
   fmt.Println("3. Once signed in, you will be on the 'Register Your Device' screen.")
   fmt.Printf("4. Enter this code: %s\n", codePair.PublicCode)
   fmt.Println("5. Click 'Register Device'.")
   fmt.Println("6. Finally, come back here and run TestStep2_CompleteProcess")
   fmt.Println("========================================================================")
   fmt.Println()
}

// 2. Run this after approving the code in your browser
func TestStep2_CompleteProcess(t *testing.T) {
   client := &http.Client{}

   stateFile := getTempStatePath()
   data, err := os.ReadFile(stateFile)
   if err != nil {
      t.Fatalf("Failed to read state file from temp dir. Run TestStep1 first. Error: %v", err)
   }

   var state AuthState
   if err := json.Unmarshal(data, &state); err != nil {
      t.Fatalf("Failed to parse state file: %v", err)
   }

   t.Log("Registering device (checking if you entered the code)...")
   regResponse, err := RegisterDevice(client, registerEndpoint, state.CodePair, state.Device)
   if err != nil {
      t.Fatalf("Failed to register device: %v\n(Did you forget to enter the code at amazon.com/mytv?)", err)
   }

   bearer := regResponse.Response.Success.Tokens.Bearer

   // Write token to temp dir for Step 3 to use
   tokenFile := getTempTokenPath()
   if err := os.WriteFile(tokenFile, []byte(bearer.AccessToken), 0644); err != nil {
      t.Fatalf("Failed to save access token to temp dir: %v", err)
   }

   fmt.Println()
   fmt.Println("=====================================================")
   fmt.Println("SUCCESS! Final Credentials Generated:")
   fmt.Printf("Access Token: %s\n", bearer.AccessToken)
   if bearer.RefreshToken != "" {
      fmt.Printf("Refresh Token: %s\n", bearer.RefreshToken)
   }
   fmt.Printf("Expires In: %s seconds\n", bearer.ExpiresIn)
   fmt.Printf("Token saved to: %s\n", tokenFile)
   fmt.Println("=====================================================")
   fmt.Println()

   // Clean up the intermediate state file
   _ = os.Remove(stateFile)
}

// 3. Run this to grab an MPD after you have generated the token in Step 2.
func TestStep3_GetMPD(t *testing.T) {
   client := &http.Client{}

   tokenFile := getTempTokenPath()
   tokenBytes, err := os.ReadFile(tokenFile)
   if err != nil || len(tokenBytes) == 0 {
      t.Fatalf("Failed to read token from %s. Please run TestStep1 and TestStep2 first.", tokenFile)
   }
   accessToken := strings.TrimSpace(string(tokenBytes))

   // For testing, pick a free/owned ASIN or Prime included ASIN (if you have Prime).
   asin := "B002Y27P3M" // Example ASIN (Iron Man / usually free with ads / Prime)

   opts := DefaultPlaybackOptions()
   opts.VideoQuality = "HD"
   opts.VideoCodec = "H264"
   opts.BitrateMode = "CVBR,CBR" // Fallback map to python logic

   t.Logf("Fetching manifest for ASIN %s...", asin)
   manifestResp, err := GetPlaybackResources(
      client,
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

   t.Logf("Raw MPD URL: %s", rawMPD)

   cleanMPD := CleanMPDURL(rawMPD)
   t.Logf("Cleaned MPD URL: %s", cleanMPD)

   fmt.Println()
   fmt.Println("=====================================================")
   fmt.Println("SUCCESS! Final DASH Manifest (MPD):")
   fmt.Printf("%s\n", cleanMPD)
   fmt.Println("=====================================================")
}
