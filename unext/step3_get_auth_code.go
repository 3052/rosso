// step3_get_auth_code.go
package unext

import (
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func Step3GetAuthCode(postAuthEndpoint, codeChallenge string) (string, error) {
   fullURL := "https://oauth.unext.jp" + postAuthEndpoint

   form := url.Values{}
   form.Set("code_challenge", codeChallenge)
   form.Set("code_challenge_method", "S256")

   req, err := http.NewRequest("POST", fullURL, strings.NewReader(form.Encode()))
   if err != nil {
      return "", fmt.Errorf("step3: creating request: %w", err)
   }

   req.Header.Set("user-agent", "U-NEXT Phone App Android12 5.71.0 sdk_gphone64_x86_64")
   req.Header.Set("content-type", "application/x-www-form-urlencoded")

   resp, err := clientDo(req)
   if err != nil {
      return "", fmt.Errorf("step3: sending request: %w", err)
   }
   defer resp.Body.Close()
   io.Copy(io.Discard, resp.Body)

   if resp.StatusCode != http.StatusFound {
      return "", fmt.Errorf("step3: expected 302, got %d", resp.StatusCode)
   }

   location := resp.Header.Get("Location")
   if location == "" {
      return "", fmt.Errorf("step3: no Location header in response")
   }

   locURL, err := url.Parse(location)
   if err != nil {
      return "", fmt.Errorf("step3: parsing Location: %w", err)
   }

   code := locURL.Query().Get("code")
   if code == "" {
      return "", fmt.Errorf("step3: code not found in Location: %s", location)
   }

   return code, nil
}
