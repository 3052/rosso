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
