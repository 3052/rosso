package plex

import (
   "encoding/json"
   "log"
   "net/http"
   "net/url"
)

// do logs and executes the request.
func do(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
}

type MatchItem struct {
   Key       string `json:"key"`
   RatingKey string `json:"ratingKey"`
   Title     string `json:"title"`
   Type      string `json:"type"`
}

type Media struct {
   Id       string    `json:"id"`
   Protocol string    `json:"protocol"`
   Part     []VodPart `json:"Part"`
}

type MetadataItem struct {
   Title string  `json:"title"`
   Media []Media `json:"Media"`
}

type User struct {
   Id        int    `json:"id"`
   Uuid      string `json:"uuid"`
   AuthToken string `json:"authToken"`
}

type VodMetadata struct {
   Metadata []MetadataItem `json:"Metadata"`
}

func GetVodMetadata(match *MatchItem, userData *User) (*VodMetadata, error) {
   log.Print("GetVodMetadata")
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "vod.provider.plex.tv",
      Path:   match.Key,
   }

   req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("accept", "application/json")
   req.Header.Set("x-plex-token", userData.AuthToken)

   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var result struct {
      MediaContainer VodMetadata `json:"MediaContainer"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   return &result.MediaContainer, nil
}

type VodPart struct {
   Id      string `json:"id"`
   Key     string `json:"key"`
   License string `json:"license"`
}
