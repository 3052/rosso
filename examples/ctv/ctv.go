package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/ctv"
   "log"
)

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
   dash, err := ctv.ParseDash(manifest)
   if err != nil {
      return err
   }
   c.Dash, err = maya.ListDash(dash)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func main() {
   maya.SetProxy("", "*.m4a", "*.m4v")
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   Dash *maya.Dash
   //------------
   Job maya.Job
   //------------
   address string
}

func (c *client) do() error {
   err := cache.Setup("rosso/ctv.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   address := maya.StringFlag(&c.address, "a", "address")
   //----------------------------------------------------------
   dash := maya.StringFlag(&c.Job.Dash, "d", "DASH ID")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
   if err != nil {
      return err
   }
   switch {
   case widevine.IsSet:
      return cache.Write(c)
   case address.IsSet:
      return c.do_address()
   case dash.IsSet:
      return with_cache(c.do_dash_id)
   }
   return maya.PrintFlags([][]*maya.Flag{{
      widevine,
      address,
      dash,
   }})
}

var cache maya.Cache

func (c *client) do_dash_id() error {
   return c.Dash.Download(&c.Job, ctv.FetchWidevine)
}
