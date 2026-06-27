package amazon

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
)

// ActorToken represents an actor-specific access token.
type ActorToken struct {
   Token string `json:"token"`
}

// GetActorToken exchanges the account refresh token and actor ID for an actor-specific access token.
func GetActorToken(refreshToken, actorId string) (*ActorToken, error) {
   payload := map[string]any{
      "actor_id":             actorId,
      "app_name":             "AIV",
      "requested_token_type": "actor_access_token",
      "source_token_type":    "refresh_token",
      "source_device_tokens": []any{
         map[string]any{
            "device_type": DeviceTypeID,
            "account_refresh_token": map[string]string{
               "token": refreshToken,
            },
         },
      },
   }
   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", HostAmazonAPI+"/auth/token", bytes.NewBuffer(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/json")
   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   // Embed our new ActorToken struct into the anonymous decoder struct
   var result struct {
      DeviceTokens []struct {
         ActorAccessToken ActorToken `json:"actor_access_token"`
      } `json:"device_tokens"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   if len(result.DeviceTokens) == 0 {
      return nil, fmt.Errorf("no device tokens returned")
   }

   token := result.DeviceTokens[0].ActorAccessToken
   return &token, nil
}

func (*ActorToken) CachePath() string {
   return "rosso/amazon/ActorToken"
}
