package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/oldflix"
   "log"
)

func (c *client) do() error {
   return nil
}

type client struct {
   Hls *oldflix.Hls
   //--------------
   Job maya.Job
   //--------------
   oldflix string
   //--------------
   dash_id string
}

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}
