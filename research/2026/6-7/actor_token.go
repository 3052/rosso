// get_actor_token.go
package amazon

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
)

type ActorTokenRequest struct {
   AppName            string                   `json:"app_name"`
   AppVersion         string                   `json:"app_version"`
   SourceTokenType    string                   `json:"source_token_type"`
   RequestedTokenType string                   `json:"requested_token_type"`
   SourceDeviceTokens []SourceDeviceTokenParam `json:"source_device_tokens"`
   ActorId            string                   `json:"actor_id"`
   AgeInfo            map[string]any           `json:"age_info"`
}

type SourceDeviceTokenParam struct {
   DeviceType          string              `json:"device_type"`
   AccountRefreshToken AccountRefreshToken `json:"account_refresh_token"`
}

type AccountRefreshToken struct {
   Token string `json:"token"`
}

type ActorTokenResponse struct {
   ActorType    string                `json:"actor_type"`
   TokenType    string                `json:"token_type"`
   DeviceTokens []DeviceTokenResponse `json:"device_tokens"`
}

type DeviceTokenResponse struct {
   DeviceType       string `json:"device_type"`
   ActorAccessToken struct {
      ExpiresIn int    `json:"expires_in"`
      Token     string `json:"token"`
   } `json:"actor_access_token"`
}

func GetActorToken(client *http.Client, refreshToken, actorId, deviceType string) (*ActorTokenResponse, error) {
   reqBody := ActorTokenRequest{
      AppName:            "com.amazon.avod.thirdpartyclient",
      AppVersion:         "130050002",
      SourceTokenType:    "refresh_token",
      RequestedTokenType: "actor_access_token",
      SourceDeviceTokens: []SourceDeviceTokenParam{
         {
            DeviceType: deviceType, // e.g. "A43PXU4ZN2AL1"
            AccountRefreshToken: AccountRefreshToken{
               Token: refreshToken,
            },
         },
      },
      ActorId: actorId,
      AgeInfo: make(map[string]any),
   }

   jsonData, err := json.Marshal(reqBody)
   if err != nil {
      return nil, fmt.Errorf("failed to marshal request: %w", err)
   }

   req, err := http.NewRequest("POST", "https://api.amazon.com/auth/token", bytes.NewBuffer(jsonData))
   if err != nil {
      return nil, fmt.Errorf("failed to create request: %w", err)
   }

   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("User-Agent", "AmazonWebView/MAPClientLib/130050002/Android/11/sdk_gphone_x86_64")
   req.Header.Set("x-amzn-identity-auth-domain", "api.amazon.com")

   resp, err := client.Do(req)
   if err != nil {
      return nil, fmt.Errorf("request failed: %w", err)
   }
   defer resp.Body.Close()

   bodyBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, fmt.Errorf("failed to read response: %w", err)
   }

   var result ActorTokenResponse
   if err := json.Unmarshal(bodyBytes, &result); err != nil {
      return nil, fmt.Errorf("failed to decode response: %w", err)
   }

   return &result, nil
}
