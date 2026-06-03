// codepair.go
package amazon

import (
   "bytes"
   "context"
   "encoding/json"
   "fmt"
   "net/http"
)

// CreateCodePair makes a request to /auth/create/codepair
func CreateCodePair(ctx context.Context, client *http.Client, apiBaseURL string, device map[string]interface{}) (map[string]interface{}, error) {
   url := fmt.Sprintf("https://%s/auth/create/codepair", apiBaseURL)

   payload := map[string]interface{}{
      "code_data": device,
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
   if err != nil {
      return nil, err
   }

   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("Accept-Language", "en-US")

   resp, err := client.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var result map[string]interface{}
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   if _, hasError := result["error"]; hasError {
      return nil, fmt.Errorf("api error: %v - %v", result["error"], result["error_description"])
   }

   return result, nil
}
