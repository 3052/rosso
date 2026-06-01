// post_password.go
package amazon

import (
   "fmt"
   "io"
   "net/http"
   "net/url"
   "os"
   "path/filepath"
   "strings"
)

func PostPassword(s *Session, action string, inputs map[string]string) error {
   data := url.Values{}
   data.Set("password", s.Password)
   keys := []string{
      "anti-csrftoken-a2z",
      "appAction",
      "appActionToken",
      "email",
      "metadata1",
      "openid.return_to",
      "workflowState",
   }
   for _, key := range keys {
      data.Set(key, inputs[key])
   }
   req, err := http.NewRequest("POST", action, strings.NewReader(data.Encode()))
   if err != nil {
      return err
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   resp, err := s.Client.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()

   body, _ := io.ReadAll(resp.Body)
   bodyStr := string(body)

   // Check for raw HTTP errors (like 404 or 500)
   if resp.StatusCode >= 400 {
      errFile := filepath.Join(os.TempDir(), "error_post_password.html")
      os.WriteFile(errFile, body, 0644)
      if strings.Contains(bodyStr, "Looking for Something?") {
         return fmt.Errorf("password submission failed: 404 Not Found. Bad action URL: %s", action)
      }
      return fmt.Errorf("password submission failed with HTTP %d. See %s", resp.StatusCode, errFile)
   }

   // Check if we are still on the sign-in/verification page
   finalPath := resp.Request.URL.Path
   if strings.Contains(finalPath, "/ap/signin") || strings.Contains(finalPath, "/ap/cvf") {
      errDetails := CheckAmazonErrors(body)
      if errDetails == nil {
         errDetails = fmt.Errorf("still on sign-in or verification page")
      }
      errFile := filepath.Join(os.TempDir(), "error_post_password.html")
      os.WriteFile(errFile, body, 0644)
      return fmt.Errorf("password submission failed: %v. Response saved to %s", errDetails, errFile)
   }

   return nil
}
