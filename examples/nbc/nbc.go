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
   address  string
   cache    maya.Cache
   dash     string
   flag     maya.FlagSet
   widevine string
}

type device string

///

func (c *client) do() error {
   if err := c.cache.Setup("rosso/nbc"); err != nil {
      return err
   }
   address := c.flag.String(&c.address, "a", "address")
   dash := c.flag.String(&c.dash, "d", "DASH ID")
   widevine := c.flag.String(&c.widevine, "w", "Widevine")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   switch {
   case widevine.IsSet:
      return c.cache.Encode(device(c.widevine))
   case address.IsSet:
      return c.do_address()
   case dash.IsSet:
      return c.do_dash()
   }
   return maya.PrintFlags([]maya.FlagSet{{
      widevine,
      address,
      dash,
   }})
}

func (c *client) do_address() error {
   name, err := nbc.GetName(c.address)
   if err != nil {
      return err
   }
   metadata, err := nbc.FetchMetadata(name)
   if err != nil {
      return err
   }
   stream, err := metadata.FetchStream()
   if err != nil {
      return err
   }
   dash, err := maya.ListDash(stream.GetManifest())
   if err != nil {
      return err
   }
   return c.cache.Encode(dash)
}

func (c *client) do_dash() error {
   var dash maya.Dash
   err := c.cache.Decode(&c.job, &dash)
   if err != nil {
      return err
   }
   return dash.Download(c.dash, &c.job, nbc.FetchWidevine)
}
