package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/tubi"
   "log"
)

type client struct {
   cache    maya.Cache
   dash     string
   flag     maya.FlagSet
   tubi_id  int
   widevine string
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/tubi"); err != nil {
      return err
   }
   tubi_id := c.flag.Int(&c.tubi_id, "t", "Tubi ID")
   dash := c.flag.String(&c.dash, "d", "DASH ID")
   widevine := c.flag.String(&c.widevine, "w", "Widevine folder")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   switch {
   case widevine.IsSet:
      return c.cache.Encode(widevine_device(c.widevine))
   case tubi_id.IsSet:
      return c.do_tubi()
   case dash.IsSet:
      return c.do_dash()
   }
   return maya.PrintFlags([]maya.FlagSet{{
      widevine,
      tubi_id,
      dash,
   }})
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
