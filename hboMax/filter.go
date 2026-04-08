package hboMax

import "slices"

func FilterAndSort(values []*Included) []*Included {
   values = slices.DeleteFunc(values, func(i *Included) bool {
      if i.Attributes == nil {
         return true // Remove videos with nil attributes.
      }
      // Check if the current types are in our valid slices
      isValidVideo := slices.Contains(valid_types, i.Attributes.VideoType)
      isValidShow := slices.Contains(valid_types, i.Attributes.ShowType)
      return !isValidVideo && !isValidShow
   })
   slices.SortFunc(values, func(a, b *Included) int {
      if a.Attributes == nil || b.Attributes == nil {
         return 0 // Consider them equal if attributes are missing.
      }
      return a.Attributes.EpisodeNumber - b.Attributes.EpisodeNumber
   })
   return values
}

var valid_types = []string{
   "EPISODE",
   "MOVIE",
}

type Included struct {
   Attributes *struct {
      EpisodeNumber int
      Name          string
      SeasonNumber  int
      ShowType      string
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
