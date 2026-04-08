package hboMax

import (
   "slices"
)

// validVideoTypes acts as a set to hold the video types we want to keep
var validVideoTypes = []string{
   "EPISODE",
   "MOVIE",
}

type Included struct {
   Attributes *struct {
      EpisodeNumber int
      Name          string
      SeasonNumber  int
      ShowType string
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

func FilterAndSort(values []*Included) []*Included {
   values = slices.DeleteFunc(values, func(i *Included) bool {
      if i.Attributes == nil {
         return true // Remove videos with nil attributes.
      }
      // We return 'true' to delete if the video's type is NOT in our slice.
      return !slices.Contains(validVideoTypes, i.Attributes.VideoType)
   })
   slices.SortFunc(values, func(a, b *Included) int {
      if a.Attributes == nil || b.Attributes == nil {
         return 0 // Consider them equal if attributes are missing.
      }
      return a.Attributes.EpisodeNumber - b.Attributes.EpisodeNumber
   })
   return values
}
