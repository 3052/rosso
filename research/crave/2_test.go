package crave

import (
   "encoding/json"
   "os"
   "testing"
)

func TestTwo(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(cache + "/rosso/crave.json")
   if err != nil {
      t.Fatal(err)
   }
   var zero_data zero
   err = json.Unmarshal(data, &zero_data)
   if err != nil {
      t.Fatal(err)
   }
   magic_link_token, err := zero_data.magic_link_token()
   if err != nil {
      t.Fatal(err)
   }
   resp, err := two(magic_link_token)
   if err != nil {
      t.Fatal(err)
   }
   err = resp.Write(os.Stdout)
   if err != nil {
      t.Fatal(err)
   }
}
