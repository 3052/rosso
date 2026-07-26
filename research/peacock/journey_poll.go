package peacock

import (
   "context"
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
func (c *Client) PollJourney(ctx context.Context, journeyID string) (*JourneyStatusResponse, error) {
   req, err := c.newRequest(ctx, http.MethodGet, sasBase+"/companion-service/journeys/"+journeyID, nil)
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
func (c *Client) WaitForJourney(ctx context.Context, journeyID string, pollInterval time.Duration) (*JourneyStatusResponse, error) {
   ticker := time.NewTicker(pollInterval)
   defer ticker.Stop()

   for {
      status, err := c.PollJourney(ctx, journeyID)
      if err != nil {
         return nil, err
      }
      if status.Status == JourneyStatusCompleted || status.Status == JourneyStatusExpired {
         return status, nil
      }

      select {
      case <-ctx.Done():
         return nil, ctx.Err()
      case <-ticker.C:
      }
   }
}
