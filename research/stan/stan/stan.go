package main

import (
   "41.neocities.org/maya"
   "log"
   "os"
)

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   Widevine   maya.FlagString
   cache      maya.Cache
   activation maya.FlagBool
}

func (*client) CachePath() string {
   return "rosso/examples/stan/client"
}

func (c *client) do() error {
   if err := c.cache.Setup(); err != nil {
      return err
   }
   if err := c.cache.Decode(c); err != nil {
      return c.cache.Encode(c)
   }
   flags := maya.FlagSet{
      {Name: "widevine-folder", Value: &c.Widevine},
      {Name: "activation", Value: &c.activation},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   if flags.IsSet(&c.Widevine) {
      return c.cache.Encode(c)
   }
   if c.activation {
      return c.do_activation()
   }
   return flags.Usage(os.Stderr, "stan")
}

func (*client) do_activation() error {
   return nil
}
