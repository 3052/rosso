package stan

import "testing"

var programs = []struct {
   quality string
   url     []string
}{
   {
      quality: "high",
      url:     []string{"https://play.stan.com.au/programs/331144"},
   },
   {
      quality: "ultra",
      url: []string{
         "https://play.stan.com.au/programs/6299871",
         "https://stan.com.au/watch/beast-2026",
      },
   },
}

func TestProgram(t *testing.T) {
   t.Log(programs)
}
