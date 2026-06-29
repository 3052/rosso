package amazon

import "testing"

func Test(t *testing.T) {
   var _ = []struct {
      region   string
      title_id string
      uhd      bool
   }{
      {
         region:   "GB",
         title_id: "amzn1.dv.gti.775a185a-8920-4711-8dbf-d3791538d5af",
         uhd:      false,
      },
      {
         region:   "US",
         title_id: "amzn1.dv.gti.af991753-e4cf-4d28-880d-dfca3d1e8d24",
         uhd:      true,
      },
   }
}
