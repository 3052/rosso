package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/ctv"
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
   flag    maya.FlagSet
}

///

func (c *client) do() error {
   if err := c.cache.Setup("rosso/ctv"); err != nil {
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

func (c *client) do_dash() error {
   var dash maya.Dash
   err := c.cache.Decode(&c.job, &dash)
   if err != nil {
      return err
   }
   return dash.Download(c.dash, &c.job, ctv.FetchWidevine)
}

func (c *client) do_address() error {
   path, err := ctv.GetPath(c.address)
   if err != nil {
      return err
   }
   resolve, err := ctv.Resolve(path)
   if err != nil {
      return err
   }
   axis, err := resolve.AxisContent()
   if err != nil {
      return err
   }
   playback, err := axis.Playback()
   if err != nil {
      return err
   }
   manifest, err := axis.Manifest(playback)
   if err != nil {
      return err
   }
   dash, err := maya.ListDash(manifest)
   if err != nil {
      return err
   }
   return c.cache.Encode(dash)
}
