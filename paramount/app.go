package paramount

import (
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/url"
)

var Apps = []App{
   {
      Id:      "com.cbs.app",
      Host:    "www.paramountplus.com",
      Secret:  "7081400bd4143bf3",
      Version: "Paramount+ 16.8.0",
   },
   {
      Id:      "com.cbs.tve",
      Host:    "www.cbs.com",
      Secret:  "cef32931dc01412e",
      Version: "CBS 15.6.0",
   },
   {
      Id:      "com.cbs.ca",
      Host:    "www.paramountplus.com",
      Secret:  "1c5d27627d71b420",
      Version: "Paramount+ 16.8.0",
   },
}

type App struct {
   Id      string
   Host    string
   Secret  string
   Version string
}

func GetApp(id string) (*App, error) {
   for _, each := range Apps {
      if each.Id == id {
         return &each, nil
      }
   }
   return nil, fmt.Errorf("app not found %q", id)
}

// WARNING IF YOU RUN THIS TOO MANY TIMES YOU WILL GET AN IP BAN. HOWEVER THE BAN
// IS ONLY FOR THE ANDROID CLIENT NOT WEB CLIENT
func (a *App) FetchCbsCom(username, password string) (*Cookie, error) {
   at, err := get_at(a.Secret)
   if err != nil {
      return nil, err
   }
   body := url.Values{
      "j_username": {username},
      "j_password": {password},
   }.Encode()
   target := &url.URL{
      Scheme:   "https",
      Host:     a.Host,
      Path:     "/apps-api/v2.0/androidphone/auth/login.json",
      RawQuery: url.Values{"at": {at}}.Encode(),
   }
   req, err := newPostRequest(
      target.String(),
      map[string]string{
         "content-type": "application/x-www-form-urlencoded",
         "user-agent":   "!", // randomly fails if this is missing
      },
      []byte(body),
   )
   if err != nil {
      return nil, err
   }
   resp, err := doRequest(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   _, err = io.Copy(io.Discard, resp.Body)
   if err != nil {
      return nil, err
   }
   for _, each := range resp.Cookies() {
      if each.Name == "CBS_COM" {
         return &Cookie{Name: each.Name, Value: each.Value}, nil
      }
   }
   return nil, errors.New("CBS_COM cookie not present")
}

func (a *App) FetchPlayReady(contentId string, cbsCom *Cookie) (*Session, error) {
   return a.fetch_session("xboxone", contentId, cbsCom)
}

func (a *App) FetchStreamingUrl(contentId string, cbsCom *Cookie) (*Session, error) {
   result, err := a.fetch_session("androidphone", contentId, cbsCom)
   if err != nil {
      return nil, err
   }
   if result.StreamingUrl == "" {
      return nil, errors.New("streamingUrl (MPD) is missing")
   }
   return result, nil
}

func (a *App) FetchWidevine(contentId string, cbsCom *Cookie) (*Session, error) {
   return a.fetch_session("androidphone", contentId, cbsCom)
}

func (a *App) fetch_session(platform, contentId string, cbs_com *Cookie) (*Session, error) {
   at, err := get_at(a.Secret)
   if err != nil {
      return nil, err
   }
   endpoint := "anonymous-session-token.json"
   header := map[string]string{}
   if cbs_com != nil {
      endpoint = "session-token.json"
      header["cookie"] = cbs_com.String()
   }
   target := &url.URL{
      Scheme: "https",
      Host:   a.Host,
      Path:   fmt.Sprintf("/apps-api/v3.1/%s/irdeto-control/%s", platform, endpoint),
      RawQuery: url.Values{
         "at":        {at},
         "contentId": {contentId},
      }.Encode(),
   }
   req, err := newGetRequest(target.String(), header)
   if err != nil {
      return nil, err
   }
   resp, err := doRequest(req)
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

type Cookie struct {
   Name  string
   Value string
}

func (*Cookie) CachePath() string {
   return "rosso/paramount/Cookie"
}

func (c *Cookie) String() string {
   return fmt.Sprintf("%v=%v", c.Name, c.Value)
}

type Session struct {
   LsSession    string `json:"ls_session"`
   Message      string
   StreamingUrl string // MPD
   Url          string // License Server
}

func (s *Session) Fetch(body []byte) ([]byte, error) {
   req, err := newPostRequest(
      s.Url, map[string]string{"authorization": "Bearer " + s.LsSession}, body,
   )
   if err != nil {
      return nil, err
   }
   resp, err := doRequest(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

// app.go
