package rakuten

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type Session struct {
   Market Market `json:"market"`
}

type Market struct {
   Code string `json:"code"`
}

func StartSession(marketCode string) (*Session, error) {
   link := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/me/start",
   }
   values := url.Values{}
   values.Set("market_code", marketCode)
   link.RawQuery = values.Encode()

   payload := map[string]any{
      "device_identifier": "web",
      "device_metadata": map[string]any{
         "app_version":   "app_version",
         "brand":         "brand",
         "model":         "model",
         "os":            "os",
         "serial_number": "serial_number",
         "uid":           "uid",
         "year":          0,
      },
   }
   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "content-type": "application/json",
   }

   resp, err := maya.Post(link, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var respWrapper struct {
      Data Session `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&respWrapper); err != nil {
      return nil, err
   }
   return &respWrapper.Data, nil
}
