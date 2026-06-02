package amazon

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
)

type CodePairResponse struct {
   PublicCode       string `json:"public_code"`
   PrivateCode      string `json:"private_code"`
   Error            string `json:"error,omitempty"`
   ErrorDescription string `json:"error_description,omitempty"`
}

type RegisterResponse struct {
   Response struct {
      Success struct {
         Tokens struct {
            Bearer struct {
               AccessToken  string `json:"access_token"`
               RefreshToken string `json:"refresh_token"`
               ExpiresIn    string `json:"expires_in"` // Set to string based on Amazon's JSON response
            } `json:"bearer"`
         } `json:"tokens"`
      } `json:"success"`
   } `json:"response"`
   Error            string `json:"error,omitempty"`
   ErrorDescription string `json:"error_description,omitempty"`
}

// Define the device identity we are pretending to be
var defaultDevice = map[string]string{
   "domain":        "Device",
   "app_name":      "com.amazon.amazonvideo.livingroom",
   "app_version":   "1.1",
   "device_model":  "LG-Tv",
   "os_version":    "6.0.1",
   "device_type":   "A71I8788P1ZV8",
   "device_name":   "My Go Device",
   "device_serial": "a906a7f9bfd6a7ab",
}

func GetCodePair() (*CodePairResponse, error) {
   bodyBytes, err := json.Marshal(map[string]any{
      "code_data": defaultDevice,
   })
   if err != nil {
      return nil, err
   }
   resp, err := http.Post(
      "https://api.amazon.com/auth/create/codepair", "",
      bytes.NewBuffer(bodyBytes),
   )
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

func RegisterDevice(codePair *CodePairResponse) (*RegisterResponse, error) {
   bodyBytes, err := json.Marshal(map[string]any{
      "auth_data": map[string]any{
         "code_pair": codePair,
      },
      "registration_data":    defaultDevice,
      "requested_token_type": []string{"bearer"},
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      http.MethodPost, "https://api.amazon.com/auth/register",
      bytes.NewBuffer(bodyBytes),
   )
   if err != nil {
      return nil, err
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      respBody, _ := io.ReadAll(resp.Body)
      return nil, fmt.Errorf("unable to register (has the code been entered?): %s [%d]", string(respBody), resp.StatusCode)
   }

   var result RegisterResponse
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   if result.Error != "" {
      return nil, fmt.Errorf("API error: %s [%s]", result.ErrorDescription, result.Error)
   }

   return &result, nil
}
