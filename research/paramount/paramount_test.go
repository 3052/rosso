package paramount

import (
   "fmt"
   "os/exec"
   "slices"
   "strings"
   "testing"
)

func run(name string, arg ...string) (string, error) {
   var data strings.Builder
   command := exec.Command(name, arg...)
   command.Stdout = &data
   fmt.Println(command.Args)
   err := command.Run()
   if err != nil {
      return "", err
   }
   return data.String(), nil
}

func TestCookie(t *testing.T) {
   username, err := run("credential", "-h=paramountplus.com", "-k=username")
   if err != nil {
      t.Fatal(err)
   }
   password, err := run("credential", "-h=paramountplus.com", "-k=password")
   if err != nil {
      t.Fatal(err)
   }
   i := slices.IndexFunc(Apps, func(a *app) bool {
      return a.id == "com.cbs.app"
   })
   cookie, err := Apps[i].CbsCom(username, password)
   if err != nil {
      t.Fatal(err)
   }
   t.Log(cookie)
}
