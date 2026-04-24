package plex

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type AnonymousUser struct {
   Id        int    `json:"id"`
   Uuid      string `json:"uuid"`
   AuthToken string `json:"authToken"`
}

func CreateAnonymousUser() (*AnonymousUser, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "plex.tv",
      Path:   "/api/v2/users/anonymous",
   }

   headers := map[string]string{
      "X-Plex-Client-Identifier": "!",
      "X-Plex-Product":           "Plex Mediaverse",
   }

   resp, err := maya.Post(endpoint, headers, nil)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var anonymous AnonymousUser
   if err := json.NewDecoder(resp.Body).Decode(&anonymous); err != nil {
      return nil, err
   }

   return &anonymous, nil
}
