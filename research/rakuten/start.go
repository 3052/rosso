package rakuten

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type StartResponse struct {
   Data SessionData `json:"data"`
}

type SessionData struct {
   Profile SessionProfile `json:"profile"`
   Market  SessionMarket  `json:"market"`
}

type SessionProfile struct {
   Classification ProfileClassification `json:"classification"`
   AudioLanguage  ProfileLanguage       `json:"audio_language"`
}

type ProfileClassification struct {
   Id string `json:"id"`
}

type ProfileLanguage struct {
   Id string `json:"id"`
}

type SessionMarket struct {
   Code string `json:"code"`
}

type StartPayload struct {
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

func StartSession(marketCode string) (*SessionData, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/me/start",
   }

   query := url.Values{}
   query.Set("market_code", marketCode)
   endpoint.RawQuery = query.Encode()

   startPayload := StartPayload{
      DeviceIdentifier: "web",
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

   payloadData, err := json.Marshal(startPayload)
   if err != nil {
      return nil, err
   }

   resp, err := maya.Post(endpoint, nil, payloadData)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var startResponse StartResponse
   if err := json.NewDecoder(resp.Body).Decode(&startResponse); err != nil {
      return nil, err
   }

   return &startResponse.Data, nil
}
