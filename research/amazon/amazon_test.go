// amazon_test.go
package amazon

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

var (
	apiBaseURL   = "api.amazon.com"
	manifestBase = "atv-ps.amazon.com"
	hostBase     = "https://www.amazon.com/gp/video"
	titleID      = "B075RND57T" // Updated TitleID
	marketplace  = "ATVPDKIKX0DER"

	// Mock Android TV device data
	deviceData = map[string]interface{}{
		"manufacturer":   "Hisense",
		"device_chipset": "mt7663",
		"domain":         "Device",
		"app_name":       "AIV",
		"os_name":        "Android",
		"app_version":    "3.12.0",
		"device_model":   "HAT4KDTV",
		"os_version":     "VIDAA",
		"device_type":    "A2RGJ95OVLR12U",
		"device_serial":  "13f5b56b4a17de5d136f0e4c28236109",
		"device_name":    "Test Hisense TV",
	}
)

func getTempFile(name string) string {
	return filepath.Join(os.TempDir(), name)
}

// 1. function for codepair that writes the result to os.TempDir
// Run with: go test -v -run TestCodePair
func TestCodePair(t *testing.T) {
	client := &http.Client{}
	ctx := context.Background()

	result, err := CreateCodePair(ctx, client, apiBaseURL, deviceData)
	if err != nil {
		t.Fatalf("CreateCodePair failed: %v", err)
	}

	fileData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal codepair result: %v", err)
	}

	outFile := getTempFile("amazon_codepair.json")
	if err := os.WriteFile(outFile, fileData, 0644); err != nil {
		t.Fatalf("Failed to write codepair to temp dir: %v", err)
	}

	publicCode, ok := result["public_code"].(string)
	if !ok {
		t.Fatalf("Failed to parse public_code from response: %v", result)
	}

	t.Logf("\n")
	t.Logf("=================================================================")
	t.Logf("ACTION REQUIRED BEFORE RUNNING THE NEXT TEST:")
	t.Logf("1. Open your browser and go to: https://www.amazon.com/mytv")
	t.Logf("2. Log in to your Amazon account")
	t.Logf("3. Enter this code: %s", publicCode)
	t.Logf("4. Once successfully registered in the browser, run:")
	t.Logf("   go test -v -run TestRegister")
	t.Logf("=================================================================\n")
	t.Logf("CodePair data written to: %s", outFile)
}

// 2. function for register that reads input from os.TempDir and writes output to os.TempDir
// Run with: go test -v -run TestRegister
func TestRegister(t *testing.T) {
	client := &http.Client{}
	ctx := context.Background()

	inFile := getTempFile("amazon_codepair.json")
	fileData, err := os.ReadFile(inFile)
	if err != nil {
		t.Fatalf("Failed to read codepair from temp dir: %v", err)
	}

	var codePair map[string]interface{}
	if err := json.Unmarshal(fileData, &codePair); err != nil {
		t.Fatalf("Failed to unmarshal codepair data: %v", err)
	}

	result, err := RegisterDevice(ctx, client, apiBaseURL, codePair, deviceData)
	if err != nil {
		t.Fatalf("RegisterDevice failed (Did you enter the code at amazon.com/mytv?): %v", err)
	}

	outData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal register result: %v", err)
	}

	outFile := getTempFile("amazon_register.json")
	if err := os.WriteFile(outFile, outData, 0644); err != nil {
		t.Fatalf("Failed to write register output to temp dir: %v", err)
	}

	// Safely extract and print the credentials to the console
	if respMap, ok := result["response"].(map[string]interface{}); ok {
		if successMap, ok := respMap["success"].(map[string]interface{}); ok {
			if tokensMap, ok := successMap["tokens"].(map[string]interface{}); ok {
				if bearerMap, ok := tokensMap["bearer"].(map[string]interface{}); ok {
					accessToken, _ := bearerMap["access_token"].(string)
					refreshToken, _ := bearerMap["refresh_token"].(string)
					expiresIn, _ := bearerMap["expires_in"].(float64)

					t.Logf("\n")
					t.Logf("=================== CREDENTIALS ===================")
					t.Logf("Access Token:  %s", accessToken)
					t.Logf("Refresh Token: %s", refreshToken)
					t.Logf("Expires In:    %v seconds", expiresIn)
					t.Logf("===================================================\n")
				}
			}
		}
	}

	t.Logf("Success! Device registered. Output written to: %s", outFile)
	t.Logf("You can now run: go test -v -run TestPlayback")
}

// 3. function for playback that reads input from os.TempDir
// Run with: go test -v -run TestPlayback
func TestPlayback(t *testing.T) {
	client := &http.Client{}
	ctx := context.Background()

	inFile := getTempFile("amazon_register.json")
	fileData, err := os.ReadFile(inFile)
	if err != nil {
		t.Fatalf("Failed to read register file from temp dir: %v", err)
	}

	var registerData map[string]interface{}
	if err := json.Unmarshal(fileData, &registerData); err != nil {
		t.Fatalf("Failed to unmarshal register data: %v", err)
	}

	// Extract the bearer token safely
	var deviceToken string
	if resp, ok := registerData["response"].(map[string]interface{}); ok {
		if success, ok := resp["success"].(map[string]interface{}); ok {
			if tokens, ok := success["tokens"].(map[string]interface{}); ok {
				if bearer, ok := tokens["bearer"].(map[string]interface{}); ok {
					if access, ok := bearer["access_token"].(string); ok {
						deviceToken = access
					}
				}
			}
		}
	}

	if deviceToken == "" {
		t.Fatalf("Failed to extract device token from %s", inFile)
	}

	metadataResp, err := EnrichItemMetadata(ctx, client, hostBase, []string{titleID})
	if err != nil {
		t.Fatalf("EnrichItemMetadata failed: %v", err)
	}

	var playbackEnvelope string
	if enrichments, ok := metadataResp["enrichments"].(map[string]interface{}); ok {
		if titleData, exists := enrichments[titleID].(map[string]interface{}); exists {
			if playbackActions, ok := titleData["playbackActions"].([]interface{}); ok && len(playbackActions) > 0 {
				firstAction := playbackActions[0].(map[string]interface{})
				if expMeta, ok := firstAction["playbackExperienceMetadata"].(map[string]interface{}); ok {
					if env, ok := expMeta["playbackEnvelope"].(string); ok {
						playbackEnvelope = env
					}
				}
			}
		}
	}

	if playbackEnvelope == "" {
		t.Fatalf("Skipping playback request: Could not get playbackEnvelope. Title might be invalid or require cookies. Response: %v", metadataResp)
	}

	t.Logf("Successfully extracted PlaybackEnvelope!")

	params := PlaybackParams{
		BaseURL:          manifestBase,
		DeviceID:         deviceData["device_serial"].(string),
		DeviceTypeID:     deviceData["device_type"].(string),
		GascEnabled:      false,
		MarketplaceID:    marketplace,
		TitleID:          titleID,
		DeviceToken:      deviceToken, // Populated from os.TempDir register file
		PlaybackEnvelope: playbackEnvelope,
		Quality:          "UHD",
		VideoCodec:       "H265",
		BitrateMode:      "CVBR",
		HDR:              "SDR",
		IsPlayReady:      false,
		PlayerType:       "html5",
	}

	manifestResp, err := GetVodPlaybackResources(ctx, client, params)
	if err != nil {
		t.Fatalf("GetVodPlaybackResources failed: %v", err)
	}

	if errByRes, ok := manifestResp["errorsByResource"].(map[string]interface{}); ok && len(errByRes) > 0 {
		t.Fatalf("Playback API returned resource errors: %v", errByRes)
	}

	t.Log("Successfully fetched playback manifest.")
}
