package rtbf

import (
   "encoding/json"
   "errors"
   "fmt"
   "net/http"
   "net/url"
)

type Identity struct {
   ErrorMessage string
   IdToken      string `json:"id_token"`
}

type Session struct {
   SessionToken string
}

// hard coded in JavaScript
const api_key = "4_Ml_fJ47GnBAW6FrPzMxh0w"

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
