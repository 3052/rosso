package amazon

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
   "strings"
)

// CodePair represents the public and private codes used for device linking.
type CodePair struct {
   PublicCode  string `json:"public_code"`
   PrivateCode string `json:"private_code"`
}

// CreateCodePair requests a public and private code pair for device linking.
func CreateCodePair() (*CodePair, error) {
   payload := map[string]interface{}{
      "code_data": map[string]string{
         "domain":        "Device",
         "device_type":   DeviceTypeID,
         "device_serial": DeviceID,
      },
   }
   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://api.amazon.com/auth/create/codepair",
      bytes.NewBuffer(body),
   )
   if err != nil {
      return nil, err
   }
   client := &http.Client{}
   resp, err := client.Do(req)
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
   data.WriteString("Please navigate to https://primevideo.com/ontv\n")
   data.WriteString("or https://amazon.com/code\n")
   data.WriteString("Enter the following code: ")
   data.WriteString(c.PublicCode)
   return data.String()
}
