package paramount

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/url"
)

func (a *App) fetch_session(platform, contentId, cbs_com string) (*Session, error) {
   at, err := get_at(a.Secret)
   if err != nil {
      return nil, err
   }
   endpoint := "anonymous-session-token.json"
   header := map[string]string{}
   if cbs_com != "" {
      endpoint = "session-token.json"
      header["cookie"] = cbs_com
   }
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   a.Host,
         Path:   fmt.Sprintf("/apps-api/v3.1/%s/irdeto-control/%s", platform, endpoint),
         RawQuery: url.Values{
            "at":        {at},
            "contentId": {contentId},
         }.Encode(),
      },
      header,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Session
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   return &result, nil
}

// WARNING IF YOU RUN THIS TOO MANY TIMES YOU WILL GET AN IP BAN. HOWEVER THE BAN
// IS ONLY FOR THE ANDROID CLIENT NOT WEB CLIENT
func (a *App) FetchCbsCom(username, password string) (string, error) {
   at, err := get_at(a.Secret)
   if err != nil {
      return "", err
   }
   body := url.Values{
      "j_username": {username},
      "j_password": {password},
   }.Encode()
   resp, err := maya.Post(
      &url.URL{
         Scheme:   "https",
         Host:     a.Host,
         Path:     "/apps-api/v2.0/androidphone/auth/login.json",
         RawQuery: url.Values{"at": {at}}.Encode(),
      },
      map[string]string{
         "content-type": "application/x-www-form-urlencoded",
         "user-agent":   "!", // randomly fails if this is missing
      },
      []byte(body),
   )
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   var result struct {
      Message string
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return "", err
   }
   if result.Message != "" {
      return "", errors.New(result.Message)
   }
   for _, cookie := range resp.Cookies() {
      if cookie.Name == "CBS_COM" {
         return cookie.String(), nil
      }
   }
   return "", errors.New("named cookie not present")
}

func (s *Session) Fetch(body []byte) ([]byte, error) {
   target, err := url.Parse(s.Url)
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      target, map[string]string{"authorization": "Bearer " + s.LsSession}, body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}
