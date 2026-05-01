package crave

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type FirstContent struct {
   Id string `json:"id"`
}

type Media struct {
   FirstContent FirstContent `json:"firstContent"`
   Id           string       `json:"id"`
}

type SessionContext struct {
   UserLanguage string `json:"userLanguage"`
   UserMaturity string `json:"userMaturity"`
}

type ShowpageVariables struct {
   Ids            []string       `json:"ids"`
   SessionContext SessionContext `json:"sessionContext"`
}

type ShowpageRequest struct {
   Query     string            `json:"query"`
   Variables ShowpageVariables `json:"variables"`
}

func GetShowpage(id string) ([]Media, error) {
   endpoint := url.URL{
      Scheme: "https",
      Host:   "rte-api.bellmedia.ca",
      Path:   "/graphql",
   }

   payload := ShowpageRequest{
      Query: `query GetShowpage($sessionContext: SessionContext!, $ids: [String!]!) {
   medias(sessionContext: $sessionContext, ids: $ids) {
      firstContent {
         id
      }
      id
   }
}


`,
      Variables: ShowpageVariables{
         Ids: []string{id},
         SessionContext: SessionContext{
            UserLanguage: "EN",
            UserMaturity: "ADULT",
         },
      },
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "authorization": "Bearer eyAicGxhdGZvcm0iOiAicGxhdGZvcm1fd2ViIiB9",
   }

   resp, err := maya.Post(&endpoint, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var parsed struct {
      Data struct {
         Medias []Media `json:"medias"`
      } `json:"data"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
      return nil, err
   }

   return parsed.Data.Medias, nil
}
