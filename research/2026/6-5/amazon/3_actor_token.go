package amazon

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
)

type ActorTokenResponse struct {
   DeviceTokens []struct {
      DeviceType       string `json:"device_type"`
      ActorAccessToken struct {
         Token string `json:"token"`
      } `json:"actor_access_token"`
   } `json:"device_tokens"`
}

func GetActorAccessToken(refreshToken, actorId string) (string, error) {
   url := "https://api.amazon.com/auth/token"

   payload := map[string]interface{}{
      "app_name":             "com.amazon.avod.thirdpartyclient",
      "app_version":          "130050002",
      "source_token_type":    "refresh_token",
      "requested_token_type": "actor_access_token",
      "actor_id":             actorId,
      "source_device_tokens": []map[string]interface{}{
         {
            "device_type": "A43PXU4ZN2AL1",
            "account_refresh_token": map[string]string{
               "token": refreshToken,
            },
         },
      },
   }

   body, _ := json.Marshal(payload)
   req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
   if err != nil {
      return "", err
   }

   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("x-amzn-identity-auth-domain", "api.amazon.com")
   req.Header.Set("User-Agent", "AmazonWebView/MAPClientLib/130050002/Android/11/sdk_gphone_x86_64")

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   var tokenResp ActorTokenResponse
   if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
      return "", err
   }

   for _, dt := range tokenResp.DeviceTokens {
      if dt.DeviceType == "A43PXU4ZN2AL1" && dt.ActorAccessToken.Token != "" {
         return dt.ActorAccessToken.Token, nil
      }
   }

   return "", fmt.Errorf("actor access token not found in response")
}
