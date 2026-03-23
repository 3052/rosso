package crave

import (
   "encoding/json"
   "fmt"
   "log"
   "net/http"
   "net/url"
   "os"
   "slices"
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
   publicMovieURL := "https://www.crave.ca/en/movie/goldeneye-38860"
   // Magic happens here
   manifestURL, err := final_tokens.GetManifestFromURL(publicMovieURL)
   if err != nil {
      log.Fatalf("Error retrieving manifest: %v", err)
   }
   fmt.Println("Success!")
   fmt.Println("DASH Manifest URL:", manifestURL)
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
