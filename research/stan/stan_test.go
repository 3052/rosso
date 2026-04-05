package stan

import (
   "154.pages.dev/text"
   "fmt"
   "os"
   "testing"
   "time"
)

var program_ids = []int64{
   // play.stan.com.au/programs/1540676
   1540676,
   // play.stan.com.au/programs/1768588
   1768588,
}

func TestProgram(t *testing.T) {
   for _, program_id := range program_ids {
      var program LegacyProgram
      err := program.New(program_id)
      if err != nil {
         t.Fatal(err)
      }
      name, err := text.Name(Namer{program})
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%q\n", name)
      time.Sleep(time.Second)
   }
}

func TestCode(t *testing.T) {
   var code ActivationCode
   err := code.New()
   if err != nil {
      t.Fatal(err)
   }
   code.Unmarshal()
   fmt.Println(code)
   os.WriteFile("code.json", code.Data, 0666)
}

func TestToken(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   var code ActivationCode
   code.Data, err = os.ReadFile("code.json")
   if err != nil {
      t.Fatal(err)
   }
   code.Unmarshal()
   token, err := code.Token()
   if err != nil {
      t.Fatal(err)
   }
   os.WriteFile(home + "/stan.json", token.Data, 0666)
}
