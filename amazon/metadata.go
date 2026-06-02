// metadata.go
package amazon

import (
   "context"
   "crypto/rand"
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
)

// GenerateAmazonRequestID generates a 20-character random alphanumeric uppercase string
func GenerateAmazonRequestID() string {
   const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
   b := make([]byte, 20)
   _, _ = rand.Read(b)
   for i := range b {
      b[i] = charset[int(b[i])%len(charset)]
   }
   return string(b)
}

// EnrichItemMetadata makes a request to /api/enrichItemMetadata to get playback actions including the PlaybackEnvelope.
func EnrichItemMetadata(ctx context.Context, client *http.Client, hostURL string, titleIDs []string) (map[string]interface{}, error) {
   endpoint := fmt.Sprintf("%s/api/enrichItemMetadata", hostURL)

   req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
   if err != nil {
      return nil, err
   }

   metadataToEnrich := map[string]string{
      "placement": "HOVER",
      "playback":  "true",
      "preroll":   "true",
      "trailer":   "true",
      "watchlist": "true",
   }

   metaBytes, err := json.Marshal(metadataToEnrich)
   if err != nil {
      return nil, err
   }

   titlesBytes, err := json.Marshal(titleIDs)
   if err != nil {
      return nil, err
   }

   u, err := url.Parse(hostURL)
   if err != nil {
      return nil, err
   }

   q := req.URL.Query()
   q.Add("metadataToEnrich", string(metaBytes))
   q.Add("titleIDsToEnrich", string(titlesBytes))
   q.Add("currentUrl", fmt.Sprintf("https://%s/", u.Host))

   req.URL.RawQuery = q.Encode()

   req.Header.Set("device-memory", "8")
   req.Header.Set("downlink", "10")
   req.Header.Set("dpr", "2")
   req.Header.Set("ect", "4g")
   req.Header.Set("rtt", "50")
   req.Header.Set("viewport-width", "671")
   req.Header.Set("x-amzn-client-ttl-seconds", "15")
   req.Header.Set("x-amzn-requestid", GenerateAmazonRequestID())
   req.Header.Set("x-requested-with", "XMLHttpRequest")

   resp, err := client.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   var result map[string]interface{}
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   return result, nil
}
