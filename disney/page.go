package disney

import (
   "encoding/json"
   "errors"
   "log"
   "net/http"
   "strings"
)

func do(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
}

type Error struct {
   Code        string // 2026-04-05
   Description string // 2026-04-05
}

func (e *Error) Error() string {
   var data strings.Builder
   data.WriteString("code: ")
   data.WriteString(e.Code)
   data.WriteString("\ndescription: ")
   data.WriteString(e.Description)
   return data.String()
}

type Page struct {
   Actions []struct {
      InternalTitle string // movie
   }
   Containers []struct {
      Seasons []struct { // series
         Visuals struct {
            Name string
         }
         Id string
      }
   }
   Visuals struct {
      Restriction struct {
         Message string // 2026-08-16
      }
   }
}

func (p *Page) String() string {
   var data strings.Builder
   if len(p.Containers[0].Seasons) >= 1 {
      var line bool
      for _, seasonItem := range p.Containers[0].Seasons {
         if line {
            data.WriteString("\n\n")
         } else {
            line = true
         }
         data.WriteString("name: ")
         data.WriteString(seasonItem.Visuals.Name)
         data.WriteString("\nid: ")
         data.WriteString(seasonItem.Id)
      }
   } else {
      data.WriteString(p.Actions[0].InternalTitle)
   }
   return data.String()
}

type Token struct {
   AccessTokenType string
   AccessToken     string
   RefreshToken    string
}

// request: Account
func (t *Token) FetchPage(entity string) (*Page, error) {
   if err := t.assert("Account"); err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "GET",
      "https://disney.api.edge.bamgrid.com/explore/v1.12/page/"+entity+"?limit=0",
      nil,
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+t.AccessToken)
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         Errors []Error // 2026-04-11
         Page   Page
      }
      Errors []Error // 2026-05-03
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   if len(result.Data.Errors) >= 1 {
      return nil, &result.Data.Errors[0]
   }
   if restriction := result.Data.Page.Visuals.Restriction.Message; restriction != "" {
      return nil, errors.New(restriction)
   }
   return &result.Data.Page, nil
}

func (t *Token) assert(expected string) error {
   if t.AccessTokenType != expected {
      return errors.New("expected token type " + expected)
   }
   return nil
}

// page.go
