package disney

import (
   _ "embed"
   "errors"
   "log"
   "net/http"
   "net/url"
   "path"
   "strings"
)

// ZGlzbmV5JmJyb3dzZXImMS4wLjA
// disney&browser&1.0.0
const client_api_key = "ZGlzbmV5JmJyb3dzZXImMS4wLjA.Cu56AgSfBTDag5NiRA81oLHkDZfu5L3CKadnefEAY84"

//go:embed authenticateWithOtp.gql
var mutation_authenticate_with_otp string

//go:embed login.gql
var mutation_login string

//go:embed loginWithActionGrant.gql
var mutation_login_with_action_grant string

//go:embed refreshToken.gql
var mutation_refresh_token string

//go:embed registerDevice.gql
var mutation_register_device string

//go:embed requestOtp.gql
var mutation_request_otp string

//go:embed switchProfile.gql
var mutation_switch_profile string

// https://disneyplus.com/browse/entity-7df81cf5-6be5-4e05-9ff6-da33baf0b94d
// https://disneyplus.com/cs-cz/browse/entity-7df81cf5-6be5-4e05-9ff6-da33baf0b94d
// https://disneyplus.com/play/7df81cf5-6be5-4e05-9ff6-da33baf0b94d
func GetEntityId(rawUrl string) (string, error) {
   parsed, err := url.Parse(rawUrl)
   if err != nil {
      return "", err
   }
   base := path.Base(parsed.Path)
   if !strings.HasPrefix(base, "entity-") {
      return "", errors.New("entity value missing from URL")
   }
   return base, nil
}

func do(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
}

type AuthenticateWithOtp struct {
   ActionGrant string
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

type Login struct {
   Account struct {
      Profiles []Profile
   }
}

type LoginWithActionGrant struct {
   Account struct {
      Profiles []Profile
   }
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
         Message string
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

type Profile struct {
   Name string
   Id   string
}

func (p *Profile) String() string {
   var data strings.Builder
   data.WriteString("name: ")
   data.WriteString(p.Name)
   data.WriteString("\nid: ")
   data.WriteString(p.Id)
   return data.String()
}

type RequestOtp struct {
   Accepted bool
}

func (r *RequestOtp) String() string {
   if r.Accepted {
      return "accepted = true"
   }
   return "accepted = false"
}

type Season struct {
   Items []struct {
      Actions []struct {
         InternalTitle string
      }
   }
}

func (s Season) String() string {
   var (
      data strings.Builder
      line bool
   )
   for _, item := range s.Items {
      for _, action := range item.Actions {
         if line {
            data.WriteByte('\n')
         } else {
            line = true
         }
         data.WriteString(action.InternalTitle)
      }
   }
   return data.String()
}

// disney.go
