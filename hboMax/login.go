package hboMax

import (
   "encoding/json"
   "fmt"
   "log"
   "net/http"
   "strings"
)

// doReq handles executing the HTTP request and logging the method/URL
func doReq(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
}

// APIError represents a single error object from the Max API
type APIError struct {
   Code   string `json:"code"`
   Detail string `json:"detail"`
}

// APIErrors represents a collection of API errors and implements the error interface
type APIErrors []APIError

func (e APIErrors) Error() string {
   var b strings.Builder
   for i, err := range e {
      if i > 0 {
         b.WriteString(", ")
      }
      b.WriteString(err.Code)
      b.WriteString(": ")
      b.WriteString(err.Detail)
   }
   return b.String()
}

type Cookie struct {
   Name  string
   Value string
}

func (c *Cookie) String() string {
   return fmt.Sprintf("%v=%v", c.Name, c.Value)
}

type Login struct {
   Token string
}

// you must
// /authentication/linkDevice/initiate
// first or this will always fail
func LoginRequest(st *Cookie) (*Login, error) {
   req, err := http.NewRequest(http.MethodPost, "https://default.prd.api.hbomax.com/authentication/linkDevice/login", nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("cookie", st.String())
   req.Header.Set("x-device-info", device_info)
   req.Header.Set("x-disco-client", disco_client)
   req.Header.Set("x-disco-params", disco_params)

   resp, err := doReq(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var result struct {
      Errors APIErrors `json:"errors"`
      Data   struct {
         Attributes Login `json:"attributes"`
      } `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   if len(result.Errors) > 0 {
      return nil, result.Errors
   }
   return &result.Data.Attributes, nil
}
