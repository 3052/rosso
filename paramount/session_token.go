package paramount

import (
   "bytes"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

type SessionToken struct {
   Errors       string
   LsSession    string `json:"ls_session"`
   StreamingUrl string // MPD
   Url          string
}

// do we always need to check streamingUrl? no

// can androidphone ls_session be used with PlayReady? no

// can we hard code the license URL? yes but its pointless because `url` must
// match `ls_session`

// do we actually need xboxone? yes for PlayReady

// should we cache session token? no because we need `androidphone` for MPD and
// `xboxone` for PlayReady

// what is `xboxone` MPD? 1080p

// what is `androidphone` MPD? 2160p

//-------------------------------------------------------------------------------

// what is the SL2000 max?

// what is the L3 max?

func FetchSessionToken(at, contentId string, cbsCookie *http.Cookie) (*SessionToken, error) {
   endpoint := "anonymous-session-token.json"
   if cbsCookie != nil {
      endpoint = "session-token.json"
   }
   url_data := &url.URL{
      Scheme: "https",
      Host:   "www.paramountplus.com",
      Path:   fmt.Sprintf("/apps-api/v3.1/xboxone/irdeto-control/%s", endpoint),
   }
   query := url_data.Query()
   query.Set("at", at)
   query.Set("contentId", contentId)
   url_data.RawQuery = query.Encode()
   req, err := http.NewRequest(http.MethodGet, url_data.String(), nil)
   if err != nil {
      return nil, err
   }
   if cbsCookie != nil {
      req.AddCookie(cbsCookie)
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result SessionToken
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Errors != "" {
      return nil, errors.New(result.Errors)
   }
   if result.StreamingUrl == "" {
      return nil, errors.New("streaming URL is empty")
   }
   return &result, nil
}
