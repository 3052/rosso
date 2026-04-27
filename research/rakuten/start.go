package rakuten

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type MarketCode string
type ClassificationId int
type LanguageId string

type Classification struct {
   NumericalId ClassificationId `json:"numerical_id"`
}

type Language struct {
   Id LanguageId `json:"id"`
}

type Profile struct {
   Classification Classification `json:"classification"`
   AudioLanguage  Language       `json:"audio_language"`
}

type Session struct {
   Profile Profile `json:"profile"`
}

func StartSession(market MarketCode) (*Session, error) {
   endpoint := url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/me/start",
   }

   query := url.Values{}
   query.Set("market_code", string(market))
   endpoint.RawQuery = query.Encode()

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

   resp, err := maya.Post(&endpoint, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var wrapper struct {
      Data Session `json:"data"`
   }
   decoder := json.NewDecoder(resp.Body)
   if err := decoder.Decode(&wrapper); err != nil {
      return nil, err
   }

   return &wrapper.Data, nil
}
