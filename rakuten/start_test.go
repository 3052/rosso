package rakuten

import "testing"

func TestStart(t *testing.T) {
   profile_data, err := FetchProfile("cz")
   if err != nil {
      t.Fatal(err)
   }
   t.Logf("%+v", profile_data)
}
