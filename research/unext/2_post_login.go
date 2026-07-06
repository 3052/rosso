package unext

import (
   "fmt"
   "net/http"
   "net/url"
   "strings"
)

// PostLogin authenticates the user and retrieves necessary auth cookies
func PostLogin(client *http.Client, csrfToken, recaptchaResponse, loginID, password string) error {
   urlStr := "https://account.unext.jp/login"

   data := url.Values{}
   data.Set("_csrf", csrfToken)
   data.Set("backurl", "https://video.unext.jp/title/SID0020149")
   data.Set("g-recaptcha-response", recaptchaResponse)
   data.Set("login_id", loginID)
   data.Set("password", password)

   req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
   if err != nil {
      return err
   }

   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
   req.Header.Set("Accept-Language", "en-US,en;q=0.5")
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   req.Header.Set("Origin", "https://account.unext.jp")
   req.Header.Set("Referer", "https://account.unext.jp/login?&backurl=https%3A%2F%2Fvideo.unext.jp%2Ftitle%2FSID0020149")
   req.Header.Set("Upgrade-Insecure-Requests", "1")

   // CheckRedirect can be left default, the CookieJar will store the _ut, _st, _upm, _ubr cookies
   resp, err := client.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusFound {
      return fmt.Errorf("login failed with status code: %d", resp.StatusCode)
   }

   return nil
}
