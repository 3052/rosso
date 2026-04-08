package canal

import (
   "encoding/json"
   "errors"
   "net/http"
   "net/url"
   "strings"
)

func (s *Session) Search(query string) ([]Collection, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme:   "https",
         Host:     "tvapi-hlm2.solocoo.tv",
         Path:     "/v1/search",
         RawQuery: url.Values{"query": {query}}.Encode(),
      },
      Header: http.Header{},
   }
   req.Header.Set("authorization", "Bearer "+s.Token)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Collection []Collection
      Message    string
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   return result.Collection, nil
}

func (a *Asset) String() string {
   var data strings.Builder
   data.WriteString("title = ")
   data.WriteString(a.Title)
   data.WriteString("\ntype = ")
   data.WriteString(a.Type)
   data.WriteString("\nid = ")
   data.WriteString(a.Id)
   return data.String()
}

type Asset struct {
   Title string
   Type  string
   Id    string
}

type Collection struct {
   Assets []Asset
}
