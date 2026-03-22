package crave

import (
   "encoding/json"
   "os"
   "testing"
)

func TestZero(t *testing.T) {
   username, err := run("credential", "-h=crave.ca", "-k=username")
   if err != nil {
      t.Fatal(err)
   }
   password, err := run("credential", "-h=crave.ca")
   if err != nil {
      t.Fatal(err)
   }
   zero_data, err := fetch_zero(username, password)
   if err != nil {
      t.Fatal(err)
   }
   data, err := json.Marshal(zero_data)
   if err != nil {
      t.Fatal(err)
   }
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(cache + "/rosso/crave.json", data, os.ModePerm)
   if err != nil {
      t.Fatal(err)
   }
}
