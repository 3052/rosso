package hboMax

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
   "strings"
)

func (l Login) Search(query string) ([]Included, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme: "https",
         Host:   "default.prd.api.hbomax.com",
         Path:   "/cms/routes/search/result",
         RawQuery: url.Values{
            "contentFilter[query]": {query},
            "include":              {"default"},
         }.Encode(),
      },
      Header: http.Header{},
   }
   req.Header.Set("authorization", "Bearer "+l.Data.Attributes.Token)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   var result struct {
      Included []Included
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return result.Included, nil
}

func (i *Included) String() string {
   data := &strings.Builder{}
   if i.Attributes.EpisodeNumber >= 1 {
      fmt.Fprintln(data, "episode number =", i.Attributes.EpisodeNumber)
   }
   fmt.Fprintln(data, "name =", i.Attributes.Name)
   if i.Attributes.SeasonNumber >= 1 {
      fmt.Fprintln(data, "season number =", i.Attributes.SeasonNumber)
   }
   if i.Attributes.VideoType != "" {
      fmt.Fprintln(data, "video type =", i.Attributes.VideoType)
   }
   fmt.Fprint(data, "id = ", i.Id)
   if i.Relationships.Edit != nil {
      fmt.Fprint(data, "\nedit id = ", i.Relationships.Edit.Data.Id)
   }
   return data.String()
}

type Included struct {
   Attributes *struct {
      EpisodeNumber int
      Name          string
      SeasonNumber  int
      VideoType     string
   }
   Id            string
   Relationships *struct {
      Edit *struct {
         Data struct {
            Id string
         }
      }
   }
}
