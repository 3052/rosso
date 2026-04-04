package crave

import (
   "bytes"
   "encoding/base64"
   "encoding/json"
   "errors"
   "net/http"
   "strconv"
   "strings"
)

// https://crave.ca/movie/anaconda-2025-59881
// https://crave.ca/play/anaconda-2025-3300246
//
// https://crave.ca/movie/goldeneye-38860
// https://crave.ca/play/goldeneye-938361
func ParseMediaId(urlData string) (int, error) {
   idx := strings.LastIndex(urlData, "-")
   if idx == -1 {
      return 0, strconv.ErrSyntax
   }
   return strconv.Atoi(urlData[idx+1:])
}

type Media struct {
   FirstContent struct {
      Id int `json:"id,string"`
   }
   Id int `json:"id,string"`
}

func FetchMedia(id int) (*Media, error) {
   body, err := marshal(map[string]any{
      "query": get_showpage,
      "variables": map[string]any{
         "sessionContext": map[string]string{
            "userLanguage": Language,
            "userMaturity": "ADULT",
         },
         "ids": []string{strconv.Itoa(id)},
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://rte-api.bellmedia.ca/graphql", bytes.NewBuffer(body),
   )
   if err != nil {
      return nil, err
   }
   bearer := base64.StdEncoding.EncodeToString(
      []byte(`{ "platform": "platform_web" }`),
   )
   req.Header.Set("Authorization", "Bearer "+bearer)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         Medias []Media
      }
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   if len(result.Data.Medias) == 0 || result.Data.Medias[0].FirstContent.Id == 0 {
      return nil, errors.New("content ID not found in GraphQL response")
   }
   return &result.Data.Medias[0], nil
}
