package oldflix

import (
   "os"
   "os/exec"
   "strings"
   "testing"
)

func TestWatch(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(cache + "/rosso/oldflix")
   if err != nil {
      t.Fatal(err)
   }
   var login_data Login
   login_data.Token = string(data)
   // https://oldflix.com.br/browse/play/5d5d54a4d55dc050f8468513
   browse_data, err := login_data.FetchBrowse("5d5d54a4d55dc050f8468513")
   if err != nil {
      t.Fatal(err)
   }
   original, err := browse_data.GetOriginal()
   if err != nil {
      t.Fatal(err)
   }
   watch, err := browse_data.FetchWatch(original.Id, login_data.Token)
   if err != nil {
      t.Fatal(err)
   }
   t.Log(watch)
}

func TestLogin(t *testing.T) {
   username, err := run("credential", "-h=oldflix.com.br", "-k=username")
   if err != nil {
      t.Fatal(err)
   }
   password, err := run("credential", "-h=oldflix.com.br", "-k=password")
   if err != nil {
      t.Fatal(err)
   }
   login_data, err := FetchLogin(username, password)
   if err != nil {
      t.Fatal(err)
   }
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(
      cache+"/rosso/oldflix", []byte(login_data.Token), os.ModePerm,
   )
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
