package plex

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type User struct {
   Id        int    `json:"id"`
   Uuid      string `json:"uuid"`
   AuthToken string `json:"authToken"`
}

func CreateAnonymousUser() (*User, error) {
   targetUrl := &url.URL{
      Scheme: "https",
      Host:   "plex.tv",
      Path:   "/api/v2/users/anonymous",
   }

   headers := map[string]string{
      "x-plex-client-identifier": "!",
      "x-plex-product":           "Plex Mediaverse",
   }

   resp, err := maya.Post(targetUrl, headers, nil)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var user User
   if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
      return nil, err
   }

   return &user, nil
}
