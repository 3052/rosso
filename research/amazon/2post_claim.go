// post_claim.go
package amazon

import (
   "fmt"
   "io/ioutil"
   "net/http"
   "net/url"
   "strings"
)

// PostClaim submits the email address (the "claim").
func PostClaim(client *http.Client, pageData *PageData, email string) (*PageData, error) {
   if pageData == nil || pageData.ActionURL == "" {
      return nil, fmt.Errorf("invalid page data or missing action URL")
   }
   // Explicitly check for the required workflowState parameter and return an error if missing
   ws, ok := pageData.HiddenParams["workflowState"]
   if !ok || ws == "" {
      return nil, fmt.Errorf("missing required hidden parameter: workflowState")
   }
   data := url.Values{}
   data.Set("workflowState", ws)
   req, err := http.NewRequest("POST", pageData.ActionURL, strings.NewReader(data.Encode()))
   if err != nil {
      return nil, fmt.Errorf("error creating request: %w", err)
   }
   // Only User-Agent and Content-Type needed for headers
   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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

   return ExtractPageData(string(bodyBytes)), nil
}
