package stan

import (
   "bytes"
   "encoding/json"
   "fmt"
   "log"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

func do(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
}

type APIError struct {
   Code string
}

type APIErrors []APIError

func (e APIErrors) Error() string {
   codes := make([]string, len(e))
   for i, err := range e {
      codes[i] = err.Code
   }
   return fmt.Sprintf("stan: %s", strings.Join(codes, ", "))
}

// AppSession is the full session response from /login/v1/sessions/app.
// FIX: expanded to match Python _oauth() response fields.
type AppSession struct {
   JwToken string
   Renew   int64
   Now     int64
   UserId  string
   Profile Profile
   Errors  APIErrors
}

func (a *AppSession) FetchMedia(id int) (*Media, error) {
   req, err := http.NewRequest(
      "GET", "https://api.stan.com.au/concurrency/v1/streams", nil,
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("x-forwarded-for", "1.128.0.0")
   req.URL.RawQuery = url.Values{
      "capabilities.drm": {"widevine"}, // need for media.drm
      "format":           {"dash"},     // hls
      "jwToken":          {a.JwToken},
      "programId":        {strconv.Itoa(id)},

      //dragon high PASS
      //beast high PASS
      // beast ultra

      // high
      //play.stan.com.au/programs/331144

      //ultra
      //stan.com.au/watch/beast-2026
      //play.stan.com.au/programs/6299871

      //"quality":          {"sd"}, // 540p
      //"quality":          {"high"}, // 1080p
      "quality": {"ultra"}, // 2160p
      //"quality":          {"auto"}, // 2160p

   }.Encode()
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var result struct {
      Media  Media
      Errors APIErrors
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   if len(result.Errors) > 0 {
      return nil, result.Errors
   }

   return &result.Media, nil
}

type Media struct {
   Drm *struct {
      CustomData string
      KeyId      string
   }
   VideoUrl string
}

func (m *Media) License(data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      // final slash is needed
      "POST", "https://lic.drmtoday.com/license-proxy-widevine/cenc/",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("dt-custom-data", m.Drm.CustomData)
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      License []byte
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return result.License, nil
}

// Profile represents a Stan user profile.
type Profile struct {
   Id            string
   Name          string
   IsKidsProfile bool
}
