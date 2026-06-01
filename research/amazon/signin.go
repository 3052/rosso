// get_signin.go
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

func GetSignIn(s *Session) (string, map[string]string, error) {
   req0, _ := http.NewRequest("GET", "https://www.amazon.com/", nil)
   req0.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   req0.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
   resp0, err := s.Client.Do(req0)
   if err == nil {
      io.Copy(io.Discard, resp0.Body)
      resp0.Body.Close()
   }

   returnUrl := url.QueryEscape("https://www.amazon.com/gp/video/detail/" + s.VideoID + "?ref_=nav_custrec_signin")
   targetURL := "https://www.amazon.com/ap/signin?openid.return_to=" + returnUrl + "&openid.identity=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0%2Fidentifier_select&openid.assoc_handle=usflex&openid.mode=checkid_setup&openid.claimed_id=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0%2Fidentifier_select&openid.ns=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0"

   req, err := http.NewRequest("GET", targetURL, nil)
   if err != nil {
      return "", nil, err
   }

   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
   req.Header.Set("Accept-Language", "en-US,en;q=0.5")
   req.Header.Set("Upgrade-Insecure-Requests", "1")
   req.Header.Set("Sec-Fetch-Dest", "document")
   req.Header.Set("Sec-Fetch-Mode", "navigate")
   req.Header.Set("Sec-Fetch-Site", "same-origin")

   resp, err := s.Client.Do(req)
   if err != nil {
      return "", nil, err
   }
   defer resp.Body.Close()

   body, _ := io.ReadAll(resp.Body)
   action, inputs := ExtractForm(string(body), "signIn")

   if action == "" {
      errFile := filepath.Join(os.TempDir(), "error_get_signin.html")
      os.WriteFile(errFile, body, 0644)
      return "", nil, fmt.Errorf("signIn form not found. Response body saved to %s", errFile)
   }

   if strings.HasPrefix(action, "/") {
      action = "https://www.amazon.com" + action
   }

   return action, inputs, nil
}
