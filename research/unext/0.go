package main

import (
   "bytes"
   "net/http"
   "net/url"
   "os"
)

func main() {
   client := &http.Client{}

   reqURL := &url.URL{
      Scheme: "https",
      Host:   "myaccount.unext.jp",
      Path:   "/api/migrateTokens",
   }
   bodyData := []byte(`{"forceMigration":true}`)
   req, err := http.NewRequest("POST", reqURL.String(), bytes.NewBuffer(bodyData))
   if err != nil {
      panic(err)
   }
   req.Header.Add("accept", "*/*")
   req.Header.Add("accept-encoding", "identity")
   req.Header.Add("accept-language", "en-US,en;q=0.5")
   req.Header.Add("cache-control", "max-age=0")
   req.Header.Add("content-type", "application/json")
   req.Header.Add("origin", "https://myaccount.unext.jp")
   req.Header.Add("priority", "u=4")
   req.Header.Add("sec-fetch-dest", "empty")
   req.Header.Add("sec-fetch-mode", "cors")
   req.Header.Add("sec-fetch-site", "same-origin")
   req.AddCookie(&http.Cookie{Name: "_st", Value: "05d2e4e77e84ad9f5d51fd559f9ac4dbbc728edf233ead075e8cf442ebeeed54a%3A2%3A%7Bi%3A0%3Bs%3A3%3A%22_st%22%3Bi%3A1%3Bs%3A36%3A%22ee9f4211-13c3-4c0c-94fd-83fbebc867cc%22%3B%7D"})
   req.AddCookie(&http.Cookie{Name: "__td_signed", Value: "true"})
   req.Header.Add("te", "trailers")
   req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   resp, err := client.Do(req)
   if err != nil {
      panic(err)
   }
   if err := resp.Write(os.Stdout); err != nil {
      panic(err)
   }
}
