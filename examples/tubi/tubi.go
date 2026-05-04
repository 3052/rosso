package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/tubi"
   "log"
)

func (c *client) do() error {
   err := cache.Setup("rosso/tubi.xml")
   if err != nil {
      return err
   }
   cache_err := cache.Read(c)
   //----------------------------------------------------------
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   tubi_id := maya.IntFlag(&c.tubi_id, "t", "Tubi ID")
   //------------------------------------------------
   dash := maya.StringFlag(&c.Job.Dash, "d", "DASH ID")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
   var (
      action    func() error
      use_cache = true
   )
   switch {
   case widevine.IsSet:
      action = c.do_write
      use_cache = false
   case tubi_id.IsSet:
      action = c.do_tubi
      use_cache = false
   case dash.IsSet:
      action = c.do_dash
   }
   if action != nil {
      if use_cache && cache_err != nil {
         return cache_err
      }
      return action()
   }
   return maya.PrintFlags([][]*maya.Flag{{
      widevine,
      tubi_id,
      dash,
   }})
}

func (c *client) do_write() error {
   return cache.Write(c)
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
   Dash          *maya.Dash
   LicenseServer *tubi.LicenseServer
   //-------------------------------
   Job maya.Job
   //-------------------------------
   tubi_id int
}
