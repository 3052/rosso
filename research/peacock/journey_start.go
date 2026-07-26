package peacock

import (
   "bytes"
   "context"
   "encoding/json"
   "fmt"
   "net/http"
)

// JourneyStartResponse is the response from starting a companion-device activation journey.
type JourneyStartResponse struct {
   JourneyID           string `json:"journeyId"`
   OneTimePassword     string `json:"oneTimePassword"`
   PollingPeriodSecs   int    `json:"pollingPeriodSecs"`
   RemainingOTPTTLSecs int    `json:"remainingOtpTtlSecs"`
   Status              string `json:"status"`
}

// StartJourney creates a new companion-device activation journey.
// The returned OneTimePassword is the 6-character code the user must enter at peacocktv.com/tv.
func (c *Client) StartJourney(ctx context.Context) (*JourneyStartResponse, error) {
   body, _ := json.Marshal(map[string]string{"deviceId": c.DeviceID})

   req, err := c.newRequest(ctx, http.MethodPost, sasBase+"/companion-service/journeys?type=sign-in", bytes.NewReader(body))
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/vnd.companionservice.v1+json")
   req.Header.Set("Accept", "application/vnd.companionservice.v1+json")

   resp, err := c.HTTP.Do(req)
   if err != nil {
      return nil, fmt.Errorf("start journey: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
      return nil, fmt.Errorf("start journey: bad status %d", resp.StatusCode)
   }

   var out JourneyStartResponse
   if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
      return nil, fmt.Errorf("start journey: decode: %w", err)
   }
   return &out, nil
}
