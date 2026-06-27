package amazon

import (
   "encoding/json"
   "fmt"
   "net/http"
)

// Profile represents an Amazon actor profile.
type Profile struct {
   ProfileID        string `json:"profileId"`
   IsDefaultProfile bool   `json:"isDefaultProfile"`
}

// GetPrimaryProfile uses the account access token to fetch available profiles and returns the primary profile.
func GetPrimaryProfile(accountAccessToken string, deviceID DeviceID) (*Profile, error) {
   url := HostATVExt + "/lrcedge/getDataByJavaTransform/v1/lr/profiles/profileSelection"
   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      return nil, err
   }
   query := req.URL.Query()
   query.Add("deviceTypeID", DeviceTypeID)
   query.Add("deviceID", string(deviceID))
   req.Header.Set("Authorization", "Bearer "+accountAccessToken)
   req.URL.RawQuery = query.Encode()
   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   // Embed our new Profile struct into the anonymous decoder struct
   var result struct {
      Resource struct {
         Profiles []Profile `json:"profiles"`
      } `json:"resource"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   for _, profile := range result.Resource.Profiles {
      if profile.IsDefaultProfile {
         return &profile, nil
      }
   }

   if len(result.Resource.Profiles) > 0 {
      return &result.Resource.Profiles[0], nil
   }

   return nil, fmt.Errorf("no profiles found")
}

func (*Profile) CachePath() string {
   return "rosso/amazon/Profile"
}
