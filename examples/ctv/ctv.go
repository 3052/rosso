package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/ctv"
   "log"
   "os"
)

type client struct {
   cache          maya.Cache
   WidevineFolder maya.Flag[string]
   Address        maya.Flag[string]
   DashId         maya.Flag[string]
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/ctv"); err != nil {
      return err
   }
   if err := maya.ParseFlags(os.Args[1:], c); err != nil {
      return err
   }
   switch {
   case c.WidevineFolder.Set:
      return c.cache.Encode(WidevineFolder(c.WidevineFolder.Value))
   case c.Address.Set:
      return c.do_address()
   case c.DashId.Set:
      return c.do_dash_id()
   }
   return maya.FormatFlags(os.Stderr, "ctv", c)
}

type WidevineFolder string

func (c *client) do_dash_id() error {
   var (
      manifest maya.Manifest
      widevine WidevineFolder
   )
   err := c.cache.Decode(&manifest, &widevine)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.DashId.Value, &manifest, &maya.Options{
      Device:  string(widevine),
      Drm:     maya.DrmWidevine,
      License: ctv.FetchWidevine,
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
   path, err := ctv.GetPath(c.Address.Value)
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
   maya_manifest, err := maya.ListDash(manifest)
   if err != nil {
      return err
   }
   return c.cache.Encode(maya_manifest)
}
