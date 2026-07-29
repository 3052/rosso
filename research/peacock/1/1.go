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

const bodyData = `{"auth":{"personaId":"3c0ce3d6-75bd-48c6-b3a9-fabce4c2ff83"},"device":{"id":"d81aabd73e994093"}}`

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
      Path:   "/auth/tokens",
   }

   hash := md5.Sum([]byte(bodyData))
   contentMD5 := hex.EncodeToString(hash[:])

   req, err := http.NewRequest("PATCH", reqURL.String(), bytes.NewBufferString(bodyData))
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
   req.Header.Add("x-skyott-usertoken", "165-L6EX7Sx+KllCZ1f75HKenh+zrrA7UVVuCYGcDExKipH+JZLKe69zC8NJMA3DIbP/j45y+M4H0lMMQs2Wi9y40v5NSGEGLiYyBLz/v1R7xl4c7Ev+QYBl2q3viyn1HRrm3b6g0hGPVJLWm25uaETb178cRUuxzwa+STElXRYS7AiU/ILfG+zNTywx2A67zYOo5t47OInQz6kIn/5MHl4UsjlGmIOTnO/oop0SK3lRHq5IvnGEV3dNXhBCZREbfi1MTorCYg/wIeT6a7d+u29eCedJti3+3oIpdH8HqzXna34XHRZVHgKMUfFgo4Xt8gJJByMUYeIkNOKOuDjZi06LbdNWViZTCJ9j10qqTE9GadJnGPhmOWKxSqtFUwmh8QCd/3hwDGClT3UHYLiHD4bqiAR0rdoOuEgSZ1mk5vQyS5anP0yaNIImR16x4SpKoKonm496PimH0ypOPGlGP8+ITA==")
   resp, err := client.Do(req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()

   if err := resp.Write(os.Stdout); err != nil {
      panic(err)
   }
}
