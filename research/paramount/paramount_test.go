package paramount

import (
   "fmt"
   "os/exec"
   "strings"
   "testing"
)

func TestCookie(t *testing.T) {
   username, err := run("credential", "-h=paramountplus.com", "-k=username")
   if err != nil {
      t.Fatal(err)
   }
   password, err := run("credential", "-h=paramountplus.com", "-k=password")
   if err != nil {
      t.Fatal(err)
   }
   const id = "com.cbs.app"
   app_data, ok := Apps[id]
   if !ok {
      t.Fatal(id)
   }
   cookie, err := app_data.CbsCom(username, password)
   if err != nil {
      t.Fatal(err)
   }
   t.Log(cookie)
}

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
