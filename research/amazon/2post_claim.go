// post_claim.go
package amazon

import (
   "fmt"
   "io/ioutil"
   "net/http"
   "net/url"
   "strings"
)

// PostClaim submits the email address (the "claim") using the dynamic Action URL and hidden tokens.
func PostClaim(client *http.Client, pageData *PageData, email string) (*PageData, error) {
   if pageData == nil || pageData.ActionURL == "" {
      return nil, fmt.Errorf("invalid page data or missing action URL")
   }

   data := url.Values{}
   // Populate dynamically extracted hidden fields required by Amazon
   for k, v := range pageData.HiddenParams {
      data.Set(k, v)
   }
   data.Set("email", email)

   req, err := http.NewRequest("POST", pageData.ActionURL, strings.NewReader(data.Encode()))
   if err != nil {
      return nil, fmt.Errorf("error creating request: %w", err)
   }

   req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
   req.Header.Set("accept-language", "en-US,en;q=0.5")
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   req.Header.Set("origin", "https://www.amazon.com")
   req.Header.Set("upgrade-insecure-requests", "1")
   req.Header.Set("sec-fetch-dest", "document")
   req.Header.Set("sec-fetch-mode", "navigate")
   req.Header.Set("sec-fetch-site", "same-origin")
   req.Header.Set("sec-fetch-user", "?1")

   resp, err := client.Do(req)
   if err != nil {
      return nil, fmt.Errorf("error executing request: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   bodyBytes, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      return nil, fmt.Errorf("error reading response body: %w", err)
   }

   // Extract the new action URL and tokens for the password step
   return ExtractPageData(string(bodyBytes)), nil
}
