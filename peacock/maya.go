package peacock

import (
   "encoding/json"
   "errors"
   "net/http"
   "net/url"
   "strings"
)

func FetchIdSession(user, password string) (*http.Cookie, error) {
   data := url.Values{
      "userIdentifier": {user},
      "password":       {password},
   }.Encode()
   req, err := http.NewRequest(
      "POST", "https://rango.id.peacocktv.com/signin/service/international",
      strings.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   req.Header.Set("x-skyott-proposition", "NBCUOTT")
   req.Header.Set("x-skyott-provider", "NBCU")
   req.Header.Set("x-skyott-territory", Territory)
   resp, err := http.DefaultClient.Do(req)
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
   if resp.StatusCode != http.StatusCreated {
      return nil, errors.New(result.Properties.Errors.CategoryErrors[0].Code)
   }
   for _, cookie := range resp.Cookies() {
      if cookie.Name == "idsession" {
         return cookie, nil
      }
   }
   return nil, http.ErrNoCookie
}
