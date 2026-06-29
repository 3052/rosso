package amazon

import (
   "bytes"
   "encoding/json"
   "errors"
   "fmt"
   "log"
   "net/http"
   "strings"
)

const ( // API Hosts
   HostAmazonAPI = "https://api.amazon.com"
   HostATVPS     = "https://atv-ps.amazon.com"
)

const DeviceID = "deviceID"

// the wrong DTID will fail the license request. if you change the DTID you
// need to relog. also if you get a failed license request try provision again.
// this might be UHD also
// > amazon-device -dtid A3GTP8TAF8V3YG
// manufacturer name: Hisense TV
// model number: HU43K3110FW
var Devices = []Device{
   {
      Manufacturer:  "Hisense",
      Model:         "HE55A7000EUWTS",
      SecurityLevel: 3000,
      DeviceTypeID:  "A3REWRVYBYPKUM",
   },
   {
      Manufacturer:  "Hisense",
      Model:         "HU50A6100UW",
      SecurityLevel: 3000,
      DeviceTypeID:  "AAJ692ZPT1X85",
   },
   {
      Manufacturer:  "Hisense",
      Model:         "HU32E5600FHWV",
      SecurityLevel: 3000,
      DeviceTypeID:  "A2RGJ95OVLR12U",
   },
   {
      Manufacturer:  "EXPRESS LUCK TECHNOLOGY LIMITED",
      Model:         "LE-*",
      SecurityLevel: 3000,
      DeviceTypeID:  "A3NM0WFSU3DLT5",
   },
}

// doRequest wraps the http.Client Do method to log every outgoing request.
func doRequest(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   client := &http.Client{}
   return client.Do(req)
}

func marshal(value any) ([]byte, error) {
   return json.MarshalIndent(value, "", " ")
}

// ActorToken represents an actor-specific access token.
type ActorToken struct {
   Token string `json:"token"`
}

// GetActorToken exchanges the account refresh token and actor ID for an actor-specific access token.
func GetActorToken(tokens *TokenPair, profile *Profile, deviceTypeID string) (*ActorToken, error) {
   payload := map[string]any{
      "actor_id":             profile.ProfileID,
      "app_name":             "AIV",
      "requested_token_type": "actor_access_token",
      "source_token_type":    "refresh_token",
      "source_device_tokens": []any{
         map[string]any{
            "device_type": deviceTypeID,
            "account_refresh_token": map[string]string{
               "token": tokens.RefreshToken,
            },
         },
      },
   }
   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", HostAmazonAPI+"/auth/token", bytes.NewBuffer(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/json")

   resp, err := doRequest(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   // Embed our new ActorToken struct into the anonymous decoder struct
   var result struct {
      DeviceTokens []struct {
         ActorAccessToken ActorToken `json:"actor_access_token"`
      } `json:"device_tokens"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   if len(result.DeviceTokens) == 0 {
      return nil, fmt.Errorf("no device tokens returned")
   }

   token := result.DeviceTokens[0].ActorAccessToken
   return &token, nil
}

func (*ActorToken) CachePath() string {
   return "rosso/amazon/ActorToken"
}

// CodePair represents the public and private codes used for device linking.
type CodePair struct {
   PublicCode  string `json:"public_code"`
   PrivateCode string `json:"private_code"`
}

// CreateCodePair requests a public and private code pair for device linking.
func CreateCodePair(deviceTypeID string) (*CodePair, error) {
   if deviceTypeID == "" {
      return nil, errors.New("deviceTypeID cannot be empty")
   }

   payload := map[string]any{
      "code_data": map[string]string{
         "domain":        "Device",
         "device_type":   deviceTypeID,
         "device_serial": DeviceID,
      },
   }
   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", HostAmazonAPI+"/auth/create/codepair",
      bytes.NewBuffer(body),
   )
   if err != nil {
      return nil, err
   }

   resp, err := doRequest(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   // Decode directly into our new struct type
   var result CodePair
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   return &result, nil
}

func (*CodePair) CachePath() string {
   return "rosso/amazon/CodePair"
}

func (c *CodePair) String() string {
   var data strings.Builder
   data.WriteString("Please navigate to https://amazon.com/gp/video/ontv\n")
   data.WriteString("Enter the following code: ")
   data.WriteString(c.PublicCode)
   return data.String()
}

// Device represents the metadata for a supported hardware device.
type Device struct {
   Manufacturer  string
   Model         string
   SecurityLevel int
   DeviceTypeID  string
}

func (*PlaybackExperienceMetadata) CachePath() string {
   return "rosso/amazon/PlaybackExperienceMetadata"
}

func (*TokenPair) CachePath() string {
   return "rosso/amazon/TokenPair"
}

// Refresh exchanges the existing refresh token for a new access token
// using the /auth/token endpoint, mutating the TokenPair in-place.
func (t *TokenPair) Refresh() error {
   if t == nil || t.RefreshToken == "" {
      return fmt.Errorf("invalid token pair or missing refresh token")
   }

   payload := map[string]string{
      "app_name":             "AIV",
      "requested_token_type": "access_token",
      "source_token":         t.RefreshToken,
      "source_token_type":    "refresh_token",
   }
   body, err := json.Marshal(payload)
   if err != nil {
      return err
   }
   req, err := http.NewRequest(
      "POST", HostAmazonAPI+"/auth/token", bytes.NewBuffer(body),
   )
   if err != nil {
      return err
   }
   req.Header.Set("content-type", "application/json")

   resp, err := doRequest(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()

   // Decode into an anonymous struct handling the expected Python response keys
   var result struct {
      AccessToken string `json:"access_token"`
      TokenType   string `json:"token_type"`
      Error       string `json:"error"`
      ErrorDesc   string `json:"error_description"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return err
   }

   // Handle API errors as seen in the Python code
   if result.Error != "" {
      return fmt.Errorf("failed to refresh device token: %s [%s]", result.ErrorDesc, result.Error)
   }

   if result.TokenType != "bearer" {
      return fmt.Errorf("unexpected returned refreshed token type: %s", result.TokenType)
   }

   // Mutate the struct in-place with the new access token
   t.AccessToken = result.AccessToken

   return nil
}
