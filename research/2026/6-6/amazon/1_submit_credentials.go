package amazon

import (
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

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

   resp, err := client.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   _, _ = io.Copy(io.Discard, resp.Body)

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
