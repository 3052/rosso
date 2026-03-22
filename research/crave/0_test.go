package crave

import (
   "fmt"
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
   fmt.Println(zero_data)
}
