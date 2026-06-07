package amazon

import (
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

// VerifyOTP submits the OTP verification form.
// It returns the redirect URL (to be used in GetAuthorizationCode) and updated cookies.
func VerifyOTP(formValues url.Values, cookies []*http.Cookie) (string, []*http.Cookie, error) {
   verifyUrl := "https://www.amazon.com/ap/cvf/verify"

   req, err := http.NewRequest("POST", verifyUrl, strings.NewReader(formValues.Encode()))
   if err != nil {
      return "", nil, err
   }

   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 11; sdk_gphone_x86_64 Build/RSR1.240422.006; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/83.0.4103.106 Mobile Safari/537.36")
   req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
   req.Header.Set("X-Requested-With", "com.amazon.avod.thirdpartyclient")
   req.Header.Set("Origin", "https://www.amazon.com")

   for _, cookie := range cookies {
      req.AddCookie(cookie)
   }

   // Create a custom client that stops at the first redirect to capture the Location header
   client := &http.Client{
      CheckRedirect: func(req *http.Request, via []*http.Request) error {
         return http.ErrUseLastResponse
      },
   }

   resp, err := client.Do(req)
   if err != nil {
      return "", nil, err
   }
   defer resp.Body.Close()

   // Drain the body so the connection can be reused
   _, _ = io.Copy(io.Discard, resp.Body)

   if resp.StatusCode != http.StatusFound {
      return "", nil, fmt.Errorf("expected 302 redirect, got status code: %d", resp.StatusCode)
   }

   location := resp.Header.Get("Location")
   if location == "" {
      return "", nil, errors.New("location header not found in the 302 response")
   }

   return location, resp.Cookies(), nil
}
