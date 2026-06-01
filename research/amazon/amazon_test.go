// amazon_test.go
package amazon

import (
   "encoding/json"
   "fmt"
   "io/ioutil"
   "net/http"
   "net/http/cookiejar"
   "net/url"
   "os"
   "os/exec"
   "path/filepath"
   "strings"
   "testing"
   "time"
)

// Helper function to get the temporary file path for storing cookies
func getCookieFilePath() string {
   return filepath.Join(os.TempDir(), "amazon_cookies.json")
}

// Credential represents the JSON structure returned by credential.exe
type Credential struct {
   Username string `json:"username"`
   Password string `json:"password"`
}

// Helper function to fetch credentials dynamically from credential.exe
func fetchCredentials() (string, string, error) {
   cmd := exec.Command("credential.exe", "-j=amazon.com")
   output, err := cmd.Output()
   if err != nil {
      return "", "", fmt.Errorf("failed to execute credential.exe: %w\nOutput: %s", err, string(output))
   }

   var creds []Credential
   if err := json.Unmarshal(output, &creds); err != nil {
      return "", "", fmt.Errorf("failed to parse credentials JSON: %w", err)
   }

   if len(creds) == 0 {
      return "", "", fmt.Errorf("no credentials returned from credential.exe")
   }

   return creds[0].Username, creds[0].Password, nil
}

// Helper function to save cookies to a JSON file
func saveCookies(filename string, cookies []*http.Cookie) error {
   data, err := json.MarshalIndent(cookies, "", "  ")
   if err != nil {
      return err
   }
   return ioutil.WriteFile(filename, data, 0644)
}

// Helper function to load cookies from a JSON file
func loadCookies(filename string) ([]*http.Cookie, error) {
   data, err := ioutil.ReadFile(filename)
   if err != nil {
      return nil, err
   }
   var cookies []*http.Cookie
   err = json.Unmarshal(data, &cookies)
   return cookies, err
}

// TEST 1: Perform the login flow using dynamically fetched credentials and save the cookies to a temp file.
// Runs twice to ensure stability against Amazon's anti-bot mechanics.
func TestLoginAndSaveSession(t *testing.T) {
   // --- Tweak these variables if Amazon acts up ---
   const numRuns = 2
   const sleepBetweenSteps = 4 * time.Second // Simulates human typing speed
   const sleepBetweenRuns = 15 * time.Second // Longer delay between full login attempts to avoid IP flags
   // -----------------------------------------------

   fmt.Println("Fetching credentials from credential.exe...")
   email, password, err := fetchCredentials()
   if err != nil {
      t.Fatalf("Credential fetch failed: %v", err)
   }
   fmt.Printf("Using email: %s\n", email)

   for i := 1; i <= numRuns; i++ {
      fmt.Printf("\n--- Starting Login Run %d ---\n", i)

      // Create a fresh cookie jar and client for each run to avoid state contamination
      jar, err := cookiejar.New(nil)
      if err != nil {
         t.Fatal(err)
      }

      client := &http.Client{
         Jar: jar,
      }

      fmt.Println("Step 1: Fetching initial sign-in page...")
      pageData, err := GetSigninPage(client)
      if err != nil {
         t.Fatalf("Run %d: GetSigninPage failed: %v", i, err)
      }

      fmt.Printf("Sleeping for %v (Simulating typing email)...\n", sleepBetweenSteps)
      time.Sleep(sleepBetweenSteps)

      fmt.Println("Step 2: Submitting email claim...")
      nextPageData, err := PostClaim(client, pageData, email)
      if err != nil {
         t.Fatalf("Run %d: PostClaim failed: %v", i, err)
      }

      fmt.Printf("Sleeping for %v (Simulating typing password)...\n", sleepBetweenSteps)
      time.Sleep(sleepBetweenSteps)

      fmt.Println("Step 3: Submitting password...")
      err = PostSignin(client, nextPageData, email, password)
      if err != nil {
         t.Fatalf("Run %d: PostSignin failed: %v", i, err)
      }

      // Extract cookies for the amazon domain
      amazonURL, _ := url.Parse("https://www.amazon.com")
      cookies := client.Jar.Cookies(amazonURL)

      // Save the session cookies to the temp file
      cookieFile := getCookieFilePath()
      err = saveCookies(cookieFile, cookies)
      if err != nil {
         t.Fatalf("Run %d: Failed to save cookies: %v", i, err)
      }

      fmt.Printf("Run %d: Successfully saved %d cookies to %s\n", i, len(cookies), cookieFile)

      // Sleep before the next run (if not the last run)
      if i < numRuns {
         fmt.Printf("Sleeping for %v to cool down Amazon's rate-limiting before the next run...\n", sleepBetweenRuns)
         time.Sleep(sleepBetweenRuns)
      }
   }
}

// TEST 2: Read the cookies from the temp file and request playback resources
func TestGetPlaybackResources(t *testing.T) {
   cookieFile := getCookieFilePath()

   // Check if the cookie file exists first
   if _, err := os.Stat(cookieFile); os.IsNotExist(err) {
      t.Fatalf("Cookie file %s does not exist. Please run TestLoginAndSaveSession first.", cookieFile)
   }

   // Load cookies from file
   cookies, err := loadCookies(cookieFile)
   if err != nil {
      t.Fatalf("Failed to load cookies: %v", err)
   }

   jar, err := cookiejar.New(nil)
   if err != nil {
      t.Fatal(err)
   }

   // Apply the loaded cookies to the relevant Amazon domains
   amazonURL, _ := url.Parse("https://www.amazon.com")
   apiURL, _ := url.Parse("https://atv-ps.amazon.com")
   jar.SetCookies(amazonURL, cookies)
   jar.SetCookies(apiURL, cookies)

   client := &http.Client{
      Jar: jar,
   }

   fmt.Println("Step 4: Requesting playback resources with loaded session...")
   playbackData, err := GetPlaybackResources(client)
   if err != nil {
      t.Fatalf("GetPlaybackResources failed: %v", err)
   }

   // Convert response to string to search for the MPD extension
   responseStr := string(playbackData)

   // Validate the response contains the expected DASH manifest (.mpd)
   if !strings.Contains(responseStr, ".mpd") {
      t.Fatalf("Test Failed: .mpd manifest URL not found in the response. Response snippet: %s", responseStr[:500])
   }

   fmt.Println("Success! Found .mpd manifest URL in the response.")

   // Optional: Print a snippet of the response to keep the console clean
   if len(responseStr) > 500 {
      fmt.Printf("Response snippet: %s...\n", responseStr[:500])
   } else {
      fmt.Printf("Response: %s\n", responseStr)
   }
}
