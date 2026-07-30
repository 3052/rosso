// journey_start.go
package peacock

import (
   "bytes"
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
func (c *Client) StartJourney() (*JourneyStartResponse, error) {
   body, _ := json.Marshal(map[string]string{"deviceId": c.DeviceID})

   req, err := http.NewRequest(http.MethodPost, sasBase+"/companion-service/journeys/sign-in", bytes.NewReader(body))
   if err != nil {
      return nil, err
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
