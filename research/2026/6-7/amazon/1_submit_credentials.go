package amazon

import (
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "os"
   "strings"
)

// SubmitCredentials posts the sign-in form with the user's credentials.
func SubmitCredentials(client *http.Client, sessionId string, formValues url.Values, referer string) (string, error) {
   postUrl := fmt.Sprintf("https://www.amazon.com/ap/signin/%s", sessionId)

   req, err := http.NewRequest("POST", postUrl, strings.NewReader(formValues.Encode()))
   if err != nil {
      return "", err
   }

   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 11; sdk_gphone_x86_64 Build/RSR1.240422.006; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/83.0.4103.106 Mobile Safari/537.36")
   req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
   req.Header.Set("X-Requested-With", "com.amazon.avod.thirdpartyclient")
   req.Header.Set("Origin", "https://www.amazon.com")
   req.Header.Set("Referer", referer)
   req.Header.Set("Sec-Fetch-Site", "same-origin")
   req.Header.Set("Sec-Fetch-Mode", "navigate")
   req.Header.Set("Sec-Fetch-User", "?1")
   req.Header.Set("Sec-Fetch-Dest", "document")
   req.Header.Set("Accept-Language", "en-US,en;q=0.9")

   // Set client to stop at redirect so we can capture the Location
   client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
      return http.ErrUseLastResponse
   }

   resp, err := client.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   bodyBytes, _ := io.ReadAll(resp.Body)
   html := string(bodyBytes)

   if resp.StatusCode == http.StatusOK {
      os.WriteFile("captcha_debug.html", bodyBytes, 0644)
      if strings.Contains(html, "cvf-aamation-challenge-form") || strings.Contains(html, "Authentication required") {
         return "", fmt.Errorf("CAPTCHA_REQUIRED")
      }
      return "", fmt.Errorf("expected 302 redirect, got 200 OK. Debug HTML saved to captcha_debug.html")
   }

   if resp.StatusCode != http.StatusFound {
      return "", fmt.Errorf("expected 302 redirect, got status code: %d", resp.StatusCode)
   }

   location := resp.Header.Get("Location")
   if location == "" {
      return "", errors.New("location header not found in the 302 response")
   }

   if strings.HasPrefix(location, "/") {
      location = "https://www.amazon.com" + location
   }

   return location, nil
}
