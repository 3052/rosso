// helpers_test.go
package unext

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
   "os/exec"
)

func newProxyClient() (*http.Client, error) {
   cmd := exec.Command("credential.exe", "-j", "api.nordvpn.com")
   output, err := cmd.Output()
   if err != nil {
      return nil, fmt.Errorf("failed to get proxy credentials: %w", err)
   }

   var creds []nordCredential
   if err := json.Unmarshal(output, &creds); err != nil {
      return nil, fmt.Errorf("failed to parse credentials: %w", err)
   }

   if len(creds) == 0 {
      return nil, fmt.Errorf("no credentials returned")
   }

   proxyURL, err := url.Parse("https://" + creds[0].Username + ":" + creds[0].Password + "@jp528.proxy.nordvpn.com:89")
   if err != nil {
      return nil, fmt.Errorf("failed to parse proxy URL: %w", err)
   }

   return &http.Client{
      Transport: &http.Transport{
         Proxy: http.ProxyURL(proxyURL),
      },
   }, nil
}

type nordCredential struct {
   Date     string `json:"date"`
   Host     string `json:"host"`
   Password string `json:"password"`
   Username string `json:"username"`
}
