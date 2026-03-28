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
   with_cache := cache.Read(c)
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   tubi_id := maya.IntFlag(&c.tubi_id, "t", "Tubi ID")
   //------------------------------------------------
   dash_id := maya.StringFlag(&c.dash_id, "d", "DASH ID")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
   switch {
   case widevine.IsSet:
      return cache.Write(c)
   case tubi_id.IsSet:
      return c.do_tubi()
   case dash_id.IsSet:
      return with_cache(c.do_dash_id)
   }
   return maya.PrintFlags([][]*maya.Flag{{
      widevine,
      tubi_id,
      dash_id,
   }})
}

func (c *client) do_tubi() error {
   content, err := tubi.FetchContent(c.tubi_id)
   if err != nil {
      return err
   }
   c.VideoResource = &content.VideoResources[0]
   c.Dash, err = c.VideoResource.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(
      c.Dash.Body, c.Dash.Url, c.dash_id, c.VideoResource.Widevine,
   )
}

var cache maya.Cache

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.mp4")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   Dash          *tubi.Dash
   VideoResource *tubi.VideoResource
   //-------------------------------
   Job maya.Job
   //-------------------------------
   tubi_id int
   //-------------------------------
   dash_id string
}
