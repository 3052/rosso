package crave

import (
   "encoding/json"
   "log"
   "net/http"
   "net/url"
   "os"
   "os/exec"
   "slices"
   "strings"
   "testing"
)

func TestFour(t *testing.T) {
   log.SetFlags(log.Ltime)
   http.DefaultTransport = &http.Transport{
      DisableKeepAlives: true, // github.com/golang/go/issues/25793
      Proxy: func(req *http.Request) (*url.URL, error) {
         log.Println(req.Method, req.URL)
         return nil, nil
      },
   }
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
   magic_link_token, err := GenerateMagicLink(auth_tokens.AccessToken)
   if err != nil {
      t.Fatal(err)
   }
   sso_tokens, err := MagicLinkLogin(magic_link_token)
   if err != nil {
      t.Fatal(err)
   }
   // 4. profiles
   profiles, err := GetProfiles(auth_tokens.AccountID, sso_tokens.AccessToken)
   if err != nil {
      t.Fatal(err)
   }
   i := slices.IndexFunc(profiles, func(p *Profile) bool {
      return p.HasPin == false
   })
   final_tokens, err := ProfileLogin(sso_tokens.RefreshToken, profiles[i].ID, "")
   if err != nil {
      t.Fatal(err)
   }
   username, err := run("credential", "-h=api.nordvpn.com", "-k=username")
   if err != nil {
      t.Fatal(err)
   }
   password, err := run("credential", "-h=api.nordvpn.com")
   if err != nil {
      t.Fatal(err)
   }
   proxy := url.URL{
      Scheme: "https",
      User: url.UserPassword(username, password),
      Host: "ca1103.nordvpn.com:89",
   }
   http.DefaultTransport = &http.Transport{
      Proxy: func(req *http.Request) (*url.URL, error) {
         if req.Method == "" {
            req.Method = "GET"
         }
         log.Println(req.Method, req.URL)
         return &proxy, nil
      },
   }
   resp, err := final_tokens.four()
   if err != nil {
      t.Fatal(err)
   }
   err = resp.Write(os.Stdout)
   if err != nil {
      t.Fatal(err)
   }
}

func run(name string, arg ...string) (string, error) {
   var data strings.Builder
   command := exec.Command(name, arg...)
   command.Stdout = &data
   err := command.Run()
   if err != nil {
      return "", err
   }
   return data.String(), nil
}
