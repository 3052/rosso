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
   var login_data Login
   login_data.Token = string(data)
   // https://oldflix.com.br/browse/play/5d5d54a4d55dc050f8468513
   browse_data, err := login_data.FetchBrowse("5d5d54a4d55dc050f8468513")
   if err != nil {
      t.Fatal(err)
   }
   original, err := browse_data.GetOriginal()
   if err != nil {
      t.Fatal(err)
   }
   watch, err := browse_data.Watch(original.Id, login_data.Token)
   if err != nil {
      t.Fatal(err)
   }
   t.Log(watch)
}
