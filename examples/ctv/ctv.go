package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/ctv"
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
   if err := c.cache.Setup("rosso/ctv"); err != nil {
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

func (c *client) do_address() error {
   address, err := c.address.ParseUrl()
   if err != nil {
      return err
   }
   resolve, err := ctv.Resolve(address.Path)
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
   maya_manifest, err := maya.ListDash(manifest)
   if err != nil {
      return err
   }
   return c.cache.Encode(maya_manifest)
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type widevine_device string

func (c *client) do_dash() error {
   var (
      manifest maya.Manifest
      device   widevine_device
   )
   err := c.cache.Decode(&device, &manifest)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash.Value, &manifest, &maya.Options{
      Device:  string(device),
      Drm:     maya.DrmWidevine,
      License: ctv.FetchWidevine,
   })
}
