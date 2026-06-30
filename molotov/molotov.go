// constants.go
package molotov

import (
   "log"
   "net/http"
)

// DeviceID is the centralized value used for the x-device-id header across all requests.
const DeviceID = "x-device-id"

// doRequest logs the method and URL, then performs the HTTP request.
func doRequest(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   client := &http.Client{}
   return client.Do(req)
}
