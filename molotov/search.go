package molotov

import (
   "encoding/json"
   "fmt"
   "log"
   "net/http"
   "net/url"
)

const x_forwarded_for = "178.132.106.134"

// doRequest logs the method and URL, then performs the HTTP request.
func doRequest(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   client := &http.Client{}
   return client.Do(req)
}

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
   params.Set("query", query)
   fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
   req, err := http.NewRequest("GET", fullURL, nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+signinResp.AccessToken)
   req.Header.Set("x-application-id", "molotov")
   req.Header.Set("x-device-app", "web")
   req.Header.Set("x-forwarded-for", x_forwarded_for)
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
                     // OnPlay []SearchComponent `json:"on_play"`
                     OnClick []SearchComponent `json:"on_click"`
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

   for _, section := range envelope.Content.Sections {
      for _, comp := range section.Components {
         for _, action := range comp.Body.Actions.OnClick {
            // Only append the item if it actually contains the Movie tracking data
            // This filters out the blank "navigation" objects
            if action.Endpoint.Payload.Payload.UiElement != "" && action.Endpoint.Payload.Payload.AssetId != "" {
               results = append(results, action)
            }
         }
      }
   }

   return results, nil
}

type SigninResponse struct {
   AccessToken  string `json:"access_token"`
   RefreshToken string `json:"refresh_token"`
}

// search.go
