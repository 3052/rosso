// request_cms_content.go
package tubi

import (
   "compress/gzip"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

// CMSResponse represents the relevant parts of the JSON response.
type CMSResponse struct {
   VideoResources []VideoResource `json:"video_resources"`
}

type VideoResource struct {
   Type          string        `json:"type"`
   LicenseServer LicenseServer `json:"license_server"`
}

type LicenseServer struct {
   URL string `json:"url"`
}

// GetCMSContent constructs the request URL using net/url, fetches the CMS data,
// and returns the dynamic Widevine License Server URL.
func GetCMSContent() (string, error) {
   // 1. Construct the URL cleanly using net/url
   u, err := url.Parse("https://uapi.adrise.tv/cms/content")
   if err != nil {
      return "", fmt.Errorf("failed to parse base URL: %w", err)
   }

   q := u.Query()
   q.Set("content_id", "610572")
   q.Set("deviceId", "!")
   q.Add("limit_resolutions[]", "h264_1080p")
   q.Add("limit_resolutions[]", "h265_1080p")
   q.Set("platform", "web")
   q.Add("video_resources[]", "dash")
   q.Add("video_resources[]", "dash_widevine")
   u.RawQuery = q.Encode()

   // 2. Create the HTTP request
   req, err := http.NewRequest("GET", u.String(), nil)
   if err != nil {
      return "", fmt.Errorf("failed to create request: %w", err)
   }

   req.Header.Set("accept-encoding", "gzip")
   req.Header.Set("user-agent", "Go-http-client/2.0")

   // 3. Execute the request
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return "", fmt.Errorf("request failed: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   // 4. Handle gzip decompression
   var reader io.Reader = resp.Body
   if resp.Header.Get("content-encoding") == "gzip" {
      gzReader, err := gzip.NewReader(resp.Body)
      if err != nil {
         return "", fmt.Errorf("failed to create gzip reader: %w", err)
      }
      defer gzReader.Close()
      reader = gzReader
   }

   // 5. Parse the JSON response
   var cmsResp CMSResponse
   if err := json.NewDecoder(reader).Decode(&cmsResp); err != nil {
      return "", fmt.Errorf("failed to decode JSON response: %w", err)
   }

   // 6. Extract the Widevine license server URL
   for _, resource := range cmsResp.VideoResources {
      if resource.Type == "dash_widevine" && resource.LicenseServer.URL != "" {
         return resource.LicenseServer.URL, nil
      }
   }

   return "", fmt.Errorf("widevine license URL not found in CMS response")
}
