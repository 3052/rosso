package peacock

import (
   "encoding/json"
   "fmt"
   "net/http"
   "time"
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
   PollingPeriodSecs   int    `json:"pollingPeriodSecs"`
   RemainingOTPTTLSecs int    `json:"remainingOtpTtlSecs"`
   OneTimePassword     string `json:"oneTimePassword"`
   Status              string `json:"status"`
   HouseholdID         string `json:"householdId,omitempty"`
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

// WaitForJourney polls until the journey reaches a terminal state (COMPLETED / EXPIRED).
// A timeout is required to prevent indefinite blocking.
func (c *Client) WaitForJourney(journeyID string, pollInterval, timeout time.Duration) (*JourneyStatusResponse, error) {
   ticker := time.NewTicker(pollInterval)
   defer ticker.Stop()

   timeoutTimer := time.NewTimer(timeout)
   defer timeoutTimer.Stop()

   for {
      status, err := c.PollJourney(journeyID)
      if err != nil {
         return nil, err
      }
      if status.Status == JourneyStatusCompleted || status.Status == JourneyStatusExpired {
         return status, nil
      }

      select {
      case <-timeoutTimer.C:
         return nil, fmt.Errorf("wait for journey: timed out after %s", timeout)
      case <-ticker.C:
      }
   }
}
