package hboMax

import (
   "fmt"
   "strings"
)

const Markets = "amer apac emea latam"

type Cookie struct {
   Name  string
   Value string
}

func (*Cookie) CachePath() string {
   return "rosso/hboMax/Cookie"
}

func (c *Cookie) String() string {
   return fmt.Sprintf("%v=%v", c.Name, c.Value)
}

type Initiate struct {
   LinkingCode string
   TargetUrl   string
}

func (i *Initiate) String() string {
   var data strings.Builder
   data.WriteString("target URL: ")
   data.WriteString(i.TargetUrl)
   data.WriteString("\nlinking code: ")
   data.WriteString(i.LinkingCode)
   return data.String()
}

type Login struct {
   Token string
}

func (*Login) CachePath() string {
   return "rosso/hboMax/Login"
}

type Playback struct {
   Drm struct {
      Schemes struct {
         PlayReady *Scheme
         Widevine  *Scheme
      }
   }
   Errors   APIErrors `json:"errors"`
   Fallback struct {
      Manifest struct {
         Url string // _fallback.mpd:1080p, .mpd:4K
      }
   }
   Manifest struct {
      Url string // 1080p
   }
}

type Scheme struct {
   LicenseUrl string
}

// String implements the fmt.Stringer interface to provide a clean visual
// output for the Entity
func (e *Entity) String() string {
   data := &strings.Builder{}
   if e.Attributes.MaterialType == "EPISODE" {
      fmt.Fprintf(data, "Episode: %d\n", e.Attributes.EpisodeNumber)
   }
   if e.Attributes.ShowType != "" {
      fmt.Fprintf(data, "Show Type: %s\n", e.Attributes.ShowType)
   } else if e.Attributes.VideoType != "" {
      fmt.Fprintf(data, "Video Type: %s\n", e.Attributes.VideoType)
   }
   fmt.Fprintf(data, "Name: %s\n", e.Attributes.Name)
   if e.Type == "video" {
      fmt.Fprintf(data, "Edit ID: %s\n", e.Relationships.Edit.Data.Id)
   } else {
      fmt.Fprintf(data, "ID: %s\n", e.Id)
   }
   return strings.TrimSpace(data.String())
}

// types.go
