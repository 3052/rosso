package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/tubi"
   "fmt"
   "log"
)

type client struct {
   cache    maya.Cache
   flag     maya.FlagSet
   dash     maya.Flag
   tubi_id  maya.Flag
   widevine maya.Flag
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/tubi"); err != nil {
      return err
   }
   c.flag.AddValue(&c.widevine, "w", "Widevine")
   c.flag.AddValue(&c.tubi_id, "t", "Tubi ID")
   c.flag.AddValue(&c.dash, "d", "DASH ID")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   switch {
   case c.widevine.Set:
      return c.cache.Encode(widevine_device(c.widevine.Value))
   case c.tubi_id.Set:
      return c.do_tubi()
   case c.dash.Set:
      return c.do_dash()
   }
   fmt.Println(c.flag)
   return nil
}

func (c *client) do_tubi() error {
   tubi_id, err := c.tubi_id.ParseInt()
   if err != nil {
      return err
   }
   content, err := tubi.GetContent(tubi_id)
   if err != nil {
      return err
   }
   video := content.VideoResources[0]
   manifest, err := maya.ListDash(&video.Manifest.Url.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest, video.LicenseServer)
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
      device   widevine_device
      manifest maya.Manifest
      server   tubi.LicenseServer
   )
   err := c.cache.Decode(&device, &manifest, &server)
   if err != nil {
      return err
   }
   license := func(body []byte) ([]byte, error) {
      return tubi.AcquireLicense(&server, body)
   }
   return maya.DownloadDash(c.dash.Value, &manifest, &maya.Options{
      Device:  string(device),
      Drm:     maya.DrmWidevine,
      License: license,
   })
}
