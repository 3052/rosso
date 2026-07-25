package peacock

import (
   "encoding/json"
   "errors"
   "net/http"
   "net/url"
   "strings"
)

var Territory = "US"

type IdSession struct {
   Cookie *http.Cookie
}

func FetchIdSession(user, password string) (*IdSession, error) {
   body := url.Values{
      "userIdentifier": {user},
      "password":       {password},
   }.Encode()
   target := url.URL{
      Scheme: "https",
      Host:   "rango.id.peacocktv.com",
      Path:   "/signin/service/international",
   }
   req, err := http.NewRequest(http.MethodPost, target.String(), strings.NewReader(body))
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   req.Header.Set("x-skyott-proposition", "NBCUOTT")
   req.Header.Set("x-skyott-provider", "NBCU")
   req.Header.Set("x-skyott-territory", Territory)
   resp, err := doRequest(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Properties SignInErrors
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusCreated {
      return nil, &result.Properties
   }
   for _, c := range resp.Cookies() {
      if c.Name == "idsession" {
         return &IdSession{Cookie: c}, nil
      }
   }
   return nil, errors.New("idsession cookie not present")
}

func (*IdSession) CachePath() string {
   return "rosso/peacock/IdSession"
}

type SignInErrors struct {
   Code   string
   Errors struct {
      CategoryErrors []struct {
         Code    string
         Message string
      }
   }
}

func (e *SignInErrors) Error() string {
   if e.Code != "" {
      return e.Code
   }
   var parts []string
   for _, ce := range e.Errors.CategoryErrors {
      parts = append(parts, ce.Code+": "+ce.Message)
   }
   return strings.Join(parts, "; ")
}
