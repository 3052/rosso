package amazon

import "testing"

func TestVod(t *testing.T) {
   vod_data, err := create_vod()
   if err != nil {
      t.Fatal(err)
   }
   t.Logf("%+v", vod_data)
}
