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

func TestUrl(t *testing.T) {
   t.Log("https://crave.ca/movie/goldeneye-38860")
}

func TestPasswordLogin(t *testing.T) {
   username, err := run("credential", "-h=crave.ca", "-k=username")
   if err != nil {
      t.Fatal(err)
   }
   password, err := run("credential", "-h=crave.ca", "-k=password")
   if err != nil {
      t.Fatal(err)
   }
   auth_tokens, err := Login(username, password)
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
   err = os.WriteFile(cache+"/rosso/crave.json", data, os.ModePerm)
   if err != nil {
      t.Fatal(err)
   }
}

func TestFinalTokens(t *testing.T) {
   log.SetFlags(log.Ltime)
   http.DefaultTransport = &http.Transport{
      DisableKeepAlives: true, // github.com/golang/go/issues/25793
      Proxy: func(req *http.Request) (*url.URL, error) {
         log.Println(req.Method, req.URL)
         return nil, nil
      },
   }
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(cache + "/rosso/crave.json")
   if err != nil {
      t.Fatal(err)
   }
   var auth_tokens Account
   err = json.Unmarshal(data, &auth_tokens)
   if err != nil {
      t.Fatal(err)
   }
   profiles, err := auth_tokens.FetchProfiles()
   if err != nil {
      t.Fatal(err)
   }
   i := slices.IndexFunc(profiles, func(p *Profile) bool {
      return p.HasPin == false
   })
   //////////////////////////////////////////////////////////////
   err = auth_tokens.Login(profiles[i].Id)
   if err != nil {
      t.Fatal(err)
   }
   data, err = json.Marshal(auth_tokens)
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(cache+"/rosso/crave-final.json", data, os.ModePerm)
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
