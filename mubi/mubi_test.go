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
         "AT", "BE", "BR", "CA", "CL", "CO", "DE", "GB", "IE", "IT", "MX", "NL",
         "PE", "TR", "US",
      },
   },
   {
      url: "https://mubi.com/films/perfect-days",
      locations: []string{
         "AT", "BR", "CL", "CO", "DE", "GB", "IE", "MX", "PE", "TR",
      },
   },
}
