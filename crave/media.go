package crave

import (
   "bytes"
   "encoding/base64"
   "encoding/json"
   "errors"
   "fmt"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

type Media struct {
   FirstContent struct {
      Id int `json:"id,string"`
   }
   Id int `json:"id,string"`
}

func ParseMedia(rawUrl string) (*Media, error) {
   parsedUrl, err := url.Parse(rawUrl)
   if err != nil {
      return nil, err
   }
   pathParts := strings.Split(parsedUrl.Path, "/")
   if len(pathParts) < 3 {
      return nil, fmt.Errorf("invalid url path structure")
   }
   urlType := pathParts[1]
   lastSegment := pathParts[len(pathParts)-1]
   dashIndex := strings.LastIndex(lastSegment, "-")
   if dashIndex == -1 {
      return nil, fmt.Errorf("id not found in url")
   }
   idString := lastSegment[dashIndex+1:]
   parsedId, err := strconv.Atoi(idString)
   if err != nil {
      return nil, fmt.Errorf("failed to parse id: %v", err)
   }
   m := &Media{}
   switch urlType {
   case "movie":
      m.Id = parsedId
   case "play":
      m.FirstContent.Id = parsedId
   default:
      return nil, fmt.Errorf("unknown media type: %s", urlType)
   }
   return m, nil
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
