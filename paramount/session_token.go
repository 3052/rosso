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
   StreamingUrl string // MPD
   Url          string
   LsSession    string `json:"ls_session"`
}

// 1. do we always need to check streamingUrl ?

// 2. can androidphone ls_session be used with PlayReady ?

// 3. can xboxone ls_session be used with Widevine ?

// 4. can we hard code the license URL ?

// 5. do we actually need xboxone ?

// https://cbsi.live.ott.irdeto.com/playready/rightsmanager.asmx?
// AccountId=cbsi&
// ContentId=wjQ4RChi6BHHu4MVTncppVuCwu44uq2Q&
// CrmId=cbsi&
// SubContentType=Default
func (s *SessionToken) Send(body []byte) ([]byte, error) {
   req, err := http.NewRequest("POST", s.Url, bytes.NewReader(body))
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+s.LsSession)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

// - androidphone
//    - 2160p
//    - Widevine
// - xboxone
//    - 1080p
//    - PlayReady
func FetchSessionToken(at, contentId string, cbsCookie *http.Cookie) (*SessionToken, error) {
   endpoint := "anonymous-session-token.json"
   if cbsCookie != nil {
      endpoint = "session-token.json"
   }
   url_data := &url.URL{
      Scheme: "https",
      Host:   "www.paramountplus.com",
      Path:   fmt.Sprintf("/apps-api/v3.1/androidphone/irdeto-control/%s", endpoint),
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
   // I DONT THINK THIS IS AN ERROR IF WE JUST NEED ls_session
   if result.StreamingUrl == "" {
      return nil, errors.New("streaming URL is empty")
   }
   return &result, nil
}
