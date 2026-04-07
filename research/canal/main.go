package main

import (
   "log"
   "net/http"
   "time"
)

type App struct {
   Client        *http.Client
   DeviceSerial  string
   TVApiBaseURL  string
   ProvisionData string
   SsoToken      string
   BearerToken   string
   Ticket        string
}

func main() {
   app := &App{
      Client: &http.Client{
         Timeout: 15 * time.Second,
      },
      DeviceSerial: "w76d15b90-3215-11f1-87ca-01f0af932fb7", // Replace with a dynamic UUID if needed
   }
   app.TVApiBaseURL = "https://tvapi-hlm2.solocoo.tv"

   log.Println("3/6 Requesting Demo SSO Token...")
   if err := app.Demo(); err != nil {
      log.Fatalf("Demo failed: %v", err)
   }
   log.Println("5/6 Initializing Login (Fetching Ticket)...")
   if err := app.LoginInit(); err != nil {
      log.Fatalf("LoginInit failed: %v", err)
   }
   log.Println("6/6 Submitting Credentials...")
   // Feed your credentials here
   if err := app.LoginSubmit("27@riseup.net", "***REMOVED***"); err != nil {
      log.Fatalf("LoginSubmit failed: %v", err)
   }
   log.Println("Successfully logged in!")
}

// setCommonHeaders applies headers that are present across the requests
func setCommonHeaders(req *http.Request) {
   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   req.Header.Set("Accept", "application/json, text/plain, */*")
   req.Header.Set("Accept-Language", "en-US,en;q=0.5")
   req.Header.Set("Origin", "https://play.canalplus.cz")
   req.Header.Set("Referer", "https://play.canalplus.cz/")
}
