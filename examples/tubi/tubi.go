package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/tubi"
   "log"
   "os"
)

type client struct {
   cache          maya.Cache
   WidevineFolder maya.Flag[string]
   ContentId      maya.Flag[int]
   DashId         maya.Flag[string]
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/tubi"); err != nil {
      return err
   }
   if err := maya.ParseFlags(os.Args[1:], c); err != nil {
      return err
   }
   switch {
   case c.WidevineFolder.Set:
      return c.cache.Encode(WidevineFolder(c.WidevineFolder.Value))
   case c.ContentId.Set:
      return c.do_content_id()
   case c.DashId.Set:
      return c.do_dash_id()
   }
   return maya.FormatFlags(os.Stderr, "tubi", c)
}

func (c *client) do_content_id() error {
   content, err := tubi.GetContent(c.ContentId.Value)
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

type WidevineFolder string

func (c *client) do_dash_id() error {
   var (
      manifest maya.Manifest
      server   tubi.LicenseServer
      widevine WidevineFolder
   )
   err := c.cache.Decode(&manifest, &server, &widevine)
   if err != nil {
      return err
   }
   license := func(body []byte) ([]byte, error) {
      return tubi.AcquireLicense(&server, body)
   }
   return maya.DownloadDash(c.DashId.Value, &manifest, &maya.Options{
      Device:  string(widevine),
      Drm:     maya.DrmWidevine,
      License: license,
   })
}
