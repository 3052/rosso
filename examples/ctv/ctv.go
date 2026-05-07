package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/ctv"
   "log"
   "net/url"
)

func (c *client) do_dash() error {
   if c.err != nil {
      return c.err
   }
   var dash maya.Dash
   err := c.cache.Decode(&dash)
   if err != nil {
      return err
   }
   return dash.Download(&c.job, ctv.FetchWidevine)
}

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
   err     error
   job     maya.Job
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/ctv"); err != nil {
      return err
   }
   address := maya.StringFlag(&c.address, "a", "address")
   c.err = c.cache.Decode(&c.job)
   dash := maya.StringFlag(&c.job.Dash, "d", "DASH ID")
   widevine := maya.StringFlag(&c.job.Widevine, "w", "Widevine")
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
   dash, err := maya.ListDash(func() (*url.URL, error) {
      return ctv.GetManifest(manifest)
   })
   if err != nil {
      return err
   }
   return c.cache.Encode(dash)
}
