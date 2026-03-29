package paramount

import (
   "os"
   "testing"
)

func Test(t *testing.T) {
   data, err := os.ReadFile("base.apk")
   if err != nil {
      t.Fatal(err)
   }
   results, err := ExtractDexHexBytes(data)
   if err != nil {
      t.Fatal(err)
   }
   for result := range results {
      t.Log(result)
   }
}
