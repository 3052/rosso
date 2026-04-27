package rakuten

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type Language struct {
   Abbr string `json:"abbr"`
}

type Market struct {
   Code                 string    `json:"code"`
   DefaultAudioLanguage *Language `json:"default_audio_language"`
}

type Classification struct {
   NumericalId int `json:"numerical_id"`
}

type Profile struct {
   Classification *Classification `json:"classification"`
}

type SessionResponse struct {
   Profile *Profile `json:"profile"`
   Market  *Market  `json:"market"`
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

type StartPayload struct {
   DeviceIdentifier string          `json:"device_identifier"`
   DeviceMetadata   *DeviceMetadata `json:"device_metadata"`
}

func CreateSession(marketCode string) (*SessionResponse, error) {
   query := make(url.Values)
   query.Set("market_code", marketCode)

   endpoint := &url.URL{
      Scheme:   "https",
      Host:     "gizmo.rakuten.tv",
      Path:     "/v3/me/start",
      RawQuery: query.Encode(),
   }

   payload := &StartPayload{
      DeviceIdentifier: "web",
      DeviceMetadata: &DeviceMetadata{
         AppVersion:   "app_version",
         Brand:        "brand",
         Model:        "model",
         Os:           "os",
         SerialNumber: "serial_number",
         Uid:          "uid",
         Year:         0,
      },
   }

   bodyBytes, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "content-type": "application/json",
   }

   resp, err := maya.Post(endpoint, headers, bodyBytes)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var wrapper struct {
      Data *SessionResponse `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
      return nil, err
   }

   return wrapper.Data, nil
}
