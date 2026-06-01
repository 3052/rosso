// register_device.go
package amazon

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
)

type RegisterResponse struct {
   Response struct {
      Success struct {
         Tokens struct {
            Bearer struct {
               AccessToken string `json:"access_token"`
               ExpiresIn   int    `json:"expires_in"`
            } `json:"bearer"`
         } `json:"tokens"`
      } `json:"success"`
   } `json:"response"`
}

func RegisterDevice(client *http.Client, endpoint string, codePair *CodePairResponse, device map[string]string) (*RegisterResponse, error) {
   payload := map[string]interface{}{
      "auth_data": map[string]interface{}{
         "code_pair": codePair,
      },
      "registration_data":    device,
      "requested_token_type": []string{"bearer"},
      "requested_extensions": []string{"device_info", "customer_info"},
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

   if resp.StatusCode != http.StatusOK {
      respBody, _ := io.ReadAll(resp.Body)
      return nil, fmt.Errorf("unable to register: %s [%d]", string(respBody), resp.StatusCode)
   }

   var result RegisterResponse
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   return &result, nil
}
