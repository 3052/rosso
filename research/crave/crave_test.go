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

func TestContent(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(cache + "/rosso/crave-final.json")
   if err != nil {
      t.Fatal(err)
   }
   var final_tokens TokenResponse
   err = json.Unmarshal(data, &final_tokens)
   if err != nil {
      t.Fatal(err)
   }
   log.SetFlags(log.Ltime)
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
   publicUrl := "https://www.crave.ca/en/movie/goldeneye-38860"
   // Magic happens here
   mediaId, err := extractMediaId(publicUrl)
   if err != nil {
      t.Fatal(err)
   }
   contentId, err := GetContentId(mediaId)
   if err != nil {
      t.Fatal(err)
   }
   pkgID, destID, err := GetPlaybackDetails(contentId)
   if err != nil {
      t.Fatal(err)
   }
   manifest_url, err := final_tokens.GetManifest(contentId, pkgID, destID)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println("DASH Manifest URL:", manifest_url)
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
   var auth_tokens TokenResponse
   err = json.Unmarshal(data, &auth_tokens)
   if err != nil {
      t.Fatal(err)
   }
   profiles, err := GetProfiles(auth_tokens.AccountId, auth_tokens.AccessToken)
   if err != nil {
      t.Fatal(err)
   }
   i := slices.IndexFunc(profiles, func(p *Profile) bool {
      return p.HasPin == false
   })
   final_tokens, err := ProfileLogin(auth_tokens.RefreshToken, profiles[i].Id)
   if err != nil {
      t.Fatal(err)
   }
   data, err = json.Marshal(final_tokens)
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(cache + "/rosso/crave-final.json", data, os.ModePerm)
   if err != nil {
      t.Fatal(err)
   }
}
func TestPasswordLogin(t *testing.T) {
   username, err := run("credential", "-h=crave.ca", "-k=username")
   if err != nil {
      t.Fatal(err)
   }
   password, err := run("credential", "-h=crave.ca")
   if err != nil {
      t.Fatal(err)
   }
   auth_tokens, err := PasswordLogin(username, password)
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
