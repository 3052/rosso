package rakuten

import (
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func CreateLicense(streaming *StreamingInfo, body []byte) ([]byte, error) {
   location := &url.URL{
      Scheme: "https",
      Host:   "prod-kami.wuaki.tv",
      Path:   "/v1/licensing/wvm/" + url.PathEscape(streaming.Id),
   }
   query := url.Values{}
   query.Set("uuid", streaming.Id)
   location.RawQuery = query.Encode()

   headers := map[string]string{
      "content-type": "application/x-protobuf",
   }

   resp, err := maya.Post(location, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   body, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   return body, nil
}
