package rtbf

import (
   "bytes"
   "encoding/json"
   "errors"
   "fmt"
   "net/http"
   "net/url"
)

func (i *Identity) Session() (*Session, error) {
   data, err := json.Marshal(map[string]any{
      "device": map[string]string{
         "deviceId": "",
         "type":     "WEB",
      },
      "jwt": i.IdToken,
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://exposure.api.redbee.live", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.URL.Path = "/v2/customer/RTBF/businessunit/Auvio/auth/gigyaLogin"
   req.Header.Set("content-type", "application/json")
   resp, err := http.DefaultClient.Do(req)
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

type Identity struct {
   ErrorMessage string
   IdToken      string `json:"id_token"`
}

type Session struct {
   SessionToken string
}

// hard coded in JavaScript
const api_key = "4_Ml_fJ47GnBAW6FrPzMxh0w"

func FetchAssetId(path string) (string, error) {
   resp, err := http.Get(
      "https://bff-service.rtbf.be/auvio/v1.23/pages" + path,
   )
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
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

func (a *Account) Identity() (*Identity, error) {
   resp, err := http.PostForm(
      "https://login.auvio.rtbf.be/accounts.getJWT", url.Values{
         "APIKey":      {api_key},
         "login_token": {a.SessionInfo.CookieValue},
      },
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Identity
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.ErrorMessage != "" {
      return nil, errors.New(result.ErrorMessage)
   }
   return &result, nil
}

func FetchAccount(id, password string) (*Account, error) {
   resp, err := http.PostForm(
      "https://login.auvio.rtbf.be/accounts.login", url.Values{
         "APIKey":   {api_key},
         "loginID":  {id},
         "password": {password},
      },
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Account
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.ErrorMessage != "" {
      return nil, errors.New(result.ErrorMessage)
   }
   return &result, nil
}

type Entitlement struct {
   AssetId   string
   Formats   []Format
   Message   string
   PlayToken string
}

type Format struct {
   Format       string
   MediaLocator string // MPD
}

func (s *Session) Entitlement(assetId string) (*Entitlement, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme: "https",
         Host:   "exposure.api.redbee.live",
         Path: fmt.Sprintf(
            "/v2/customer/RTBF/businessunit/Auvio/entitlement/%v/play", assetId,
         ),
      },
      Header: http.Header{},
   }
   req.Header.Set("x-forwarded-for", "91.90.123.17")
   req.Header.Set("authorization", "Bearer "+s.SessionToken)
   resp, err := http.DefaultClient.Do(&req)
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

func (e *Entitlement) GetDash() (*Format, error) {
   for _, format_data := range e.Formats {
      if format_data.Format == "DASH" {
         return &format_data, nil
      }
   }
   return nil, errors.New("DASH format not found")
}

func (f *Format) GetManifest() (*url.URL, error) {
   return url.Parse(f.MediaLocator)
}

func GetPath(urlData string) (string, error) {
   url_parse, err := url.Parse(urlData)
   if err != nil {
      return "", err
   }
   if url_parse.Scheme == "" {
      return "", errors.New("invalid URL: scheme is missing")
   }
   return url_parse.Path, nil
}

type Account struct {
   ErrorMessage string
   SessionInfo  struct {
      CookieValue string
   }
}
