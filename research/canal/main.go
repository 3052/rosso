package main

import (
   "log"
   "net/http"
   "net/http/cookiejar"
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
   jar, _ := cookiejar.New(nil)

   app := &App{
      Client: &http.Client{
         Timeout: 15 * time.Second,
         Jar:     jar,
      },
      DeviceSerial: "w76d15b90-3215-11f1-87ca-01f0af932fb7", // Replace with a randomly generated UUID if desired
   }

   log.Println("1/6 Fetching Config...")
   if err := app.GetConfig(); err != nil {
      log.Fatalf("Config failed: %v", err)
   }

   log.Println("2/6 Provisioning Device...")
   if err := app.Provision(); err != nil {
      log.Fatalf("Provision failed: %v", err)
   }

   log.Println("3/6 Requesting Demo SSO Token...")
   if err := app.Demo(); err != nil {
      log.Fatalf("Demo failed: %v", err)
   }

   log.Println("4/6 Starting Session...")
   if err := app.Session(); err != nil {
      log.Fatalf("Session failed: %v", err)
   }

   log.Println("5/6 Initializing Login (Fetching Ticket)...")
   if err := app.LoginInit(); err != nil {
      log.Fatalf("LoginInit failed: %v", err)
   }

   log.Println("6/6 Submitting Credentials...")
   // Put the target email and password here
   if err := app.LoginSubmit("27@riseup.net", "***REMOVED***"); err != nil {
      log.Fatalf("LoginSubmit failed: %v", err)
   }

   log.Println("Successfully completed the login flow and retrieved the final SSO Token!")
}

// setCommonHeaders applies headers that are present across most requests
func setCommonHeaders(req *http.Request) {
   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   req.Header.Set("Accept", "application/json, text/plain, */*")
   req.Header.Set("Accept-Language", "en-US,en;q=0.5")
   req.Header.Set("Origin", "https://play.canalplus.cz")
   req.Header.Set("Referer", "https://play.canalplus.cz/")
}
