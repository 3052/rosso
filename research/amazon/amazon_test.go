// amazon_test.go
package amazon

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/http/cookiejar"
   "net/url"
   "os"
   "os/exec"
   "path/filepath"
   "strings"
   "testing"
)

type SavedState struct {
   Cookies  []*http.Cookie `json:"cookies"`
   VideoID  string         `json:"video_id"`
   DeviceID string         `json:"device_id"`
}

type Credential struct {
   Date     string `json:"date"`
   Host     string `json:"host"`
   Password string `json:"password"`
   Trial    string `json:"trial"`
   Username string `json:"username"`
}

type PlaybackResponse struct {
   VodPlaylistedPlaybackUrls struct {
      Result struct {
         PlaybackUrls struct {
            IntraTitlePlaylist []struct {
               Urls []struct {
                  URL string `json:"url"`
               } `json:"urls"`
            } `json:"intraTitlePlaylist"`
         } `json:"playbackUrls"`
      } `json:"result"`
   } `json:"vodPlaylistedPlaybackUrls"`
}

func getTempFilePath() string {
   return filepath.Join(os.TempDir(), "amazon_test_session.json")
}

func getCredentials() (string, string, error) {
   cmd := exec.Command("credential.exe", "-j=amazon.com")
   output, err := cmd.Output()
   if err != nil {
      return "", "", fmt.Errorf("failed to execute credential.exe: %w", err)
   }

   var creds []Credential
   if err := json.Unmarshal(output, &creds); err != nil {
      return "", "", fmt.Errorf("failed to parse credentials JSON: %w\nOutput: %s", err, string(output))
   }

   if len(creds) == 0 {
      return "", "", fmt.Errorf("no credentials found in output")
   }

   return creds[0].Username, creds[0].Password, nil
}

func TestLoginAndSave(t *testing.T) {
   email, password, err := getCredentials()
   if err != nil {
      t.Fatalf("Could not get credentials: %v", err)
   }

   videoID := "B075RND57T"

   jar, err := cookiejar.New(nil)
   if err != nil {
      t.Fatalf("Failed to create cookie jar: %v", err)
   }

   session := &Session{
      Client:   &http.Client{Jar: jar},
      Email:    email,
      Password: password,
      VideoID:  videoID,
      DeviceID: GenerateUUID(),
   }

   // Step 1: Get Sign-in page
   action, inputs, err := GetSignIn(session)
   if err != nil {
      t.Fatalf("GetSignIn failed: %v", err)
   }

   // Step 2: Post Email
   nextAction, nextInputs, err := PostEmail(session, action, inputs)
   if err != nil {
      t.Fatalf("PostEmail failed: %v", err)
   }

   // Step 3: Post Password
   err = PostPassword(session, nextAction, nextInputs)
   if err != nil {
      t.Fatalf("PostPassword failed: %v", err)
   }

   // Extract cookies to save
   amazonURL, _ := url.Parse("https://www.amazon.com")
   state := SavedState{
      Cookies:  session.Client.Jar.Cookies(amazonURL),
      VideoID:  session.VideoID,
      DeviceID: session.DeviceID,
   }

   data, err := json.Marshal(state)
   if err != nil {
      t.Fatalf("Failed to marshal state: %v", err)
   }

   filePath := getTempFilePath()
   err = os.WriteFile(filePath, data, 0600)
   if err != nil {
      t.Fatalf("Failed to write state file: %v", err)
   }

   t.Logf("Successfully logged in as %s and saved state to %s", email, filePath)
}

func TestLoadAndFetchRemaining(t *testing.T) {
   filePath := getTempFilePath()
   data, err := os.ReadFile(filePath)
   if err != nil {
      t.Fatalf("Failed to read state file: %v. Please run TestLoginAndSave first.", err)
   }

   var state SavedState
   if err := json.Unmarshal(data, &state); err != nil {
      t.Fatalf("Failed to unmarshal state: %v", err)
   }

   // Reconstruct CookieJar
   jar, err := cookiejar.New(nil)
   if err != nil {
      t.Fatalf("Failed to create cookie jar: %v", err)
   }

   amazonURL, _ := url.Parse("https://www.amazon.com")
   atvURL, _ := url.Parse("https://atv-ps.amazon.com")

   // Apply saved cookies to relevant Amazon domains
   jar.SetCookies(amazonURL, state.Cookies)
   jar.SetCookies(atvURL, state.Cookies)

   session := &Session{
      Client:   &http.Client{Jar: jar},
      VideoID:  state.VideoID,
      DeviceID: state.DeviceID,
   }

   // Step 4: Get Metadata (Sets Envelope and TargetTitleID)
   err = GetEnrichItemMetadata(session)
   if err != nil {
      t.Fatalf("GetEnrichItemMetadata failed: %v", err)
   }

   // Step 5: Get VOD Playback Resources
   response, err := GetVodPlaybackResources(session)
   if err != nil {
      t.Fatalf("GetVodPlaybackResources failed: %v", err)
   }

   if response == "" {
      t.Fatalf("Received empty response from GetVodPlaybackResources")
   }

   t.Logf("Successfully fetched playback resources! Payload length: %d bytes", len(response))

   // Parse JSON and output MPD URLs
   var pr PlaybackResponse
   if err := json.Unmarshal([]byte(response), &pr); err != nil {
      t.Fatalf("Failed to unmarshal playback response: %v", err)
   }

   found := false
   for _, playlist := range pr.VodPlaylistedPlaybackUrls.Result.PlaybackUrls.IntraTitlePlaylist {
      for _, u := range playlist.Urls {
         if strings.Contains(u.URL, ".mpd") {
            t.Logf("MPD URL Found: %s", u.URL)
            found = true
         }
      }
   }

   if !found {
      t.Log("No MPD URLs found in the response.")
   }
}
