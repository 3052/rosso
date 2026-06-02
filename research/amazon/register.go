// register.go
package amazon

import (
   "bytes"
   "context"
   "encoding/json"
   "fmt"
   "net/http"
)

// RegisterDevice makes a request to /auth/register
func RegisterDevice(ctx context.Context, client *http.Client, apiBaseURL string, codePair map[string]interface{}, device map[string]interface{}) (map[string]interface{}, error) {
   url := fmt.Sprintf("https://%s/auth/register", apiBaseURL)

   payload := map[string]interface{}{
      "auth_data": map[string]interface{}{
         "code_pair": codePair,
      },
      "registration_data":    device,
      "requested_token_type": []string{"bearer"},
      "requested_extensions": []string{"device_info", "customer_info"},
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
