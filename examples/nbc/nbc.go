package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/nbc"
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
   Widevine maya.FlagString

   address maya.FlagString
   dash    maya.FlagString

   cache maya.Cache
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/nbc"); err != nil {
      return err
   }
   if err := c.cache.Decode(c); err != nil {
      if !os.IsNotExist(err) {
         return err
      }
   }
   flags := maya.FlagSet{
      {Name: "widevine-folder", Value: &c.Widevine},
      {Name: "address", Value: &c.address},
      {Name: "dash-id", Value: &c.dash},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   switch {
   case flags.IsSet(&c.Widevine):
      return c.cache.Encode(c)
   case c.address != "":
      return c.do_address()
   case c.dash != "":
      return c.do_dash()
   }
   return flags.Usage(os.Stderr, "nbc")
}

///

func (c *client) do_address() error {
   name, err := nbc.GetName(c.Address.Value)
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
   manifest, err := maya.ListDash(stream.GetManifest())
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest)
}

func (c *client) do_dash() error {
   var (
      manifest maya.Manifest
      widevine WidevineFolder
   )
   err := c.cache.Decode(&manifest, &widevine)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.DashId.Value, &manifest, &maya.Options{
      Device:  widevine.Value,
      Drm:     maya.DrmWidevine,
      License: nbc.FetchWidevine,
   })
}
