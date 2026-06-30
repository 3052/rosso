package main

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "os"
   "strconv"
   "strings"
)

func extractInitialState(html string) (string, error) {
   _, after, found := strings.Cut(html, `window.__INITIAL_STATE__ = "`)
   if !found {
      return "", fmt.Errorf("could not find 'window.__INITIAL_STATE__ =' in the HTML")
   }

   rawEscapedJSON, _, found := strings.Cut(after, `";</script>`)
   if !found {
      return "", fmt.Errorf("could not find the closing '\";</script>' marker")
   }

   return rawEscapedJSON, nil
}

func fetchURL(url string) (string, error) {
   resp, err := http.Get(url)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return "", fmt.Errorf("unexpected HTTP status: %d", resp.StatusCode)
   }

   bodyBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      return "", err
   }
   return string(bodyBytes), nil
}

func getAssetID(deepLinkURL string) (string, error) {
   // Prepare JSON payload
   payloadData := map[string]string{
      "deeplink": deepLinkURL,
   }
   bodyData, _ := json.Marshal(payloadData)

   req, err := http.NewRequest("POST", "https://api-eu.fubo.tv/papi/v1/deeplink", bytes.NewBuffer(bodyData))
   if err != nil {
      return "", err
   }

   req.Header.Add("authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2FwaS5mdWJvLnR2Iiwic3ViIjoiNmExOGMzYmZkMmQ1MmYwMDAxMTAyNTExIiwiYXVkIjpbImtRc3JMSldSWWszSEk0RmpMRWVzWjZpUkhZUkRlYWZyIl0sImV4cCI6MTc4Mjg1NzQ2MiwiaWF0IjoxNzgyODIxNDYyLCJ0eXBlIjoiYWNjZXNzIiwiaHR0cHM6Ly9yYmFjL2FwcCI6ImZ1Ym90diIsImh0dHBzOi8vcmJhYy9yb2xlcyI6WyJlbmR1c2VyIl0sImRldmljZV9pZCI6IngtZGV2aWNlLWlkIn0.aH9vSV4aniwmaXvQFxodiWf4RBk9k_NC0V_G68yjFhA")
   req.Header.Add("x-application-id", "molotov")
   req.Header.Add("x-device-app", "web")

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return "", fmt.Errorf("post request failed: %v", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return "", fmt.Errorf("deeplink API returned status: %d", resp.StatusCode)
   }

   bodyBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      return "", fmt.Errorf("failed to read response: %v", err)
   }

   var dlResp DeepLinkResponse
   if err := json.Unmarshal(bodyBytes, &dlResp); err != nil {
      return "", fmt.Errorf("failed to parse DeepLink JSON response: %v", err)
   }

   if len(dlResp.Actions) == 0 || dlResp.Actions[0].PublicPath == "" {
      return "", fmt.Errorf("no valid actions/public_path found in response body")
   }

   // Extract the asset ID from the PublicPath (e.g. "/program-details/program/VOD_314017")
   // By splitting at "/", the last element is our "VOD_314017"
   parts := strings.Split(dlResp.Actions[0].PublicPath, "/")
   assetID := parts[len(parts)-1]

   return assetID, nil
}

func main() {
   if len(os.Args) < 2 {
      fmt.Println("Usage: go run main.go <molotov_url>")
      fmt.Println("Example: go run main.go https://www.molotov.tv/fr_fr/p/194433/le-grand-jeu")
      os.Exit(1)
   }
   publicURL := os.Args[1]

   // --- STEP 1: Fetch the HTML page ---
   fmt.Println("[1/4] Fetching Molotov public page...")
   htmlContent, err := fetchURL(publicURL)
   if err != nil {
      fmt.Printf("Error fetching URL: %v\n", err)
      os.Exit(1)
   }

   // --- STEP 2: Extract & unescape the JSON Initial State ---
   fmt.Println("[2/4] Extracting JSON initial state from HTML...")
   rawEscapedJSON, err := extractInitialState(htmlContent)
   if err != nil {
      fmt.Printf("Extraction failed: %v\n", err)
      os.Exit(1)
   }

   unescapedJSON, err := strconv.Unquote(`"` + rawEscapedJSON + `"`)
   if err != nil {
      fmt.Printf("Warning: Could not unescape string completely: %v\n", err)
      unescapedJSON = rawEscapedJSON // Fallback
   }

   // --- STEP 3: Parse JSON & find the deep link URL ---
   fmt.Println("[3/4] Parsing JSON to find deep link URL...")
   var state MolotovState
   if err := json.Unmarshal([]byte(unescapedJSON), &state); err != nil {
      fmt.Printf("Error parsing JSON: %v\n", err)
      os.Exit(1)
   }

   deepLink := ""
   // Iterate through programs map dynamically (bypasses needing to know "194433" upfront)
   for _, program := range state.Programs.Programs {
      if len(program.Channels) > 0 && len(program.Channels[0].ProgramHeader.Buttons) > 0 {
         deepLink = program.Channels[0].ProgramHeader.Buttons[0].URL
         break // Taking the first available URL as requested
      }
   }

   if deepLink == "" {
      fmt.Println("Error: Could not find the deep link URL in the extracted JSON structure.")
      os.Exit(1)
   }
   fmt.Printf("      -> Found deep link: %s\n", deepLink)

   // --- STEP 4: Request Fubo Deeplink API & extract Asset ID ---
   fmt.Println("[4/4] Sending POST request to Deeplink API...")
   assetID, err := getAssetID(deepLink)
   if err != nil {
      fmt.Printf("Error getting Asset ID: %v\n", err)
      os.Exit(1)
   }

   fmt.Println("\n==============================================")
   fmt.Printf("SUCCESS! Final Asset ID: %s\n", assetID)
   fmt.Println("==============================================")
}

// Minimal struct to extract the public_path from the DeepLink API response
type DeepLinkResponse struct {
   Actions []struct {
      PublicPath string `json:"public_path"`
   } `json:"actions"`
}

// Minimal struct to extract just the specific URL from the Molotov initial state
type MolotovState struct {
   Programs struct {
      Programs map[string]struct {
         Channels []struct {
            ProgramHeader struct {
               Buttons []struct {
                  URL string `json:"url"`
               } `json:"buttons"`
            } `json:"program_header"`
         } `json:"channels"`
      } `json:"programs"`
   } `json:"programs"`
}
