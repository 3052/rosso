package stan

import (
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
      time.Sleep(time.Second)
   }
}
