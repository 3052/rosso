package hboMax

import (
   "fmt"
   "strings"
)

// String implements the fmt.Stringer interface to provide a clean visual output for the Entity.
func (e *Entity) String() string {
   data := &strings.Builder{}

   // 1. print episode number if material type is episode
   if e.Attributes.MaterialType == "EPISODE" {
      fmt.Fprintf(data, "Episode: %d\n", e.Attributes.EpisodeNumber)
   }

   // 2. print attributes name
   fmt.Fprintf(data, "Name: %s\n", e.Attributes.Name)

   // 3 & 4. print edit ID if type is video, otherwise print ID
   if e.Type == "video" {
      fmt.Fprintf(data, "Edit ID: %s\n", e.Relationships.Edit.Data.Id)
   } else {
      fmt.Fprintf(data, "ID: %s\n", e.Id)
   }

   // 5. print either show type or video type
   if e.Attributes.ShowType != "" {
      fmt.Fprintf(data, "Show Type: %s\n", e.Attributes.ShowType)
   } else if e.Attributes.VideoType != "" {
      fmt.Fprintf(data, "Video Type: %s\n", e.Attributes.VideoType)
   }

   return strings.TrimSpace(data.String())
}
