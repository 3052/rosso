package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/nbc"
   "fmt"
   "log"
)

type client struct {
   cache    maya.Cache
   flag     maya.FlagSet
   address  maya.Flag
   dash     maya.Flag
   widevine maya.Flag
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/nbc"); err != nil {
      return err
   }
   c.flag.AddValue(&c.widevine, "w", "Widevine")
   c.flag.AddValue(&c.address, "a", "address")
   c.flag.AddValue(&c.dash, "d", "DASH ID")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   switch {
   case c.widevine.Set:
      return c.cache.Encode(widevine_device(c.widevine.Value))
   case c.address.Set:
      return c.do_address()
   case c.dash.Set:
      return c.do_dash()
   }
   fmt.Println(c.flag)
   return nil
}

type widevine_device string

func (c *client) do_dash() error {
   var (
      device   widevine_device
      manifest maya.Manifest
   )
   err := c.cache.Decode(&device, &manifest)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash.Value, &manifest, &maya.Options{
      Device:  string(device),
      Drm:     maya.DrmWidevine,
      License: nbc.FetchWidevine,
   })
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_address() error {
   address, err := c.address.ParseUrl()
   if err != nil {
      return err
   }
   metadata, err := nbc.FetchMetadata(nbc.GetName(address))
   if err != nil {
      return err
   }
   stream, err := metadata.FetchStream()
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(stream.GetManifest())
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest)
}
