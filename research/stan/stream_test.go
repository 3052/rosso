package stan

import (
   "154.pages.dev/widevine"
   "encoding/hex"
   "fmt"
   "os"
   "testing"
)

// play.stan.com.au/programs/1768588
const (
   key_id = "0b5c271e61c244a8ab81e8363a66aa35"
   program_id = 1768588
)

func TestStream(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   var token WebToken
   token.Data, err = os.ReadFile(home + "/stan.json")
   if err != nil {
      t.Fatal(err)
   }
   token.Unmarshal()
   session, err := token.Session()
   if err != nil {
      t.Fatal(err)
   }
   stream, err := session.Stream(program_id)
   if err != nil {
      t.Fatal(err)
   }
   private_key, err := os.ReadFile(home + "/widevine/private_key.pem")
   if err != nil {
      t.Fatal(err)
   }
   client_id, err := os.ReadFile(home + "/widevine/client_id.bin")
   if err != nil {
      t.Fatal(err)
   }
   var pssh widevine.PSSH
   pssh.KeyId, err = hex.DecodeString(key_id)
   if err != nil {
      t.Fatal(err)
   }
   var module widevine.CDM
   err = module.New(private_key, client_id, pssh.Encode())
   if err != nil {
      t.Fatal(err)
   }
   key, err := module.Key(stream, pssh.KeyId)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%x\n", key)
}
