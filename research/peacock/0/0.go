package main

import (
   "bytes"
   "crypto/md5"
   "crypto/tls"
   "encoding/hex"
   "net/http"
   "net/url"
   "os"
)

const bodyData = `{"auth":{"authScheme":"OAUTH","authIssuer":"NOWTV","provider":"NBCU","providerTerritory":"US","proposition":"NBCUOTT","authToken":"a2cb8bf281cbfb426ebe48e56662a41371d12200f044c11ac0657f242fb21ce3"},"device":{"type":"TV","platform":"ANDROIDTV","id":"d81aabd73e994093","drmDeviceId":"UNKNOWN"}}`

func main() {
   pemData, err := os.ReadFile("play.clients.peacocktv.com.pem")
   if err != nil {
      panic(err)
   }

   cert, err := tls.X509KeyPair(pemData, pemData)
   if err != nil {
      panic(err)
   }

   client := &http.Client{
      Transport: &http.Transport{
         TLSClientConfig: &tls.Config{
            Certificates: []tls.Certificate{cert},
            MinVersion:   tls.VersionTLS12,
         },
      },
   }

   reqURL := &url.URL{
      Scheme: "https",
      Host:   "play.clients.peacocktv.com",
      Path:   "/auth/throttled/tokens",
   }

   hash := md5.Sum([]byte(bodyData))
   contentMD5 := hex.EncodeToString(hash[:])

   req, err := http.NewRequest("POST", reqURL.String(), bytes.NewBufferString(bodyData))
   if err != nil {
      panic(err)
   }
   req.Header.Add("accept", "application/vnd.tokens.v1+json")
   req.Header.Add("accept-language", "en-US,en;q=0.9")
   req.Header.Add("content-type", "application/vnd.tokens.v1+json")
   req.Header.Add("origin", "https://tv.clients.peacocktv.com")
   req.Header.Add("referer", "https://tv.clients.peacocktv.com/")
   req.Header.Add("sec-fetch-dest", "empty")
   req.Header.Add("sec-fetch-mode", "cors")
   req.Header.Add("sec-fetch-site", "same-site")
   req.Header.Add("content-md5", contentMD5)
   req.Header.Add("x-requested-with", "com.peacocktv.peacockandroid")
   req.Header.Add("x-skyott-activeterritory", "US")
   req.Header.Add("x-skyott-broadcastregions", "INPATTERN_US_CENTRAL")
   req.Header.Add("x-skyott-device", "TV")
   req.Header.Add("x-skyott-language", "en-US")
   req.Header.Add("x-skyott-platform", "ANDROIDTV")
   req.Header.Add("x-skyott-proposition", "NBCUOTT")
   req.Header.Add("x-skyott-provider", "NBCU")
   req.Header.Add("x-skyott-territory", "US")
   resp, err := client.Do(req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()

   if err := resp.Write(os.Stdout); err != nil {
      panic(err)
   }
}
