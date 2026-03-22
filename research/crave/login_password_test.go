package crave

import (
   "encoding/json"
   "os"
   "testing"
)

func TestLoginPassword(t *testing.T) {
   username, err := run("credential", "-h=crave.ca", "-k=username")
   if err != nil {
      t.Fatal(err)
   }
   password, err := run("credential", "-h=crave.ca")
   if err != nil {
      t.Fatal(err)
   }
   client := NewClient()
   auth_tokens, err := client.PasswordLogin(username, password, "")
   if err != nil {
      t.Fatal(err)
   }
   data, err := json.Marshal(auth_tokens)
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
