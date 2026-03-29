package paramount

import (
   "errors"
   "io"
   "net/http"
   "net/url"
   "testing"
   "time"
)

func TestCbs(t *testing.T) {
   results, err := ExtractDexHexBytes("base.apk")
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
      t.Log(brands("www.cbs.com", result), result)
   }
}

func TestParamount(t *testing.T) {
   results, err := ExtractDexHexBytes("base.apk")
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
      t.Log(brands("www.paramountplus.com", result), result)
   }
}

func brands(host, app_secret string) error {
   at, err := GetAt(app_secret)
   if err != nil {
      return err
   }
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{}
   req.URL.Host = host
   req.URL.Path = "/apps-api/v3.0/androidphone/brands/.json"
   value := url.Values{}
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
