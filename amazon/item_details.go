package amazon

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
   "strings"
)

// PlaybackExperienceMetadata contains the envelope and related data needed for playback requests.
type PlaybackExperienceMetadata struct {
   PlaybackEnvelope string `json:"playbackEnvelope"`
   ExpiryTime       int64  `json:"expiryTime"`
   CorrelationId    string `json:"correlationId"`
}

// Profile represents an Amazon actor profile.
type Profile struct {
   ProfileID        string `json:"profileId"`
   IsDefaultProfile bool   `json:"isDefaultProfile"`
}

// GetPrimaryProfile uses the account access token to fetch available profiles and returns the primary profile.
func GetPrimaryProfile(tokens *TokenPair, deviceTypeID string) (*Profile, error) {
   req, err := http.NewRequest(
      "GET",
      HostATVPS+"/lrcedge/getDataByJavaTransform/v1/lr/profiles/profileSelection",
      nil,
   )
   if err != nil {
      return nil, err
   }
   query := url.Values{}
   query.Set("deviceTypeID", deviceTypeID)
   query.Set("deviceID", DeviceID)
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

// Resource represents the "resource" object returned from the detailsPageATF endpoint.
type Resource struct {
   Actions []struct {
      Metadata struct {
         PlaybackExperienceMetadata PlaybackExperienceMetadata `json:"playbackExperienceMetadata"`
      } `json:"metadata"`
   } `json:"actions"`
   ApplyHdr             bool `json:"applyHdr"`
   ApplyUhd             bool `json:"applyUhd"`
   EntitlementMessaging struct {
      EntitlementMessageSlotDetail struct {
         Message string `json:"message"`
      } `json:"ENTITLEMENT_MESSAGE_SLOT_DETAIL"`
   } `json:"entitlementMessaging"`
}

// GetItemDetails uses the actor access token to get metadata for a specific title.
// It explicitly passes UI schema flags to ensure the server returns the PlaybackEnvelope.
func GetItemDetails(token *ActorToken, titleId, deviceTypeID string) (*Resource, error) {
   req, err := http.NewRequest(
      "GET",
      HostATVPS+"/lrcedge/getDataByJavaTransform/v1/lr/detailsPage/detailsPageATF",
      nil,
   )
   if err != nil {
      return nil, err
   }
   query := url.Values{}
   query.Set("itemId", titleId)
   query.Set("presentationScheme", "android-tv-react")
   // Device parameters
   query.Set("deviceTypeID", deviceTypeID)
   query.Set("deviceID", DeviceID)

   if token != nil {
      // Critical UI and Feature flags to force the V2/V3 BuyBox response with
      // PlaybackEnvelope
      query.Set("roles", "playback-envelope-supported")
      // you can get the envelope without this, but it will be trailer:
      // resource.secondaryActions[0].presentation.label = "Watch trailer"
      req.Header.Set("Authorization", "Bearer "+token.Token)
   } else {
      query.Add("firmware", "")
      query.Add("roles", "prime-offer-supported,svod-supported")
      query.Add("clientFeatures", "EnableBuyBoxV2")
   }

   req.URL.RawQuery = query.Encode()

   resp, err := doRequest(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   // Embed our new Resource struct into the anonymous decoder struct
   var result struct {
      Resource Resource `json:"resource"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   return &result.Resource, nil
}

// GetPlaybackExperienceMetadata searches the Actions array and returns the first valid PlaybackExperienceMetadata.
func (r *Resource) GetPlaybackExperienceMetadata() (*PlaybackExperienceMetadata, error) {
   for _, action := range r.Actions {
      pem := action.Metadata.PlaybackExperienceMetadata
      if pem.PlaybackEnvelope != "" {
         return &pem, nil
      }
   }
   return nil, fmt.Errorf("playbackExperienceMetadata not found in actions")
}

func (r *Resource) String() string {
   var data strings.Builder
   if r.ApplyHdr {
      data.WriteString("HDR: true")
   } else {
      data.WriteString("HDR: false")
   }
   data.WriteByte('\n')
   if r.ApplyUhd {
      data.WriteString("UHD: true")
   } else {
      data.WriteString("UHD: false")
   }

   data.WriteByte('\n')
   data.WriteString("Message: ")
   data.WriteString(r.EntitlementMessaging.EntitlementMessageSlotDetail.Message)

   return data.String()
}
