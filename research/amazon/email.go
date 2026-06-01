// post_email.go
package amazon

import (
   "fmt"
   "io"
   "net/http"
   "net/url"
   "os"
   "strings"
)

func PostEmail(s *Session, action string, inputs map[string]string) (string, map[string]string, error) {
   data := url.Values{}
   for k, v := range inputs {
      data.Set(k, v)
   }
   data.Set("email", s.Email)

   req, err := http.NewRequest("POST", action, strings.NewReader(data.Encode()))
   if err != nil {
      return "", nil, err
   }

   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
   req.Header.Set("Accept-Language", "en-US,en;q=0.5")
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   req.Header.Set("Origin", "https://www.amazon.com")
   req.Header.Set("Upgrade-Insecure-Requests", "1")
   req.Header.Set("Sec-Fetch-Dest", "document")
   req.Header.Set("Sec-Fetch-Mode", "navigate")
   req.Header.Set("Sec-Fetch-Site", "same-origin")
   req.Header.Set("Referer", "https://www.amazon.com/ap/signin")

   resp, err := s.Client.Do(req)
   if err != nil {
      return "", nil, err
   }
   defer resp.Body.Close()

   body, _ := io.ReadAll(resp.Body)
   nextAction, nextInputs := ExtractForm(string(body), "signIn")

   if nextAction == "" {
      os.WriteFile("error_post_email.html", body, 0644)
      return "", nil, fmt.Errorf("password form not found. Response body saved to error_post_email.html")
   }

   if strings.HasPrefix(nextAction, "/") {
      nextAction = "https://www.amazon.com" + nextAction
   }

   return nextAction, nextInputs, nil
}
