package oldflix

import (
   "os"
   "os/exec"
   "strings"
   "testing"
)

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
      cache + "/rosso/oldflix", []byte(login_data.Token), os.ModePerm,
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
