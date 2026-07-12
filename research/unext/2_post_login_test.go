// 2_post_login_test.go
package unext

import (
   "os"
   "testing"
)

func TestPostLogin(t *testing.T) {
   client, err := newProxyClient()
   if err != nil {
      t.Fatalf("failed to create proxy client: %v", err)
   }

   csrfToken, err := GetLoginPage(client)
   if err != nil {
      t.Fatalf("failed to get CSRF token: %v", err)
   }

   loginID := os.Getenv("UNEXT_LOGIN_ID")
   password := os.Getenv("UNEXT_PASSWORD")
   recaptchaResponse := os.Getenv("UNEXT_RECAPTCHA_RESPONSE")

   if loginID == "" {
      t.Skip("UNEXT_LOGIN_ID not set")
   }
   if password == "" {
      t.Skip("UNEXT_PASSWORD not set")
   }
   if recaptchaResponse == "" {
      t.Skip("UNEXT_RECAPTCHA_RESPONSE not set")
   }

   err = PostLogin(client, csrfToken, recaptchaResponse, loginID, password)
   if err != nil {
      t.Fatalf("PostLogin failed: %v", err)
   }

   t.Log("login successful")
}
