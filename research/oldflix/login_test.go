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
   client_data := NewClient()
   err = client_data.Login(username, password)
   if err != nil {
      t.Fatal(err)
   }
   t.Log(client_data.Token)
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(
      cache + "/rosso/oldflix", []byte(client_data.Token), os.ModePerm,
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
