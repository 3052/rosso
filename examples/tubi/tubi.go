package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/tubi"
   "log"
)

func (c *client) do_tubi() error {
   content, err := tubi.FetchContent(c.tubi_id)
   if err != nil {
      return err
   }
   c.VideoResource = &content.VideoResources[0]
   c.Dash, err = maya.ListDash(c.VideoResource)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_dash() error {
   return c.Dash.Download(&c.Job, c.VideoResource.FetchWidevine)
}

var cache maya.Cache

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   Dash          *maya.Dash
   VideoResource *tubi.VideoResource
   //-------------------------------
   Job maya.Job
   //-------------------------------
   Proxy string
   //-------------------------------
   tubi_id int
}

func (c *client) do() error {
   err := cache.Setup("rosso/tubi.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   //----------------------------------------------------------
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   proxy := maya.StringFlag(&c.Proxy, "x", "proxy")
   //----------------------------------------------------------
   tubi_id := maya.IntFlag(&c.tubi_id, "t", "Tubi ID")
   //------------------------------------------------
   dash := maya.StringFlag(&c.Job.Dash, "d", "DASH ID")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
   err = maya.SetProxy(c.Proxy, "*.mp4")
   if err != nil {
      return err
   }
   switch {
   case widevine.IsSet:
      return cache.Write(c)
   case proxy.IsSet:
      return cache.Write(c)
   case tubi_id.IsSet:
      return c.do_tubi()
   case dash.IsSet:
      return with_cache(c.do_dash)
   }
   return maya.PrintFlags([][]*maya.Flag{{
      widevine,
      proxy,
      tubi_id,
      dash,
   }})
}
