// post_signin.go
package amazon

import (
   "fmt"
   "io/ioutil"
   "net/http"
   "net/url"
   "strings"
)

// PostSignin submits the password using the dynamic Action URL and tokens from the PostClaim response.
func PostSignin(client *http.Client, pageData *PageData, email, password string) error {
   if pageData == nil || pageData.ActionURL == "" {
      return fmt.Errorf("invalid page data or missing action URL")
   }

   data := url.Values{}
   // Populate dynamically extracted hidden fields
   for k, v := range pageData.HiddenParams {
      data.Set(k, v)
   }
   data.Set("email", email)
   data.Set("password", password)

   req, err := http.NewRequest("POST", pageData.ActionURL, strings.NewReader(data.Encode()))
   if err != nil {
      return fmt.Errorf("error creating request: %w", err)
   }

   req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
   req.Header.Set("accept-language", "en-US,en;q=0.5")
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   req.Header.Set("origin", "https://www.amazon.com")
   req.Header.Set("upgrade-insecure-requests", "1")
   req.Header.Set("sec-fetch-dest", "document")
   req.Header.Set("sec-fetch-mode", "navigate")
   req.Header.Set("sec-fetch-site", "same-origin")
   req.Header.Set("sec-fetch-user", "?1")

   // Intercept the 302 Found redirect to prevent Go from blindly following it
   // and losing the context of the successful login.
   client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
      return http.ErrUseLastResponse
   }

   resp, err := client.Do(req)
   if err != nil {
      return fmt.Errorf("error executing request: %w", err)
   }
   defer resp.Body.Close()

   body, _ := ioutil.ReadAll(resp.Body)
   fmt.Printf("Login Request Complete. Status: %s\n", resp.Status)
   fmt.Printf("Response Body Length: %d\n", len(body))

   if resp.StatusCode == 302 {
      fmt.Printf("Success! Redirected to: %s\n", resp.Header.Get("Location"))
   } else if resp.StatusCode == 200 {
      fmt.Println("Returned 200 OK. Check if credentials were correct or if a captcha was triggered.")
   }

   return nil
}
