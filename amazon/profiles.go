package amazon

import (
   "encoding/json"
   "fmt"
   "net/http"
)

// ItemDetails contains metadata for a specific title, including the playback
// envelope
type ItemDetails struct {
   PlaybackEnvelope string `json:"playbackEnvelope"`
}

// GetItemDetails uses the actor access token to get metadata for a specific title.
// It explicitly passes UI schema flags to ensure the server returns the PlaybackEnvelope.
func GetItemDetails(actorToken *ActorToken, titleId, deviceTypeID string) (*ItemDetails, error) {
   url := HostATVExt + "/lrcedge/getDataByJavaTransform/v1/lr/detailsPage/detailsPageATF"
   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      return nil, err
   }
   query := req.URL.Query()
   query.Add("itemId", titleId)
   // Critical UI and Feature flags to force the V2/V3 BuyBox response with
   // PlaybackEnvelope
   query.Add("roles", "playback-envelope-supported")
   query.Add("presentationScheme", "android-tv-react")
   // Device parameters
   query.Add("deviceTypeID", deviceTypeID)
   query.Add("deviceID", DeviceID)
   req.Header.Set("Authorization", "Bearer "+actorToken.Token)
   req.URL.RawQuery = query.Encode()

   resp, err := doRequest(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   // Embed our new ItemDetails struct into the anonymous decoder struct
   var result struct {
      Resource struct {
         Actions []struct {
            Metadata struct {
               PlaybackExperienceMetadata ItemDetails
            }
         }
      }
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   for _, action := range result.Resource.Actions {
      details := action.Metadata.PlaybackExperienceMetadata
      if details.PlaybackEnvelope != "" {
         return &details, nil
      }
   }
   return nil, fmt.Errorf("playbackEnvelope not found in primaryActions for titleId: %s", titleId)
}

func (*ItemDetails) CachePath() string {
   return "rosso/amazon/ItemDetails"
}

// Profile represents an Amazon actor profile.
type Profile struct {
   ProfileID        string `json:"profileId"`
   IsDefaultProfile bool   `json:"isDefaultProfile"`
}

// GetPrimaryProfile uses the account access token to fetch available profiles and returns the primary profile.
func GetPrimaryProfile(tokens *TokenPair, deviceTypeID string) (*Profile, error) {
   url := HostATVExt + "/lrcedge/getDataByJavaTransform/v1/lr/profiles/profileSelection"
   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      return nil, err
   }
   query := req.URL.Query()
   query.Add("deviceTypeID", deviceTypeID)
   query.Add("deviceID", DeviceID)
   req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
   req.URL.RawQuery = query.Encode()

   resp, err := doRequest(req)
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
