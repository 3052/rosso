package crave

import (
   "fmt"
   "os/exec"
   "strings"
   "testing"
)

func TestZero(t *testing.T) {
   username, err := run("credential", "-h=crave.ca", "-k=username")
   if err != nil {
      t.Fatal(err)
   }
   password, err := run("credential", "-h=crave.ca")
   if err != nil {
      t.Fatal(err)
   }
   zero_data, err := fetch_zero(username, password)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", zero_data)
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
