package peacock

import (
   "encoding/json"
   "fmt"
   "net/http"
)

// JourneyStatus enumerates the values of `status` returned by the polling endpoint.
const (
   JourneyStatusNotStarted = "NOT_STARTED"
   JourneyStatusStarted    = "STARTED"
   JourneyStatusCompleted  = "COMPLETED"
   JourneyStatusExpired    = "EXPIRED"
)

// JourneyStatusResponse is the body returned by GET /companion-service/journeys/{journeyId}.
type JourneyStatusResponse struct {
   PollingPeriodSecs    int    `json:"pollingPeriodSecs"`
   RemainingOTPTTLSsecs int    `json:"remainingOtpTtlSecs"`
   OneTimePassword      string `json:"oneTimePassword"`
   Status               string `json:"status"`
   HouseholdID          string `json:"householdId,omitempty"`
}

// PollJourney fetches the current state of an activation journey exactly once.
func (c *Client) PollJourney(journeyID string) (*JourneyStatusResponse, error) {
   req, err := c.newRequest(http.MethodGet, sasBase+"/companion-service/journeys/"+journeyID, nil)
   if err != nil {
      return nil, err
   }

   resp, err := c.HTTP.Do(req)
   if err != nil {
      return nil, fmt.Errorf("poll journey: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("poll journey: bad status %d", resp.StatusCode)
   }

   var out JourneyStatusResponse
   if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
      return nil, fmt.Errorf("poll journey: decode: %w", err)
   }
   return &out, nil
}
