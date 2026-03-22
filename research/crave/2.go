package main

import (
   "io"
   "net/http"
   "net/url"
   "os"
   "strings"
)

func main() {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Add("Accept", "*/*")
   req.Header.Add("Accept-Encoding", "identity")
   req.Header.Add("Accept-Language", "en-US,en;q=0.5")
   req.Header.Add("Authorization", "Basic Y3JhdmUtd2ViOmRlZmF1bHQ=")
   req.Header.Add("Content-Length", "137")
   req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
   req.Header.Add("Origin", "https://www.crave.ca")
   req.Header.Add("Priority", "u=4")
   req.Header.Add("Referer", "https://www.crave.ca/")
   req.Header.Add("Sec-Fetch-Dest", "empty")
   req.Header.Add("Sec-Fetch-Mode", "cors")
   req.Header.Add("Sec-Fetch-Site", "cross-site")
   req.Header.Add("Te", "trailers")
   req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   req.Method = "POST"
   req.ProtoMajor = 1
   req.ProtoMinor = 1
   req.URL = &url.URL{}
   req.URL.Host = "account.bellmedia.ca"
   req.URL.Path = "/api/login/v2.2"
   req.URL.Scheme = "https"
   req.Body = io.NopCloser(strings.NewReader(data))
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   err = resp.Write(os.Stdout)
   if err != nil {
      panic(err)
   }
}

var data = url.Values{
   "grant_type":[]string{"magic_link_token"},
   "magic_link_token":[]string{"MgZ_TtZxjd_pPLXjbvAqo_SQpTJOsz_tKPd8Mr7kr4Cnnjqm5eS125kLuecrMKbIxovx-qdUSbBn_pQ9iqlahA=="},
}.Encode()
