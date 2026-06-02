package amazon

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
)

type CodePairResponse struct {
   PublicCode       string `json:"public_code"`
   PrivateCode      string `json:"private_code"`
   Error            string `json:"error,omitempty"`
   ErrorDescription string `json:"error_description,omitempty"`
}

func GetCodePair(client *http.Client, endpoint string, device map[string]string) (*CodePairResponse, error) {
   payload := map[string]interface{}{
      "code_data": device,
   }
   bodyBytes, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(bodyBytes))
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

   var result CodePairResponse
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   if result.Error != "" {
      return nil, fmt.Errorf("unable to get code pair: %s [%s]", result.ErrorDescription, result.Error)
   }

   return &result, nil
}
