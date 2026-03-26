package mubi

import "testing"

func Test(t *testing.T) {
   t.Log(tests)
}

var tests = []struct {
   locations []string
   url       string
}{
   {
      url: "https://mubi.com/films/passages-2022",
      locations: []string{
         "Austria",
         "Belgium",
         "Brazil",
         "Canada",
         "Chile",
         "Colombia",
         "Germany",
         "Ireland",
         "Italy",
         "Mexico",
         "Netherlands",
         "Peru",
         "Turkey",
         "United Kingdom",
         "United States",
      },
   },
   {
      url: "https://mubi.com/films/perfect-days",
      locations: []string{
         "Austria",
         "Brazil",
         "Chile",
         "Colombia",
         "Germany",
         "Ireland",
         "Mexico",
         "Peru",
         "Turkey",
         "United Kingdom",
      },
   },
}
