package stan

import (
   "encoding/json"
   "log"
   "net/http"
   "net/url"
   "strings"
)

func do(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
}

type ApiError struct {
   Code string
}

type ApiErrors []ApiError

func (e ApiErrors) Error() string {
   var builder strings.Builder
   builder.WriteString("stan: ")
   for i, err := range e {
      if i > 0 {
         builder.WriteString(", ")
      }
      builder.WriteString(err.Code)
   }
   return builder.String()
}

// AppSession is the full session response from /login/v1/sessions/app.
type AppSession struct {
   JwToken string
   Renew   int64
   Now     int64
   UserId  string
   Errors  ApiErrors
}

// WebToken is the intermediate token returned from the activation URL.
// It is used as the payload for the session request.
type WebToken struct {
   JwToken string
}

// FetchSession exchanges the WebToken for a full AppSession.
//
//   includes device data in the payload,
//   and handles error responses.
func (w *WebToken) FetchSession(hdr bool) (*AppSession, error) {

   params := url.Values{
      "hdcpVersion":  {"2.3"}, // need for UHD
      "jwToken":      {w.JwToken},
      "manufacturer": {"NVIDIA"},            // need for UHD
      "model":        {"SHIELD Android TV"}, // need for UHD
      "screenSize":   {"9999x9999"},         // need for UHD
      "stanName":     {"Stan-AndroidTV"},
      "stanVersion":  {"9"}, // need for UHD
   }
   if hdr {
      params.Set("colorSpace", "hdr10")
   }

   req, err := http.NewRequest(
      "POST", "https://api.stan.com.au/login/v1/sessions/app",
      strings.NewReader(params.Encode()),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   result := &AppSession{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }

   if len(result.Errors) > 0 {
      return nil, result.Errors
   }

   return result, nil
}

// session.go
