package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/tubi"
   "log"
)

func (c *client) do_dash() error {
   var (
      dash   maya.Dash
      server tubi.LicenseServer
   )
   err := c.cache.Decode(&c.job, &dash, &server)
   if err != nil {
      return err
   }
   return dash.Download(c.dash, &c.job, func(body []byte) ([]byte, error) {
      return tubi.AcquireLicense(&server, body)
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
   cache   maya.Cache
   dash    string
   job     maya.Job
   tubi_id int
}

func (c *client) do_tubi() error {
   content, err := tubi.GetContent(c.tubi_id)
   if err != nil {
      return err
   }
   video := content.VideoResources[0]
   dash, err := maya.ListDash(&video.Manifest.Url.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(dash, video.LicenseServer)
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/tubi"); err != nil {
      return err
   }
   tubi_id := maya.IntFlag(&c.tubi_id, "t", "Tubi ID")
   widevine := maya.StringFlag(&c.job.Widevine, "w", "Widevine")
   dash := maya.StringFlag(&c.dash, "d", "DASH ID")
   if err := maya.ParseFlags(); err != nil {
      return err
   }
   switch {
   case widevine.IsSet:
      return c.cache.Encode(c.job)
   case tubi_id.IsSet:
      return c.do_tubi()
   case dash.IsSet:
      return c.do_dash()
   }
   return maya.PrintFlags([][]*maya.Flag{{
      widevine,

      tubi_id,
      dash,
   }})
}
