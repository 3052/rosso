package roku

import (
   "41.neocities.org/maya"
   "encoding/json"
   "io"
   "net/url"
   "strings"
)

type AccountActivation struct {
   Code string `json:"code"`
}

func CreateAccountActivation(token *AccountToken) (*AccountActivation, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "googletv.web.roku.com",
      Path:   "/api/v1/account/activation",
   }
   headers := map[string]string{
      "content-type":         "application/json",
      "user-agent":           "trc-googletv; production; 0",
      "x-roku-content-token": token.AuthToken,
   }

   reqBody, err := json.Marshal(map[string]string{
      "platform": "googletv",
   })
   if err != nil {
      return nil, err
   }

   resp, err := maya.Post(target, headers, reqBody)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var activation AccountActivation
   if err := json.NewDecoder(resp.Body).Decode(&activation); err != nil {
      return nil, err
   }
   return &activation, nil
}

func (*AccountActivation) CachePath() string {
   return "rosso/roku/AccountActivation"
}

func (a *AccountActivation) String() string {
   var data strings.Builder
   data.WriteString("1 Visit the URL\n")
   data.WriteString("\ttherokuchannel.com/link\n")
   data.WriteString("2 Enter the activation code\n")
   data.WriteByte('\t')
   data.WriteString(a.Code)
   return data.String()
}

type AccountToken struct {
   AuthToken  string `json:"authToken"`
   IsLoggedIn bool   `json:"isLoggedIn"`
   Ip         string `json:"ip"`
   Rida       string `json:"rida"`
}

// status can be nil
func GetAccountToken(status *ActivationStatus) (*AccountToken, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "googletv.web.roku.com",
      Path:   "/api/v1/account/token",
   }
   headers := map[string]string{
      "user-agent": "trc-googletv; production; 0",
   }
   if status != nil {
      headers["x-roku-content-token"] = status.Token
   }

   resp, err := maya.Get(target, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var token AccountToken
   if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
      return nil, err
   }
   return &token, nil
}

func (*AccountToken) CachePath() string {
   return "rosso/roku/AccountToken"
}

type ActivationStatus struct {
   Code      string    `json:"code"`
   Token     string    `json:"token"`
   CreatedAt int64     `json:"createdAt"`
   Profiles  []Profile `json:"profiles"`
   Platform  string    `json:"platform"`
   Status    string    `json:"status"`
}

func GetActivationStatus(token *AccountToken, activation *AccountActivation) (*ActivationStatus, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "googletv.web.roku.com",
      Path:   "/api/v1/account/activation/" + activation.Code,
   }
   headers := map[string]string{
      "user-agent":           "trc-googletv; production; 0",
      "x-roku-content-token": token.AuthToken,
   }

   resp, err := maya.Get(target, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var status ActivationStatus
   if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
      return nil, err
   }
   return &status, nil
}

func (*ActivationStatus) CachePath() string {
   return "rosso/roku/ActivationStatus"
}

type Drm struct {
   Widevine Widevine `json:"widevine"`
}

type Playback struct {
   Url         *Url   // MPD
   Drm         Drm    `json:"drm"`
   MediaFormat string `json:"mediaFormat"`
   TraceId     string `json:"traceId"`
}

func GetPlayback(token *AccountToken, rokuId string) (*Playback, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "googletv.web.roku.com",
      Path:   "/api/v3/playback",
   }
   headers := map[string]string{
      "content-type":         "application/json",
      "user-agent":           "trc-googletv; production; 0",
      "x-roku-content-token": token.AuthToken,
   }

   reqBody, err := json.Marshal(map[string]string{
      "mediaFormat": "DASH",
      "providerId":  "rokuavod",
      "rokuId":      rokuId,
   })
   if err != nil {
      return nil, err
   }

   resp, err := maya.Post(target, headers, reqBody)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var result Playback
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   return &result, nil
}

func (*Playback) CachePath() string {
   return "rosso/roku/Playback"
}

func (p *Playback) LicenseWidevine(challenge []byte) ([]byte, error) {
   headers := map[string]string{
      "content-type": "application/x-protobuf",
      "user-agent":   "Go-http-client/2.0",
   }

   resp, err := maya.Post(&p.Drm.Widevine.LicenseServer.Url, headers, challenge)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}

type Profile struct {
   Id      string `json:"id"`
   IsKids  bool   `json:"isKids"`
   IsOwner bool   `json:"isOwner"`
}

type Url struct {
   Url url.URL
}

func (u *Url) MarshalText() ([]byte, error) {
   return u.Url.MarshalBinary()
}

func (u *Url) UnmarshalText(text []byte) error {
   return u.Url.UnmarshalBinary(text)
}

type Widevine struct {
   LicenseServer *Url `json:"licenseServer"`
}
