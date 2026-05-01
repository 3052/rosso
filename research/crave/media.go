package crave

import (
   "41.neocities.org/maya"
   "encoding/json"
   "net/url"
)

func GetMedia(showId string) ([]Media, error) {
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
         "ids": []string{showId},
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

   var wrapper struct {
      Data struct {
         Medias []Media `json:"medias"`
      } `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
      return nil, err
   }

   return wrapper.Data.Medias, nil
}

type Media struct {
   FirstContent FirstContent `json:"firstContent"`
   Id           string       `json:"id"`
}

type FirstContent struct {
   Id string `json:"id"`
}
