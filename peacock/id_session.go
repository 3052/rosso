// id_session.go
package peacock

import (
   "encoding/json"
   "errors"
   "net/http"
   "net/url"
   "strings"
)

var Territory = "US"

type Cookie struct {
   Name  string
   Value string
}

func FetchIdSession(user, password string) (*Cookie, error) {
   body := url.Values{
      "userIdentifier": {user},
      "password":       {password},
   }.Encode()
   target := url.URL{
      Scheme: "https",
      Host:   "rango.id.peacocktv.com",
      Path:   "/signin/service/international",
   }
   req, err := http.NewRequest("POST", target.String(), strings.NewReader(body))
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
      Properties struct {
         Errors struct {
            CategoryErrors []struct {
               Code string
            }
         }
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != 201 {
      return nil, errors.New(result.Properties.Errors.CategoryErrors[0].Code)
   }
   for _, c := range resp.Cookies() {
      if c.Name == "idsession" {
         return &Cookie{Name: c.Name, Value: c.Value}, nil
      }
   }
   return nil, errors.New("idsession cookie not present")
}
