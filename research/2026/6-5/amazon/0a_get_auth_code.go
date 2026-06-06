package amazon

import (
   "fmt"
   "net/http"
   "net/url"
)

func GetAuthorizationCode(signInUrl string, cookies []*http.Cookie) (string, error) {
   req, err := http.NewRequest("GET", signInUrl, nil)
   if err != nil {
      return "", err
   }

   req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 11; sdk_gphone_x86_64 Build/RSR1.240422.006; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/83.0.4103.106 Mobile Safari/537.36")
   req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
   req.Header.Set("X-Requested-With", "com.amazon.avod.thirdpartyclient")

   for _, cookie := range cookies {
      req.AddCookie(cookie)
   }

   // Create a custom client that stops at the first redirect to capture the Location header
   client := &http.Client{
      CheckRedirect: func(req *http.Request, via []*http.Request) error {
         return http.ErrUseLastResponse
      },
   }

   resp, err := client.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusFound {
      return "", fmt.Errorf("expected 302 redirect, got status code: %d", resp.StatusCode)
   }

   location := resp.Header.Get("Location")
   if location == "" {
      return "", fmt.Errorf("location header not found in the 302 response")
   }

   parsedUrl, err := url.Parse(location)
   if err != nil {
      return "", err
   }

   authCode := parsedUrl.Query().Get("openid.oa2.authorization_code")
   if authCode == "" {
      return "", fmt.Errorf("openid.oa2.authorization_code missing from redirect location")
   }

   return authCode, nil
}
