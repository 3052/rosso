package roku

import "testing"

var tests = []struct {
   browse string
   url    string
   plan   bool
}{
   {
      url:    "https://therokuchannel.roku.com/details/50946a7e29d4869f9f81b77d4bdb5d42",
      browse: "movies",
      plan:   true,
   },
   {
      url:    "https://therokuchannel.roku.com/details/597a64a4a25c5bf6af4a8c7053049a6f",
      browse: "movies",
      plan:   false,
   },
   {
      url:    "https://therokuchannel.roku.com/details/105c41ea75775968b670fbb26978ed76",
      browse: "series",
      plan:   false,
   },
}

func Test(t *testing.T) {
   t.Log(tests)
}
