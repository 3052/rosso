package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/ctv"
   "log"
   "net/url"
)

func (c *client) do() error {
   if err := cache.Setup("rosso/ctv.xml"); err != nil {
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

var cache maya.Cache

func (c *client) do_dash_id() error {
   return c.Dash.Download(&c.Job, ctv.FetchWidevine)
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
   c.Dash, err = maya.ListDash(func() (*url.URL, error) {
      return ctv.GetManifest(manifest)
   })
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

type client struct {
   // cache
   Dash *maya.Dash
   Job  maya.Job
   // flags
   address string
   // state
   cache_err error
}
