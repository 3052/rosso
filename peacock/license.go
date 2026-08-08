package peacock

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "log"
   "net/http"
)

// doRequest logs the request method and URL, then sends the request
// using the provided http.Client.
func doRequest(client *http.Client, req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return client.Do(req)
}

// PlayoutVodResponse is the response from POST /video/playouts/vod.
type PlayoutVodResponse struct {
   Asset struct {
      Endpoints []struct {
         Cdn string `json:"cdn"`
         Url string `json:"url"`
      } `json:"endpoints"`
   } `json:"asset"`
   Protection struct {
      LicenceAcquisitionUrl string `json:"licenceAcquisitionUrl"`
   } `json:"protection"`
   ErrorCode   string `json:"errorCode"`
   Description string `json:"description"`
}

// AcquireLicense sends a Widevine license acquisition request to the licence acquisition URL
// and returns the raw license bytes. The licenceAcquisitionUrl is the full URL returned in
// the PlayoutVod response.
func (p *PlayoutVodResponse) AcquireLicense(challenge []byte) ([]byte, error) {
   if p == nil {
      return nil, fmt.Errorf("acquire license: nil playout")
   }
   if p.Protection.LicenceAcquisitionUrl == "" {
      return nil, fmt.Errorf("acquire license: empty licenceAcquisitionUrl")
   }
   if len(challenge) == 0 {
      return nil, fmt.Errorf("acquire license: empty challenge")
   }
   req, err := http.NewRequest(http.MethodPost, p.Protection.LicenceAcquisitionUrl, bytes.NewReader(challenge))
   if err != nil {
      return nil, fmt.Errorf("acquire license: create request: %w", err)
   }
   resp, err := doRequest(http.DefaultClient, req)
   if err != nil {
      return nil, fmt.Errorf("acquire license: %w", err)
   }
   defer resp.Body.Close()

   license, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, fmt.Errorf("acquire license: read body: %w", err)
   }

   // The endpoint returns a JSON error body instead of a non-200 status
   // when the request is rejected (e.g. "Unsupported browser/client").
   var lr licenseErrorResponse
   if err := json.Unmarshal(license, &lr); err == nil && lr.ErrorCode != "" {
      return nil, fmt.Errorf("acquire license: %s: %s", lr.ErrorCode, lr.Description)
   }

   return license, nil
}

// licenseErrorResponse is the error body returned by the licence
// acquisition endpoint when the request is rejected.
type licenseErrorResponse struct {
   ErrorCode   string `json:"errorCode"`
   Description string `json:"description"`
}

// license.go
