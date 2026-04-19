package nbc

import (
   "bytes"
   "io"
   "net/http"
)

func FetchWidevine(data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", "https://drmproxy.digitalsvc.apps.nbcuni.com",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }

   req.URL.Path = "/drm-proxy/license/widevine"
   req.URL.RawQuery = build_query("widevine")
   req.Header.Set("Content-Type", "application/octet-stream")

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}
