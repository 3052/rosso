package stan

import "testing"

var program_ids = []int64{
   // play.stan.com.au/programs/1540676
   1540676,
   // play.stan.com.au/programs/1768588
   1768588,
}

func TestProgram(t *testing.T) {
   t.Log(program_ids)
}
