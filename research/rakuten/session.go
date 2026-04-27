package rakuten

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type UserSessionLanguage struct {
   Id string `json:"id"`
}

type UserSessionClassification struct {
   NumericalId int `json:"numerical_id"`
}

type UserSessionProfile struct {
   Classification    UserSessionClassification `json:"classification"`
   AudioLanguage     UserSessionLanguage       `json:"audio_language"`
   SubtitlesLanguage UserSessionLanguage       `json:"subtitles_language"`
}

type UserSessionMarket struct {
   Code string `json:"code"`
}

type UserSession struct {
   Profile UserSessionProfile `json:"profile"`
   Market  UserSessionMarket  `json:"market"`
}

type UserSessionDeviceMetadata struct {
   AppVersion   string `json:"app_version"`
   Brand        string `json:"brand"`
   Model        string `json:"model"`
   Os           string `json:"os"`
   SerialNumber string `json:"serial_number"`
   Uid          string `json:"uid"`
   Year         int    `json:"year"`
}

type UserSessionPayload struct {
   DeviceIdentifier string                    `json:"device_identifier"`
   DeviceMetadata   UserSessionDeviceMetadata `json:"device_metadata"`
}

func StartUserSession(marketCode string) (*UserSession, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/me/start",
   }
   values := url.Values{}
   values.Set("market_code", marketCode)
   endpoint.RawQuery = values.Encode()

   payload := UserSessionPayload{
      DeviceIdentifier: "web",
      DeviceMetadata: UserSessionDeviceMetadata{
         AppVersion:   "app_version",
         Brand:        "brand",
         Model:        "model",
         Os:           "os",
         SerialNumber: "serial_number",
         Uid:          "uid",
         Year:         0,
      },
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   resp, err := maya.Post(endpoint, nil, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var wrapper struct {
      Data UserSession `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
      return nil, err
   }

   return &wrapper.Data, nil
}
