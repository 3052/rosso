package crave

import (
   "encoding/json"
   "fmt"
   "os"
   "testing"
)

func TestMagicLink(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(cache + "/rosso/crave.json")
   if err != nil {
      t.Fatal(err)
   }
   token := &TokenResponse{}
   err = json.Unmarshal(data, token)
   if err != nil {
      t.Fatal(err)
   }
   client := NewClient()
   magic_link_token, err := client.GenerateMagicLink(token.AccessToken)
   if err != nil {
      t.Fatal(err)
   }
   token, err = client.MagicLinkLogin(magic_link_token)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", token)
}
