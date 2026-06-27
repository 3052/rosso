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
func GetPrimaryProfile(accountAccessToken string) (*Profile, error) {
   url := HostATVExt + "/lrcedge/getDataByJavaTransform/v1/lr/profiles/profileSelection"
   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      return nil, err
   }
   query := req.URL.Query()
   query.Add("deviceTypeID", DeviceTypeID)
   query.Add("deviceID", DeviceID)
   req.Header.Set("Authorization", "Bearer "+accountAccessToken)
   req.URL.RawQuery = query.Encode()

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   // Embed our new Profile struct alongside the error Message struct
   var result struct {
      Resource struct {
         Profiles []Profile `json:"profiles"`
      } `json:"resource"`
      Message *struct {
         Body *struct {
            Code    string `json:"code"`
            Message string `json:"message"`
         } `json:"body"`
      } `json:"message"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, fmt.Errorf("failed to decode response (status %d): %w", resp.StatusCode, err)
   }

   // 1. Check for the structured JSON API error
   if result.Message != nil && result.Message.Body != nil {
      return nil, fmt.Errorf("API error [%s]: %s", result.Message.Body.Code, result.Message.Body.Message)
   }

   // 2. Check for standard HTTP errors if no JSON error message was provided
   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   // 3. Extract and return the primary profile
   for _, profile := range result.Resource.Profiles {
      if profile.IsDefaultProfile {
         return &profile, nil
      }
   }

   return nil, fmt.Errorf("default profile not found")
}

func (*Profile) CachePath() string {
   return "rosso/amazon/Profile"
}
