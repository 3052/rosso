package crave

import (
   "encoding/json"
   "fmt"
   "os"
   "testing"
)

func TestProfiles(t *testing.T) {
   // 1. Client
   client := NewClient()
   // 2. authTokens
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(cache + "/rosso/crave.json")
   if err != nil {
      t.Fatal(err)
   }
   var auth_tokens TokenResponse
   err = json.Unmarshal(data, &auth_tokens)
   if err != nil {
      t.Fatal(err)
   }
   // 3. ssoTokens
   magic_link_token, err := client.GenerateMagicLink(auth_tokens.AccessToken)
   if err != nil {
      t.Fatal(err)
   }
   sso_tokens, err := client.MagicLinkLogin(magic_link_token)
   if err != nil {
      t.Fatal(err)
   }
   // 4. profiles
   profiles, err := client.GetProfiles(auth_tokens.AccountID, sso_tokens.AccessToken)
   if err != nil {
      t.Fatal(err)
   }
   for _, profile_data := range profiles {
      fmt.Printf("%+v\n", profile_data)
   }
}
