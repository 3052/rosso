// map_md.go
package amazon

import (
   "encoding/base64"
   "encoding/json"
)

// GenerateMapMD creates the Base64 encoded JSON string required for the map-md cookie.
// It mimics the metadata of the Prime Video Android app.
func GenerateMapMD() (string, error) {
   metadata := map[string]interface{}{
      "device_registration_data": map[string]string{
         "software_version": "130050002",
      },
      "app_identifier": map[string]interface{}{
         "package": "com.amazon.avod.thirdpartyclient",
         "SHA-256": []string{
            "2f19adeb284eb36f7f07786152b9a1d14b21653203ad0b04ebbf9c73ab6d7625",
         },
         "app_version":      "458000357",
         "app_version_name": "3.0.458.357",
         "app_sms_hash":     "e0kK4QFSWp0",
         "map_version":      "MAPAndroidLib-1.3.49030.0",
      },
      "app_info": map[string]int{
         "auto_pv":                   0,
         "auto_pv_with_smsretriever": 1,
         "smartlock_supported":       0,
         "permission_runtime_grant":  2,
      },
   }

   jsonBytes, err := json.Marshal(metadata)
   if err != nil {
      return "", err
   }

   return base64.StdEncoding.EncodeToString(jsonBytes), nil
}
