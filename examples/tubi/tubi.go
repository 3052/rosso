package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/tubi"
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
   Proxy      maya.FlagString
   Widevine   maya.FlagString
   content_id maya.FlagInt
   dash       maya.FlagString

   cache maya.Cache
}

func (*client) CachePath() string {
   return "rosso/examples/tubi/client"
}

func (c *client) do() error {
   if err := c.cache.Setup(); err != nil {
      return err
   }
   if err := c.cache.Decode(c); err != nil {
      return c.cache.Encode(c)
   }
   flags := maya.FlagSet{
      {Name: "widevine-folder", Value: &c.Widevine},
      {Name: "proxy", Value: &c.Proxy},
      {Name: "content-id", Value: &c.content_id},
      {Name: "dash-id", Value: &c.dash},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   if flags.IsSet(&c.Widevine) {
      return c.cache.Encode(c)
   }
   if flags.IsSet(&c.Proxy) {
      return c.cache.Encode(c)
   }
   if err := maya.SetProxy(string(c.Proxy)); err != nil {
      return err
   }
   if c.content_id >= 1 {
      return c.do_content_id()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return flags.Usage(os.Stderr, "tubi")
}

func (c *client) do_content_id() error {
   content, err := tubi.GetContent(int(c.content_id))
   if err != nil {
      return err
   }
   video := content.VideoResources[0]
   manifest, err := maya.ListDash(&video.Manifest.Url.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest, &video.LicenseServer)
}

func (c *client) do_dash() error {
   var (
      manifest maya.Manifest
      server   tubi.LicenseServer
   )
   err := c.cache.Decode(&manifest, &server)
   if err != nil {
      return err
   }
   license := func(body []byte) ([]byte, error) {
      return tubi.AcquireLicense(&server, body)
   }
   return maya.DownloadDash(string(c.dash), &manifest, &maya.Options{
      Device:  string(c.Widevine),
      Drm:     maya.DrmWidevine,
      License: license,
   })
}
