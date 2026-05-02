// FILE: crave/media.go
package crave

import (
   "41.neocities.org/maya"
   "encoding/json"
   "net/url"
   "strconv"
)

func GetMedia(showId int) (*Media, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "rte-api.bellmedia.ca",
      Path:   "/graphql",
   }

   headers := map[string]string{
      "authorization": "Bearer eyAicGxhdGZvcm0iOiAicGxhdGZvcm1fd2ViIiB9",
   }

   bodyMap := map[string]interface{}{
      "query": get_showpage,
      "variables": map[string]interface{}{
         "ids": []string{strconv.Itoa(showId)},
         "sessionContext": map[string]interface{}{
            "userLanguage": "EN",
            "userMaturity": "ADULT",
         },
      },
   }

   body, err := json.Marshal(bodyMap)
   if err != nil {
      return nil, err
   }

   resp, err := maya.Post(endpoint, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var result struct {
      Data struct {
         Medias []Media `json:"medias"`
      } `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   return &result.Data.Medias[0], nil
}

type Media struct {
   FirstContent FirstContent `json:"firstContent"`
   Id           int          `json:"id,string"`
}

type FirstContent struct {
   Id int `json:"id,string"`
}
