package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/tubi"
   "log"
   "net/url"
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
   Widevine   maya.FlagString
   content_id maya.FlagInt
   dash_id    maya.FlagString
   bitrate    maya.FlagBool

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
      {Name: "content-id", Value: &c.content_id},
      {Name: "dash-id", Value: &c.dash_id},
      {Name: "bitrate", Value: &c.bitrate, Needs: "dash-id"},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   if flags.IsSet(&c.Widevine) {
      return c.cache.Encode(c)
   }
   if c.content_id >= 1 {
      return c.do_content_id()
   }
   if c.dash_id != "" {
      return c.do_dash_id()
   }
   return flags.Usage(os.Stderr, "tubi")
}

func (c *client) do_content_id() error {
   content, err := tubi.GetContent(int(c.content_id))
   if err != nil {
      return err
   }
   video := content.VideoResources[0]
   address, err := url.Parse(video.Manifest.Url)
   if err != nil {
      return err
   }
   manifest, err := maya.DashList(address)
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest, &video.LicenseServer)
}

func (c *client) do_dash_id() error {
   var (
      manifest maya.Manifest
      server   tubi.LicenseServer
   )
   err := c.cache.Decode(&manifest, &server)
   if err != nil {
      return err
   }
   if c.bitrate {
      return maya.DashBitrate(string(c.dash_id), &manifest)
   }
   license := func(body []byte) ([]byte, error) {
      return tubi.AcquireLicense(&server, body)
   }
   return maya.DashDownload(string(c.dash_id), &manifest, &maya.Options{
      Device:  string(c.Widevine),
      Drm:     maya.DrmWidevine,
      License: license,
   })
}
