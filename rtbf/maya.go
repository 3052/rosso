package rtbf

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
)

func (e *Entitlement) FetchWidevine(data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", "https://exposure.api.redbee.live", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.URL.Path = "/v2/license/customer/RTBF/businessunit/Auvio/widevine"
   req.URL.RawQuery = url.Values{
      "contentId":  {e.AssetId},
      "ls_session": {e.PlayToken},
   }.Encode()
   req.Header.Set("content-type", "application/x-protobuf")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      var value struct {
         Message string
      }
      err = json.NewDecoder(resp.Body).Decode(&value)
      if err != nil {
         return nil, err
      }
      return nil, errors.New(value.Message)
   }
   return io.ReadAll(resp.Body)
}
