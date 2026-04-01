package crave

import (
   "encoding/json"
   "fmt"
   "log"
   "net/http"
   "net/url"
   "os"
   "os/exec"
   "slices"
   "strings"
   "testing"
)

const public_url = "https://crave.ca/movie/goldeneye-38860"

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

func TestContent(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(cache + "/rosso/crave-final.json")
   if err != nil {
      t.Fatal(err)
   }
   var final_tokens Account
   err = json.Unmarshal(data, &final_tokens)
   if err != nil {
      t.Fatal(err)
   }
   log.SetFlags(log.Ltime)
   username, err := run("credential", "-h=api.nordvpn.com", "-k=username")
   if err != nil {
      t.Fatal(err)
   }
   password, err := run("credential", "-h=api.nordvpn.com", "-k=password")
   if err != nil {
      t.Fatal(err)
   }
   proxy := url.URL{
      Scheme: "https",
      User:   url.UserPassword(username, password),
      Host:   "ca1103.nordvpn.com:89",
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
   //////////////////////////////////////////////////
   // Magic happens here
   media_id, err := ParseMediaId(public_url)
   if err != nil {
      t.Fatal(err)
   }
   media_data, err := FetchMedia(media_id)
   if err != nil {
      t.Fatal(err)
   }
   content, err := media_data.FetchContentPackage()
   if err != nil {
      t.Fatal(err)
   }
   manifest_data, err := content.FetchManifest(
      media_data.FirstContent.Id, final_tokens.AccessToken,
   )
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println("DASH Manifest URL:", manifest_data)
}
