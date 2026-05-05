package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/tubi"
   "log"
)

func (c *client) do() error {
   if err := cache.Setup("rosso/tubi.xml"); err != nil {
      return err
   }
   c.cache_err = cache.Read(c)
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   tubi_id := maya.IntFlag(&c.tubi_id, "t", "Tubi ID")
   //------------------------------------------------
   dash := maya.StringFlag(&c.Job.Dash, "d", "DASH ID")
   if err := maya.ParseFlags(); err != nil {
      return err
   }
   switch {
   case widevine.IsSet:
      return cache.Write(c)
   case tubi_id.IsSet:
      return c.do_tubi()
   case dash.IsSet:
      return c.run(c.do_dash)
   }
   return maya.PrintFlags([][]*maya.Flag{{
      widevine,
      tubi_id,
      dash,
   }})
}

func (c *client) run(action func() error) error {
   if c.cache_err != nil {
      return c.cache_err
   }
   return action()
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
   c.Dash, err = maya.ListDash(video.GetManifest)
   if err != nil {
      return err
   }
   c.LicenseServer = &video.LicenseServer
   return cache.Write(c)
}

func (c *client) do_dash() error {
   return c.Dash.Download(&c.Job, func(body []byte) ([]byte, error) {
      return tubi.AcquireLicense(c.LicenseServer, body)
   })
}

var cache maya.Cache

type client struct {
   // cache
   Dash          *maya.Dash
   Job           maya.Job
   LicenseServer *tubi.LicenseServer
   // flags
   tubi_id int
   // state
   cache_err error
}
