package paramount

import (
   "os"
   "testing"
)

func Test(t *testing.T) {
   data, err := os.ReadFile("classes.dex")
   if err != nil {
      t.Fatal(err)
   }
   for _, result := range ExtractHexBytes(data) {
      t.Log(result)
   }
}
