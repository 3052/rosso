package disney

import (
   _ "embed"
   "errors"
   "net/url"
   "strings"
)

func (s *Stream) GetManifest() (*url.URL, error) {
   return url.Parse(s.Sources[0].Complete.Url)
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
         data.WriteString("name = ")
         data.WriteString(seasonItem.Visuals.Name)
         data.WriteString("\nid = ")
         data.WriteString(seasonItem.Id)
      }
   } else {
      data.WriteString(p.Actions[0].InternalTitle)
   }
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
         Message string
      }
   }
}

func (e *Error) Error() string {
   var data strings.Builder
   data.WriteString("code = ")
   data.WriteString(e.Code)
   data.WriteString("\ndescription = ")
   data.WriteString(e.Description)
   return data.String()
}

type Error struct {
   Code        string // 2026-04-05
   Description string // 2026-04-05
}

func (t *Token) String() string {
   var data strings.Builder
   data.WriteString("type = ")
   data.WriteString(t.AccessTokenType)
   data.WriteString("\naccess token = ")
   data.WriteString(t.AccessToken)
   if t.RefreshToken != "" {
      data.WriteString("\nrefresh token = ")
      data.WriteString(t.RefreshToken)
   }
   return data.String()
}

type Token struct {
   AccessTokenType string
   AccessToken     string
   RefreshToken    string
}

// ZGlzbmV5JmJyb3dzZXImMS4wLjA
// disney&browser&1.0.0
const client_api_key = "ZGlzbmV5JmJyb3dzZXImMS4wLjA.Cu56AgSfBTDag5NiRA81oLHkDZfu5L3CKadnefEAY84"

//go:embed authenticateWithOtp.gql
var mutation_authenticate_with_otp string

//go:embed loginWithActionGrant.gql
var mutation_login_with_action_grant string

//go:embed registerDevice.gql
var mutation_register_device string

//go:embed login.gql
var mutation_login string

//go:embed requestOtp.gql
var mutation_request_otp string

//go:embed refreshToken.gql
var mutation_refresh_token string

//go:embed switchProfile.gql
var mutation_switch_profile string

// https://disneyplus.com/browse/entity-7df81cf5-6be5-4e05-9ff6-da33baf0b94d
// https://disneyplus.com/cs-cz/browse/entity-7df81cf5-6be5-4e05-9ff6-da33baf0b94d
// https://disneyplus.com/play/7df81cf5-6be5-4e05-9ff6-da33baf0b94d
func ParseEntity(urlData string) (string, error) {
   if strings.Contains(urlData, "/play/") {
      return "", errors.New("URL is a 'play' and not a 'browse'")
   }
   // The unique marker for the ID we want is "/browse/entity-".
   const marker = "/browse/entity-"
   // strings.Cut splits the string at the first instance of the marker.
   // It returns the part before, the part after, and a boolean indicating if the marker was found.
   // We don't need the 'before' part, so we discard it with the blank identifier _.
   _, id, found := strings.Cut(urlData, marker)
   // If the marker was not found, or if the resulting ID string is empty, return an error.
   if !found || id == "" {
      return "", errors.New("failed to find a valid ID in the URL")
   }
   // The 'id' variable now holds the rest of the string after the marker.
   return id, nil
}

type AuthenticateWithOtp struct {
   ActionGrant string
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

func (p *Profile) String() string {
   var data strings.Builder
   data.WriteString("name = ")
   data.WriteString(p.Name)
   data.WriteString("\nid = ")
   data.WriteString(p.Id)
   return data.String()
}

type Profile struct {
   Name string
   Id   string
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

type Season struct {
   Items []struct {
      Actions []struct {
         InternalTitle string
      }
   }
}

type Stream struct {
   Sources []struct {
      Complete struct {
         Url string
      }
   }
}

func (t *Token) assert(expected string) error {
   if t.AccessTokenType != expected {
      return errors.New("expected token type " + expected)
   }
   return nil
}
