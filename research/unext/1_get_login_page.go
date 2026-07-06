package unext

import (
   "fmt"
   "io"
   "net/http"
   "regexp"
)

// GetLoginPage fetches the login page and extracts the _csrf token
func GetLoginPage(client *http.Client) (string, error) {
   url := "https://account.unext.jp/login?&backurl=https%3A%2F%2Fvideo.unext.jp%2Ftitle%2FSID0020149"

   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      return "", err
   }

   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
   req.Header.Set("Accept-Language", "en-US,en;q=0.5")
   req.Header.Set("Referer", "https://video.unext.jp/title/SID0020149")
   req.Header.Set("Upgrade-Insecure-Requests", "1")

   resp, err := client.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return "", err
   }

   // Extract _csrf token using regex from the hidden input
   re := regexp.MustCompile(`<input type="hidden" name="_csrf" value="(.*?)">`)
   matches := re.FindStringSubmatch(string(body))
   if len(matches) < 2 {
      return "", fmt.Errorf("could not find _csrf token on login page")
   }

   return matches[1], nil
}
