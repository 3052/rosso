// client_do.go
package unext

import (
   "log"
   "net/http"
)

// clientDo wraps client.Do with a log line so every request is visible.
func clientDo(client *http.Client, req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return client.Do(req)
}
