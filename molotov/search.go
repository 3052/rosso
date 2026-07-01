// search.go
package molotov

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
   "strings"
)

const x_forwarded_for = "178.132.106.134"

type SearchComponent struct {
   Endpoint struct {
      Payload struct {
         Payload struct {
            UiElement string `json:"ui.element"`
            AssetId   string `json:"asset.asset_id"`
         } `json:"payload"`
      } `json:"payload"`
   } `json:"endpoint"`
}

// Search searches for content (VODs or Live Streams) using a query string.
func Search(query string, signinResp *SigninResponse, userResp *UserResponse) ([]SearchComponent, error) {
   baseURL := "https://api-eu.fubo.tv/papi/v1/search/content"

   params := url.Values{}
   params.Add("category", "top_results")
   params.Add("fuzzy", "true")
   params.Add("query", query)

   fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

   req, err := http.NewRequest("GET", fullURL, nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("x-forwarded-for", x_forwarded_for)
   req.Header.Set("Authorization", "Bearer "+signinResp.AccessToken)
   req.Header.Set("x-user-id", userResp.ID)
   req.Header.Set("x-profile-id", userResp.Profiles[0].ID)
   req.Header.Set("x-device-id", DeviceID)
   req.Header.Set("x-application-id", "molotov")
   req.Header.Set("x-device-group", "desktop")
   req.Header.Set("x-device-type", "desktop")
   req.Header.Set("x-device-app", "web")
   req.Header.Set("x-client-version", "6.12.0")
   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")

   resp, err := doRequest(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var envelope struct {
      Error struct {
         Code      string `json:"code"`
         Message   string `json:"message"`
         LcMessage string `json:"lc_message"`
      } `json:"error"`
      Content struct {
         Sections []struct {
            Components []struct {
               Body struct {
                  Actions struct {
                     OnPlay []SearchComponent `json:"on_play"`
                  } `json:"actions"`
               } `json:"body"`
            } `json:"components"`
         } `json:"sections"`
      } `json:"content"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
      return nil, err
   }

   if envelope.Error.Code != "" {
      return nil, fmt.Errorf("code: %s, message: %s, lc_message: %s",
         envelope.Error.Code, envelope.Error.Message, envelope.Error.LcMessage)
   }

   if len(envelope.Content.Sections) == 0 {
      return nil, fmt.Errorf("no sections found in response")
   }

   // Extract and flatten the slice of SearchComponent items
   var results []SearchComponent
   for _, comp := range envelope.Content.Sections[0].Components {
      results = append(results, comp.Body.Actions.OnPlay...)
   }

   return results, nil
}

func (s *SearchComponent) String() string {
   var data strings.Builder
   data.WriteString("ui.element: ")
   data.WriteString(s.Endpoint.Payload.Payload.UiElement)
   data.WriteString("\nasset.asset_id: ")
   data.WriteString(s.Endpoint.Payload.Payload.AssetId)
   return data.String()
}
