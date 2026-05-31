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
   "testing"
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

// TEST 1: Perform the login flow using dynamically fetched credentials and save the cookies to a temp file
func TestLoginAndSaveSession(t *testing.T) {
   fmt.Println("Fetching credentials from credential.exe...")
   email, password, err := fetchCredentials()
   if err != nil {
      t.Fatalf("Credential fetch failed: %v", err)
   }
   fmt.Printf("Using email: %s\n", email)

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
      t.Fatalf("GetSigninPage failed: %v", err)
   }

   fmt.Println("Step 2: Submitting email claim...")
   nextPageData, err := PostClaim(client, pageData, email)
   if err != nil {
      t.Fatalf("PostClaim failed: %v", err)
   }

   fmt.Println("Step 3: Submitting password...")
   err = PostSignin(client, nextPageData, email, password)
   if err != nil {
      t.Fatalf("PostSignin failed: %v", err)
   }

   // Extract cookies for the amazon domain
   amazonURL, _ := url.Parse("https://www.amazon.com")
   cookies := client.Jar.Cookies(amazonURL)

   // Save the session cookies to the temp file
   cookieFile := getCookieFilePath()
   err = saveCookies(cookieFile, cookies)
   if err != nil {
      t.Fatalf("Failed to save cookies: %v", err)
   }

   fmt.Printf("Successfully saved %d cookies to %s\n", len(cookies), cookieFile)
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

   fmt.Println("Successfully retrieved playback resources!")
   fmt.Println(string(playbackData))
}
