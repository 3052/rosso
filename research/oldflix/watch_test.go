package oldflix

import (
   "os"
   "testing"
)

func TestWatch(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(cache + "/rosso/oldflix")
   if err != nil {
      t.Fatal(err)
   }
   client_data := NewClient()
   client_data.Token = string(data)
   // https://oldflix.com.br/browse/play/5d5d54a4d55dc050f8468513
   browse, err := client_data.BrowsePlay("5d5d54a4d55dc050f8468513")
   if err != nil {
      t.Fatal(err)
   }
   
   
   
   
   
   
}
