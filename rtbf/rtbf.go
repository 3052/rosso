package rtbf

import (
   "bytes"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/url"
   "strings"
)

// hard coded in JavaScript
const api_key = "4_Ml_fJ47GnBAW6FrPzMxh0w"

// Belgium
const x_forwarded_for = "195.0.0.0"

func FetchAssetId(path string) (string, error) {
   req, err := http.NewRequest("GET",
      "https://bff-service.rtbf.be/auvio/v1.23/pages"+path,
      nil,
   )
   if err != nil {
      return "", err
   }
   resp, err := do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return "", errors.New(resp.Status)
   }
   var page struct {
      Data struct {
         Content struct {
            AssetId string
            Media   *struct {
               AssetId string
            }
         }
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&page)
   if err != nil {
      return "", err
   }
   content := page.Data.Content
   if content.AssetId != "" {
      return content.AssetId, nil
   }
   if content.Media != nil {
      return content.Media.AssetId, nil
   }
   return "", errors.New("assetId not found")
}

func GetPath(urlData string) (string, error) {
   parsedURL, err := url.Parse(urlData)
   if err != nil {
      return "", err
   }
   if parsedURL.Scheme == "" {
      return "", errors.New("invalid URL: scheme is missing")
   }
   return parsedURL.Path, nil
}

func do(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
}

type Account struct {
   SessionInfo struct {
      CookieValue string
   }
}

func FetchAccount(id, password string) (*Account, error) {
   body := url.Values{
      "APIKey":   {api_key},
      "loginID":  {id},
      "password": {password},
   }.Encode()
   req, err := http.NewRequest("POST",
      "https://login.auvio.rtbf.be/accounts.login",
      strings.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Account
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result, nil
}

func (*Account) CachePath() string {
   return "rosso/rtbf/Account"
}

func (a *Account) Identity() (*Identity, error) {
   body := url.Values{
      "APIKey":      {api_key},
      "login_token": {a.SessionInfo.CookieValue},
   }.Encode()
   req, err := http.NewRequest("POST",
      "https://login.auvio.rtbf.be/accounts.getJWT",
      strings.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Identity
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result, nil
}

type Entitlement struct {
   AssetId string
   Formats []struct {
      Format       string
      MediaLocator string // MPD
   }
   Message   string
   PlayToken string
}

func (*Entitlement) CachePath() string {
   return "rosso/rtbf/Entitlement"
}

func (e *Entitlement) FetchWidevine(body []byte) ([]byte, error) {
   params := url.Values{
      "contentId":  {e.AssetId},
      "ls_session": {e.PlayToken},
   }
   target := "https://exposure.api.redbee.live/v2/license/customer/RTBF/businessunit/Auvio/widevine?" +
      params.Encode()

   req, err := http.NewRequest("POST", target, bytes.NewReader(body))
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/x-protobuf")
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      var value struct {
         Message string
      }
      err = json.NewDecoder(resp.Body).Decode(&value)
      if err != nil {
         return nil, err
      }
      return nil, errors.New(value.Message)
   }
   return io.ReadAll(resp.Body)
}

func (e *Entitlement) GetDash() (*url.URL, error) {
   for _, format := range e.Formats {
      if format.Format == "DASH" {
         return url.Parse(format.MediaLocator)
      }
   }
   return nil, errors.New("DASH format not found")
}

type Identity struct {
   IdToken string `json:"id_token"`
}

func (i *Identity) Session() (*Session, error) {
   body, err := json.Marshal(map[string]any{
      "device": map[string]string{
         "deviceId": "",
         "type":     "WEB",
      },
      "jwt": i.IdToken,
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest("POST",
      "https://exposure.api.redbee.live/v2/customer/RTBF/businessunit/Auvio/auth/gigyaLogin",
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/json")
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &Session{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

type Session struct {
   SessionToken string
}

func (s *Session) Entitlement(assetId string) (*Entitlement, error) {
   target := fmt.Sprintf(
      "https://exposure.api.redbee.live/v2/customer/RTBF/businessunit/Auvio/entitlement/%s/play",
      assetId,
   )
   req, err := http.NewRequest("GET", target, nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+s.SessionToken)
   req.Header.Set("x-forwarded-for", x_forwarded_for)
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Entitlement
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   return &result, nil
}
