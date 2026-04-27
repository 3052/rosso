package rakuten

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type StartResponse struct {
   Profile Profile `json:"profile"`
   Market  Market  `json:"market"`
}

type Profile struct {
   Classification Classification `json:"classification"`
   AudioLanguage  Language       `json:"audio_language"`
}

type Classification struct {
   NumericalId int `json:"numerical_id"`
}

type Language struct {
   Id string `json:"id"`
}

type Market struct {
   Code string `json:"code"`
}

type StartRequest struct {
   DeviceIdentifier string         `json:"device_identifier"`
   DeviceMetadata   DeviceMetadata `json:"device_metadata"`
}

type DeviceMetadata struct {
   AppVersion   string `json:"app_version"`
   Brand        string `json:"brand"`
   Model        string `json:"model"`
   Os           string `json:"os"`
   SerialNumber string `json:"serial_number"`
   Uid          string `json:"uid"`
   Year         int    `json:"year"`
}

func FetchStart(marketCode string, deviceIdentifier string) (*StartResponse, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/me/start",
   }

   query := url.Values{}
   query.Set("market_code", marketCode)
   target.RawQuery = query.Encode()

   payload := StartRequest{
      DeviceIdentifier: deviceIdentifier,
      DeviceMetadata: DeviceMetadata{
         AppVersion:   "app_version",
         Brand:        "brand",
         Model:        "model",
         Os:           "os",
         SerialNumber: "serial_number",
         Uid:          "uid",
         Year:         0,
      },
   }

   reqBytes, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "content-type": "application/json",
   }

   resp, err := maya.Post(target, headers, reqBytes)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var respWrapper struct {
      Data StartResponse `json:"data"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&respWrapper); err != nil {
      return nil, err
   }

   return &respWrapper.Data, nil
}
