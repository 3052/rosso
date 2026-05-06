package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/tubi"
   "log"
)

type client struct {
   cache   maya.Cache
   job     maya.Job
   tubi_id int
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_tubi() error {
   content, err := tubi.GetContent(c.tubi_id)
   if err != nil {
      return err
   }

   video := content.VideoResources[0]
   if err := c.cache.Encode(video.LicenseServer); err != nil {
      return err
   }

   dash, err := maya.ListDash(video.GetManifest)
   if err != nil {
      return err
   }

   return c.cache.Encode(dash)
}

func (c *client) do_dash() error {
   var dash maya.Dash
   if err := c.cache.Decode(&dash); err != nil {
      return err
   }

   var server tubi.LicenseServer
   if err := c.cache.Decode(&server); err != nil {
      return err
   }

   return dash.Download(&c.job, func(body []byte) ([]byte, error) {
      return tubi.AcquireLicense(&server, body)
   })
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/tubi"); err != nil {
      return err
   }
   c.cache.Decode(&c.job)

   widevine := maya.StringFlag(&c.job.Widevine, "w", "Widevine")
   tubi_id := maya.IntFlag(&c.tubi_id, "t", "Tubi ID")
   dash := maya.StringFlag(&c.job.Dash, "d", "DASH ID")

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
