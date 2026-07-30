// activate.go
package peacock

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
)

type activateRequest struct {
   DeviceID  string `json:"deviceId"`
   JourneyID string `json:"journeyId"`
}

// Activate exchanges the completed journeyId for the OAuth2 user token.
func (c *Client) Activate(journeyID string) (string, error) {
   body, _ := json.Marshal(activateRequest{
      DeviceID:  c.DeviceID,
      JourneyID: journeyID,
   })

   req, err := http.NewRequest(http.MethodPost, sasBase+"/commerce/activation/activate?type=sign-in", bytes.NewReader(body))
   if err != nil {
      return "", err
   }
   req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 12; sdk_gphone64_x86_64 Build/SE1A.220826.008; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/91.0.4472.114 Mobile Safari/537.36")
   req.Header.Set("x-skyott-platform", "ANDROIDTV")
   req.Header.Set("x-skyott-proposition", "NBCUOTT")
   req.Header.Set("x-skyott-provider", "NBCU")
   req.Header.Set("x-skyott-territory", "US")
   req.Header.Set("x-skyott-activeterritory", "US")
   req.Header.Set("x-skyott-language", "en-US")
   req.Header.Set("x-skyott-device", "TV")
   req.Header.Set("x-skyott-broadcastregions", "INPATTERN_US_CENTRAL")
   req.Header.Set("x-skyint-requestid", randomUUID())
   req.Header.Set("Content-Type", "application/vnd.activationservice.v1+json")
   req.Header.Set("Accept", "application/vnd.activationservice.v1+json")

   resp, err := c.HTTP.Do(req)
   if err != nil {
      return "", fmt.Errorf("activate: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusNoContent {
      return "", fmt.Errorf("activate: bad status %d", resp.StatusCode)
   }

   token := resp.Header.Get("x-skyid-token")
   if token == "" {
      return "", fmt.Errorf("activate: response missing x-skyid-token header")
   }

   return token, nil
}
