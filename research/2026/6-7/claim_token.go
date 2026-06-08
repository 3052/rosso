// exchange_claim_token.go
package amazon

import (
   "fmt"
   "net/http"
   "net/url"
)

func ExchangeClaimToken(client *http.Client, claimTokenUrl string) (string, error) {
   req, err := http.NewRequest("GET", claimTokenUrl, nil)
   if err != nil {
      return "", fmt.Errorf("failed to create request: %w", err)
   }
   req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 11; sdk_gphone_x86_64 Build/RSR1.240422.006; wv) AppleWebKit/537.36")
   req.Header.Set("x-requested-with", "com.amazon.avod.thirdpartyclient")

   resp, err := client.Do(req)
   if err != nil {
      return "", fmt.Errorf("request failed: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusMovedPermanently {
      loc := resp.Header.Get("Location")

      parsedLoc, err := url.Parse(loc)
      if err != nil {
         return "", fmt.Errorf("failed to parse redirect location: %w", err)
      }

      // Amazon maps the oauth response directly into the redirect location query
      authCode := parsedLoc.Query().Get("openid.oa2.authorization_code")
      if authCode == "" {
         return "", fmt.Errorf("authorization_code missing from redirect location: %s", loc)
      }

      return authCode, nil
   }

   return "", fmt.Errorf("expected 302 redirect containing authorization_code, got %d", resp.StatusCode)
}
