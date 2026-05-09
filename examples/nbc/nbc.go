package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/nbc"
   "log"
)

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   address string
   cache   maya.Cache
   dash    string
   job     maya.Job
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/nbc"); err != nil {
      return err
   }
   address := maya.StringFlag(&c.address, "a", "address")
   widevine := maya.StringFlag(&c.job.Widevine, "w", "Widevine")
   dash := maya.StringFlag(&c.dash, "d", "DASH ID")
   if err := maya.ParseFlags(); err != nil {
      return err
   }
   switch {
   case widevine.IsSet:
      return c.cache.Encode(c.job)
   case address.IsSet:
      return c.do_address()
   case dash.IsSet:
      return c.do_dash()
   }
   return maya.PrintFlags([][]*maya.Flag{{
      widevine,

      address,
      dash,
   }})
}

///

func (c *client) do_address() error {
   name, err := nbc.GetName(c.address)
   if err != nil {
      return err
   }
   metadata, err := nbc.FetchMetadata(name)
   if err != nil {
      return err
   }
   stream, err := metadata.Stream()
   if err != nil {
      return err
   }
   c.Dash, err = maya.ListDash(stream.GetManifest)
   if err != nil {
      return err
   }
   return c.cache.Write(c)
}

func (c *client) do_dash() error {
   return c.Dash.Download(&c.Job, nbc.FetchWidevine)
}
