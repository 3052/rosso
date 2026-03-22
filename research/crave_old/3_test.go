package crave

import (
   "encoding/json"
   "fmt"
   "os"
   "testing"
)

func TestThree(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(cache + "/rosso/crave.json")
   if err != nil {
      t.Fatal(err)
   }
   account_data := &account{}
   err = json.Unmarshal(data, account_data)
   if err != nil {
      t.Fatal(err)
   }
   magic_link_token, err := account_data.magic_link_token()
   if err != nil {
      t.Fatal(err)
   }
   account_data, err = two(magic_link_token)
   if err != nil {
      t.Fatal(err)
   }
   err = account_data.three()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(account_data)
}
