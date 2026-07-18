// step6_get_license.go
package unext

import (
   "bytes"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

// Step6GetLicense POSTs a Widevine license challenge to the U-NEXT license
// proxy and returns the raw license response bytes.
//
// challenge is the binary SignedMessage (protobuf) produced by a Widevine CDM.
// The play_token must match the one used to fetch the MPD.
func Step6GetLicense(client *http.Client, licenseURL *url.URL, playToken string, challenge []byte) ([]byte, error) {
   q := licenseURL.Query()
   q.Set("play_token", playToken)
   licenseURL.RawQuery = q.Encode()

   req, err := http.NewRequest("POST", licenseURL.String(), bytes.NewReader(challenge))
   if err != nil {
      return nil, fmt.Errorf("step6: creating request: %w", err)
   }

   req.Header.Set("content-type", "application/octet-stream")
   req.Header.Set("user-agent", "U-NEXT Phone App Android12 5.73.1 sdk_gphone64_x86_64")

   resp, err := client.Do(req)
   if err != nil {
      return nil, fmt.Errorf("step6: sending request: %w", err)
   }
   defer resp.Body.Close()

   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, fmt.Errorf("step6: reading response body: %w", err)
   }

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("step6: expected 200, got %d: %s", resp.StatusCode, string(body))
   }

   return body, nil
}

// WidevineLicenseURL searches the playlist for the first DASH movie profile
// with a WIDEVINE license URL and returns it (without query parameters).
// Returns an error if no such profile is found.
func (p *PlaylistUrl) WidevineLicenseURL() (*url.URL, error) {
   for _, ui := range p.UrlInfo {
      for _, mp := range ui.MovieProfile {
         if mp.Type != "DASH" {
            continue
         }
         for _, lu := range mp.LicenseUrlList {
            if lu.Type == "WIDEVINE" && lu.LicenseUrl != "" {
               u, err := url.Parse(lu.LicenseUrl)
               if err != nil {
                  return nil, fmt.Errorf("parsing license URL: %w", err)
               }
               return u, nil
            }
         }
      }
   }
   return nil, fmt.Errorf("no DASH movie profile with WIDEVINE license found")
}
