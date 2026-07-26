package peacock

import (
   "bytes"
   "context"
   "encoding/json"
   "fmt"
   "net/http"
)

type activateRequest struct {
   DeviceID  string `json:"deviceId"`
   JourneyID string `json:"journeyId"`
}

// Activate exchanges the completed journeyId for the OAuth2 user token.
// The token is stored in c.Token and returned as a string.
func (c *Client) Activate(ctx context.Context, journeyID string) (string, error) {
   body, _ := json.Marshal(activateRequest{
      DeviceID:  c.DeviceID,
      JourneyID: journeyID,
   })

   req, err := c.newRequest(ctx, http.MethodPost, sasBase+"/commerce/activation/activate?type=sign-in", bytes.NewReader(body))
   if err != nil {
      return "", err
   }
   req.Header.Set("Content-Type", "application/vnd.activationservice.v1+json")
   req.Header.Set("Accept", "application/vnd.activationservice.v1+json")

   resp, err := c.HTTP.Do(req)
   if err != nil {
      return "", fmt.Errorf("activate: %w", err)
   }
   defer resp.Body.Close()

   // We expect a 204 No Content response
   if resp.StatusCode != http.StatusNoContent {
      return "", fmt.Errorf("activate: bad status %d", resp.StatusCode)
   }

   token := resp.Header.Get("x-skyid-token")
   if token == "" {
      return "", fmt.Errorf("activate: response missing x-skyid-token header")
   }

   c.Token = token
   return token, nil
}
