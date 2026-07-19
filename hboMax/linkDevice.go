package hboMax

import (
   "encoding/json"
   "errors"
   "fmt"
   "log"
   "net/http"
)

const device_info = "!/!(!/!;!/!;!/!)"

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

type Initiate struct {
   LinkingCode string
   TargetUrl   string
}

func InitiateRequest(st *Cookie, market string) (*Initiate, error) {
   endpoint := fmt.Sprintf("https://default.beam-%v.prd.api.discomax.com/authentication/linkDevice/initiate", market)
   req, err := http.NewRequest(http.MethodPost, endpoint, nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("cookie", st.String())
   req.Header.Set("x-device-info", device_info)

   resp, err := doReq(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }

   var result struct {
      Data struct {
         Attributes Initiate
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Data.Attributes, nil
}
