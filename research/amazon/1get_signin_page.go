// get_signin_page.go
package amazon

import (
   "fmt"
   "io/ioutil"
   "net/http"
   "net/url"
)

// GetSigninPage fetches the initial sign-in page and extracts the hidden form tokens.
func GetSigninPage(client *http.Client) (*PageData, error) {
   u := &url.URL{
      Scheme: "https",
      Host:   "www.amazon.com",
      Path:   "/ap/signin",
   }
   q := u.Query()
   
   q.Set("openid.return_to", "https://www.amazon.com/gp/video/detail/B075RND57T?ref_=nav_custrec_signin")
   q.Set("openid.identity", "http://specs.openid.net/auth/2.0/identifier_select")
   q.Set("openid.assoc_handle", "usflex")
   q.Set("openid.mode", "checkid_setup")
   q.Set("openid.claimed_id", "http://specs.openid.net/auth/2.0/identifier_select")
   q.Set("openid.ns", "http://specs.openid.net/auth/2.0")
   u.RawQuery = q.Encode()
   req, err := http.NewRequest("GET", u.String(), nil)
   if err != nil {
      return nil, fmt.Errorf("error creating request: %w", err)
   }
   req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
   req.Header.Set("accept-language", "en-US,en;q=0.5")
   req.Header.Set("referer", "https://www.amazon.com/gp/video/detail/B075RND57T")
   req.Header.Set("upgrade-insecure-requests", "1")
   req.Header.Set("sec-fetch-dest", "document")
   req.Header.Set("sec-fetch-mode", "navigate")
   req.Header.Set("sec-fetch-site", "same-origin")

   resp, err := client.Do(req)
   if err != nil {
      return nil, fmt.Errorf("error executing request: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   bodyBytes, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      return nil, fmt.Errorf("error reading response body: %w", err)
   }

   return ExtractPageData(string(bodyBytes)), nil
}
