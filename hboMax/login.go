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

type Cookie struct {
   Name  string
   Value string
}

func (c *Cookie) String() string {
   return fmt.Sprintf("%v=%v", c.Name, c.Value)
}

type Error struct {
   Code   string `json:"code"`
   Detail string `json:"detail"`
}

type Errors []Error

func (e Errors) Error() string {
   parts := make([]string, len(e))
   for i, err := range e {
      parts[i] = err.Detail + " (" + err.Code + ")"
   }
   return strings.Join(parts, "; ")
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

   resp, err := doReq(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var result struct {
      Errors Errors `json:"errors"`
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
