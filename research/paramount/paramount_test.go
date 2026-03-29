package paramount

import (
   "errors"
   "io"
   "net/http"
   "net/url"
   "os"
   "testing"
   "time"
)

func TestParamount(t *testing.T) {
   data, err := os.ReadFile("base.apk")
   if err != nil {
      t.Fatal(err)
   }
   results, err := ExtractDexHexBytes(data)
   if err != nil {
      t.Fatal(err)
   }
   var sleep bool
   for result := range results {
      if sleep {
         time.Sleep(time.Second)
      } else {
         sleep = true
      }
      t.Log(brands(result), result)
   }
}

func brands(app_secret string) error {
   at, err := GetAt(app_secret)
   if err != nil {
      return err
   }
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{}
   req.URL.Host = "www.cbs.com"
   //req.URL.Host = "www.paramountplus.com"
   req.URL.Path = "/apps-api/v3.0/androidphone/brands/.json"
   value := url.Values{}
   //value["at"] = []string{"ABBvXt8DjsbrJPs2Ry2E5I74VhdrYmExRFhsT9js5qJ5y0mgFLiEriZ9KYt0oQ08yPA="}
   value["at"] = []string{at}
   req.URL.RawQuery = value.Encode()
   req.URL.Scheme = "https"
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   _, err = io.Copy(io.Discard, resp.Body)
   if err != nil {
      return err
   }
   if resp.StatusCode != http.StatusOK {
      return errors.New(resp.Status)
   }
   return nil
}
