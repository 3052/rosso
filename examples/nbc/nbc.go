package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/nbc"
   "log"
)

func (c *client) do() error {
   if err := cache.Setup("rosso/nbc.xml"); err != nil {
      return err
   }
   c.cache_err = cache.Read(c)
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   address := maya.StringFlag(&c.address, "a", "address")
   //----------------------------------------------------------
   dash := maya.StringFlag(&c.Job.Dash, "d", "DASH ID")
   if err := maya.ParseFlags(); err != nil {
      return err
   }
   switch {
   case widevine.IsSet:
      return cache.Write(c)
   case address.IsSet:
      return c.do_address()
   case dash.IsSet:
      return c.run(c.do_dash_id)
   }
   return maya.PrintFlags([][]*maya.Flag{{
      widevine,
      address,
      dash,
   }})
}

func (c *client) run(action func() error) error {
   if c.cache_err != nil {
      return c.cache_err
   }
   return action()
}

func (c *client) do_dash_id() error {
   return c.Dash.Download(&c.Job, nbc.FetchWidevine)
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
   stream, err := metadata.Stream()
   if err != nil {
      return err
   }
   c.Dash, err = maya.ListDash(stream.GetManifest)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

type client struct {
   // cache
   Dash *maya.Dash
   Job  maya.Job
   // flags
   address string
   // state
   cache_err error
}
