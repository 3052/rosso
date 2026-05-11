package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/tubi"
   "log"
)

func (c *client) do_dash() error {
   var (
      manifest maya.Manifest
      server   tubi.LicenseServer
      widevine device
   )
   err := c.cache.Decode(&manifest, &server, &widevine)
   if err != nil {
      return err
   }
   license := func(body []byte) ([]byte, error) {
      return tubi.AcquireLicense(&server, body)
   }
   return maya.DownloadDash(c.dash, &manifest, &maya.Options{
      Device:  string(widevine),
      Drm:     maya.DrmWidevine,
      License: license,
   })
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   cache    maya.Cache
   dash     string
   flag     maya.FlagSet
   tubi_id  int
   widevine string
}

type device string

func (c *client) do() error {
   if err := c.cache.Setup("rosso/tubi"); err != nil {
      return err
   }
   tubi_id := c.flag.Int(&c.tubi_id, "t", "Tubi ID")
   dash := c.flag.String(&c.dash, "d", "DASH ID")
   widevine := c.flag.String(&c.widevine, "w", "Widevine")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   switch {
   case widevine.IsSet:
      return c.cache.Encode(device(c.widevine))
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

func (c *client) do_tubi() error {
   content, err := tubi.GetContent(c.tubi_id)
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
