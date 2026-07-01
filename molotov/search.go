// search.go
package molotov

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
   "strings"
)

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
func Search(query string, signinResp *SigninResponse) ([]SearchComponent, error) {
   baseURL := "https://api-eu.fubo.tv/papi/v1/search/content"
   params := url.Values{}
   params.Add("query", query)
   fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
   req, err := http.NewRequest("GET", fullURL, nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("x-forwarded-for", x_forwarded_for)
   req.Header.Set("Authorization", "Bearer "+signinResp.AccessToken)
   req.Header.Set("x-application-id", "molotov")
   req.Header.Set("x-device-app", "web")
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
      for _, action := range comp.Body.Actions.OnPlay {
         // Only append the item if it actually contains the Movie tracking data
         // This filters out the blank "navigation" objects
         if action.Endpoint.Payload.Payload.UiElement != "" && action.Endpoint.Payload.Payload.AssetId != "" {
            results = append(results, action)
         }
      }
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
